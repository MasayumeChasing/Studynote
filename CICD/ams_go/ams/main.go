// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: main.go
// Description: Contians ams start (main) functions

package main

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/api/amsgrpc"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/grpcserver"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/httpserver"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/mqmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/taskexecutor"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/taskmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("main panic recovered, panic=%+v\n", err)
			// 堆栈已经没了, 可以打印一些 统计信息
		}
	}()
	// init log
	zaplog.Init("ams")
	defer zaplog.Close()
	zaplog.Info("zapLog init ok.")

	zaplog.Info("AMS started")

	err := constants.LoadEnvConfig()
	if err != nil {
		zaplog.Fatal("[ams] Loading env config failed")
	}
	zaplog.Info("[ams] Env config loaded")

	// init scc
	err = sccapi.Initialize(constants.EnvConfig.SccConfPath)
	if err != nil {
		zaplog.Fatalf("[ams] init scc fail, %s", err)
	}
	zaplog.Info("[ams] scc init ok")
	defer func() {
		err := sccapi.Finalize()
		if err != nil {
			zaplog.Fatalf("[ams] scc close, err:%+v", err)
		}
	}()

	am, tm, conf := InitAMS()
	if am == nil || tm == nil || conf == nil {
		zaplog.Fatalf("[ams] Init error, err:%+v", err)
	}
	zaplog.Info("Initialized AMS")
	zaplog.Info("Initializing and starting gRPC Server")

	quit := make(chan error, 1)
	exit := make(chan error, 1)
	var wg sync.WaitGroup

	// init rpc server
	wg.Add(1)
	go RpcServer(am, quit, &wg)

	// init http server
	wg.Add(1)
	go HttpServer(am, quit, &wg)

	wg.Add(1)
	go AMSServer(tm, conf, quit, &wg)

	// 等待信号退出, 或者, 等待服务因为错误退出
	waitExit(quit)

	wg.Add(1)
	go CloseServer(&wg)

	// 阻塞直到退出都已经完成
	go func() {
		wg.Wait()
		close(exit)
	}()

	waitGraceExit(exit)
}

// InitAMS
func InitAMS() (*amsgrpc.AMSgRPCInterface, *taskmanager.TaskManager,
	*configmanager.AMSConfig) {
	conf, err := configmanager.NewAMSConfig()
	if err != nil {
		newErr := fmt.Errorf("[ams] Loading config failed, %s", err)
		zaplog.Error(newErr)
		return nil, nil, nil
	}
	zaplog.Info("[ams] AMS config loaded")

	tmConf := taskmanager.TaskManagerConf{
		SLMService: conf.SLMService,
		DBInfo:     conf.MongoDB,
		JwtSLM:     conf.JwtSLM,
	}
	tm, err := taskmanager.NewTaskManager(tmConf)
	if err != nil {
		newErr := fmt.Errorf("[ams] task manager error, %s", err)
		zaplog.Error(newErr)
		return nil, nil, nil
	}
	gRPCConf := amsgrpc.GRPCInterfaceConf{
		GRPC:   conf.Server,
		JwtAMS: conf.JwtAMS,
	}
	am, err := amsgrpc.NewAMSgRPCInterface(gRPCConf, tm)
	if err != nil {
		newErr := fmt.Errorf("[ams] ams init failed")
		zaplog.Error(newErr)
		return nil, nil, nil
	}

	return am, tm, conf
}

// RPC server go routine
func RpcServer(am *amsgrpc.AMSgRPCInterface, quit chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	err := grpcserver.RunServer(am)
	if err != nil {
		zaplog.Info("ams rpc server exited. err:", err)
		if quit != nil {
			quit <- err
		}
	}
}

// HTTP server go routine
func HttpServer(am *amsgrpc.AMSgRPCInterface, quit chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	err := httpserver.RunServer(am)
	if err != nil {
		zaplog.Info("ams http server exited. err:", err)
		if quit != nil {
			quit <- err
		}
	}
}

// AMS server go routine
func AMSServer(tm *taskmanager.TaskManager, conf *configmanager.AMSConfig,
	quit chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	ret := StartAMS(conf, tm)
	err := fmt.Errorf("StartAMS exit with code [%d]", ret)
	if quit != nil {
		quit <- err
	}
}

// Close server go routine
func CloseServer(wg *sync.WaitGroup) {
	defer wg.Done()
	// 处理退出
	grpcserver.Close()
	err := httpserver.Close()
	if err != nil {
		zaplog.Errorf("http close err=%+v", err)
	}
}

// Start AMS
func StartAMS(conf *configmanager.AMSConfig, tm *taskmanager.TaskManager) int {
	te, err := taskexecutor.NewTaskExecutor(&conf.EngineInfo)
	if err != nil {
		zaplog.Fatalf("Task executor initialization failed")
		return -1
	}
	var imq *mqmanager.MQManager // MQ from which AMS receives input for search task
	var omq *mqmanager.MQManager // MQ to which AMS raises alarm based on search result
	imq, err = mqmanager.NewMQManager(&conf.InMsgSearchMQ)
	if err != nil {
		zaplog.Fatalf("MQ initialization failed for Input Consumer MQ")
		return -1
	}
	omq, err = mqmanager.NewMQManager(&conf.OutMsgAlarmMQ)
	if err != nil {
		zaplog.Fatalf("MQ initialization failed for output Producer MQ")
		return -1
	}
	mqmsgs, err := imq.ReceiveMsgs()
	if err != nil {
		zaplog.Fatalf("Consuming msg from MQ failed")
		return -1
	}
	zaplog.Info("Waiting for msg from MQ")
	for mqmsg := range mqmsgs {
		err := processMQMsg(string(mqmsg.Body), omq, te, tm)
		if err != nil {
			// Print errors and continue
			zaplog.Infof("Processing MQ msg failed, %s", err)
		}
		err = mqmsg.Ack(true)
		if err != nil {
			zaplog.Infof("Sending ACK confirmation to MQ failed, %s", err)
		}
	}
	return 0
}

func processMQMsg(msg string, omq *mqmanager.MQManager,
	te *taskexecutor.TaskExecutor, tm *taskmanager.TaskManager) error {
	defer func() {
		if err := recover(); err != nil {
			zaplog.Info("Process MQ Msg panic recovered, panic=%+v\n", err)
		}
	}()
	// Call Task executor and get search result as json
	// Then Send result in ouput mq
	result, err := te.DoImgSearch(tm, msg)
	if err != nil {
		zaplog.Errorf("Img search failed, %s", err)
		return err
	}
	if len(result) > 0 {
		err := omq.SendMsg(result)
		if err != nil {
			zaplog.Info("Img search result publish to MQ failed")
			return err
		} else {
			zaplog.Info("Published the Img Search result on Alarm MQ")
		}
	}
	return nil
}

func waitGraceExit(exit chan error) {
	if exit == nil {
		return
	}
	select {
	case _, ok := <-time.After(constants.ForciblyExitWaitTime):
		if !ok {
			zaplog.Error("channel read error, ams is exiting")
			return
		}
		zaplog.Info("ams force exited")
		return
	case _, ok := <-exit:
		if !ok {
			zaplog.Error("channel read error, ams is exiting")
			return
		}
		zaplog.Info("ams gracefully exited")
		return
	}
}

func waitExit(quit chan error) {
	if quit == nil {
		return
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	for {
		select {
		case sig, ok := <-sigCh:
			if !ok {
				return
			}
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT: // SIGKILL 信号无法捕捉和忽略
				zaplog.Infof("got signal and will to exit: %d, signal name = %s", sig, sig.String())
				return
			default:
				zaplog.Infof("other signal: %s", sig.String())
			}
		case err, ok := <-quit:
			if !ok {
				return
			}
			switch err {
			case nil:
				zaplog.Infof("wait will be exited")
				return
			default:
				zaplog.Errorf("wait will be exited, caused by err:%s", err)
				return
			}
		}
	}
}
