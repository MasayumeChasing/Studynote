// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: httpserver.go
// Description: Contians http server functions

// Package httpserver implements http server functions
package httpserver

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/api/amsgrpc"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"context"
	"net/http"
	"strconv"
	"time"

	pb "codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/apis/golang/ams"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	server *http.Server
)

// Run 启动http服务
func RunServer(am *amsgrpc.AMSgRPCInterface) error {
	ctx := context.Background()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) { // 允许所有请求头字段放进来
		return s, true
	}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	rpcAddr := constants.EnvConfig.ListenIP + ":" + strconv.Itoa(int(configmanager.GetConf().Server.RpcPort))
	err := pb.RegisterAlertManagerServiceHandlerFromEndpoint(ctx, mux, rpcAddr, opts)
	if err != nil {
		return err
	}

	httpAddr := constants.EnvConfig.ListenIP + ":" + strconv.Itoa(int(configmanager.GetConf().Server.HttpPort))
	server = &http.Server{Addr: httpAddr, Handler: mux}
	zaplog.Info("ams http server start at ", httpAddr)
	err = server.ListenAndServe()
	return err
}

// Close 关闭http服务
func Close() error {
	if server == nil {
		return nil
	}

	timeout := configmanager.GetConf().Server.ForciblyExitWaitTime
	if timeout == 0 {
		timeout = 10 // Exit wait time
	}

	// 使用context控制srv.Shutdown的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
