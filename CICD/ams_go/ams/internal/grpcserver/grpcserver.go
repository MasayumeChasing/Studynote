// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: grpcserver.go
// Description: Contians grpc server functions

// Package grpcserver implements grpc server functions
package grpcserver

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/api/amsgrpc"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/apis/golang/ams"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/x/rpc"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var (
	grpcServer *grpc.Server
)

// Run grpc Server
func RunServer(am *amsgrpc.AMSgRPCInterface) error {
	netListener, err := getNetListener(constants.EnvConfig.ListenIP,
		am.Conf.GRPC.RpcPort)
	if err != nil {
		return err
	}

	stat := rpc.NewStat()
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(rpc.Recover()),
			grpc.ChainUnaryInterceptor(rpc.DebugPrint, rpc.ConvertError),
			grpc.StatsHandler(stat))
	ams.RegisterAlertManagerServiceServer(grpcServer, am)
	zaplog.Info("Staring gRPC Server...")
	err = grpcServer.Serve(netListener)
	if err != nil {
		newErr := fmt.Errorf("gRPC Server stopped: %s", err)
		zaplog.Error(newErr)
		return newErr
	}
	return nil
}

// Close 关闭grpc服务
func Close() {
    if grpcServer != nil {
	    grpcServer.GracefulStop()
    }
}

func getNetListener(ip string, port uint16) (net.Listener, error) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		newErr := fmt.Errorf("TCP Listen failed : %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	return listen, nil
}
