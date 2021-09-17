// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: configmanager.go
// Description: Contians config manager functions

// Package configmanager implements configuration management functions
package configmanager

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/apollo"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"fmt"
	"strconv"
)

const (
	MQProducer         = "producer"
	MQConsumer         = "consumer"
	MQProducerExchange = "ivm.ams.alarm.exchange"
	MQConsumerExchange = "device.alarm.exchange"
)

type AMSConfig struct {
	apr           *apollo.ApolloReader
	Server        ServerInfo
	InMsgSearchMQ MQInfo
	OutMsgAlarmMQ MQInfo
	MongoDB       MongoDBInfo
	parseErr      error
	SLMService    SLMServiceInfo
	JwtAMS        JwtInfo // For gRPC APIs in AMS
	JwtSLM        JwtInfo // For gRPC APIs in SLM
	EngineInfo    ESEngineInfo
}

type ServerInfo struct {
	RpcPort              uint16
	HttpPort             uint16
	ForciblyExitWaitTime int
}

type MQInfo struct {
	IP           string
	Port         uint16
	UserName     string
	Password     string
	ExchangeName string
	Type         string
	QueueName    string
	BindingKey   string
	RoutingKey   string
	Mode         string
}

type MongoDBInfo struct {
	IP       string
	Port     uint16
	UserName string
	Password string
}

type SLMServiceInfo struct {
	IP   string
	Port uint16
}

type JwtInfo struct {
	Secret string
}

type ESEngineInfo struct {
	Endpoints []string
	Username string
	Password string
}

// Maintaining as global variable because callback (readCommonConfig) needs to
// access it. There is no way to register callback with arguments in apollo.
var globalConf AMSConfig
var globalConfInited bool = false

// Get Config object.
func GetConf() *AMSConfig {
	return &globalConf
}

// New AMS Config
func NewAMSConfig() (*AMSConfig, error) {
	var err error
	if globalConfInited == false {
		globalConf.apr, err = apollo.NewApolloReader()
		if err != nil {
			newErr := fmt.Errorf("Apollo client init failed, %s", err)
			zaplog.Error(newErr)
			return nil, newErr
		}
		readConfig()
		if globalConf.parseErr != nil {
			newErr := fmt.Errorf("Apollo fetch config failed,",
				globalConf.parseErr)
			return nil, newErr
		}
		globalConf.apr.AddNSListener(constants.ApolloNsCommon, readCommonConfig)
		globalConf.apr.AddNSListener(constants.ApolloNsRabbitMQ,
			readRabbitMQConfig)
		globalConfInited = true
	}
	return &globalConf, nil
}

func readConfig() {
	readCommonConfig()
	readRabbitMQConfig()
	// TODO handle parseERR
	// TODO need to reinit MongoDB, MQ based on new config
}

func getPort(ns string, k string) uint16 {
	portStr := globalConf.apr.Read(ns, k)
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		globalConf.parseErr = fmt.Errorf("[%s] Port [%s] is invalid, %s",
			k, portStr, err.Error())
		return 0
	}
	if portInt != int(uint16(portInt)) {
		globalConf.parseErr = fmt.Errorf("Invalid port number [%d] in [%s]",
			portInt, k)
		return 0
	}
	if portInt == 0 {
		globalConf.parseErr = fmt.Errorf("Zero is configured as port in [%s]",
			k)
		return 0
	}
	return uint16(portInt)
}

func readCommonConfig() {
	zaplog.Info("Reading common config")
	ns := constants.ApolloNsCommon
	port := getPort(ns, constants.KeyGRPCPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid GRPC Port,", globalConf.parseErr)
		return
	}
	globalConf.Server.RpcPort = port

	port = getPort(ns, constants.KeyHttpPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid HTTP Port,", globalConf.parseErr)
		return
	}
	globalConf.Server.HttpPort = port

	globalConf.MongoDB.IP = globalConf.apr.Read(ns, constants.KeyMongodbIP)
	port = getPort(ns, constants.KeyMongodbPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid Mongodb Port,", globalConf.parseErr)
		return
	}
	globalConf.MongoDB.Port = port
	globalConf.MongoDB.UserName =
		globalConf.apr.Read(ns, constants.KeyMongodbUsername)
	globalConf.MongoDB.Password =
		globalConf.apr.Read(ns, constants.KeyMongodbPassword)

	globalConf.SLMService.IP = globalConf.apr.Read(ns, constants.KeySlmServerIP)
	port = getPort(ns, constants.KeySlmServerPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid SLMService Port,", globalConf.parseErr)
		return
	}
	globalConf.SLMService.Port = port

	globalConf.JwtAMS.Secret = globalConf.apr.Read(ns, constants.KeyAMSJwtSecret)
	globalConf.JwtSLM.Secret = globalConf.apr.Read(ns, constants.KeySLMJwtSecret)

	timeoutStr := globalConf.apr.Read(ns, constants.KeyServerForciblyExitWaitTime)
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		zaplog.Error("Invalid ForciblyExitWaitTime,", globalConf.parseErr)
		return
	}
	globalConf.Server.ForciblyExitWaitTime = timeout
	globalConf.EngineInfo.Endpoints = []string{globalConf.apr.Read(ns,
                                           constants.KeyEngineEndpoints)}
	globalConf.EngineInfo.Username = globalConf.apr.Read(ns,
                                                    constants.KeyEngineUsername)
	globalConf.EngineInfo.Password = globalConf.apr.Read(ns,
                                                    constants.KeyEnginePassword)
}

func readRabbitMQConfig() {
	zaplog.Info("Reading rabbitmq config")
	ns := constants.ApolloNsRabbitMQ
	globalConf.InMsgSearchMQ.IP =
		globalConf.apr.Read(ns, constants.KeyMQInMsgIP)
	port := getPort(ns, constants.KeyMQInMsgPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid InMsgMQ Port", globalConf.parseErr)
		return
	}
	globalConf.InMsgSearchMQ.Port = port
	globalConf.InMsgSearchMQ.UserName =
		globalConf.apr.Read(ns, constants.KeyMQInMsgUsername)
	globalConf.InMsgSearchMQ.Password =
		globalConf.apr.Read(ns, constants.KeyMQInMsgPassword)
	globalConf.InMsgSearchMQ.Type =
		globalConf.apr.Read(ns, constants.KeyMQInMsgType)
	globalConf.InMsgSearchMQ.QueueName =
		globalConf.apr.Read(ns, constants.KeyMQInMsgQueueName)
	globalConf.InMsgSearchMQ.BindingKey =
		globalConf.apr.Read(ns, constants.KeyMQInMsgBindingKey)
	globalConf.InMsgSearchMQ.RoutingKey =
		globalConf.apr.Read(ns, constants.KeyMQInMsgRoutingKey)
	globalConf.InMsgSearchMQ.Mode = MQConsumer
	globalConf.InMsgSearchMQ.ExchangeName = MQConsumerExchange

	globalConf.OutMsgAlarmMQ.IP =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgIP)
	port = getPort(ns, constants.KeyMQInMsgPort)
	if globalConf.parseErr != nil {
		zaplog.Error("Invalid OutMsgMQ Port", globalConf.parseErr)
		return
	}
	globalConf.OutMsgAlarmMQ.Port = port
	globalConf.OutMsgAlarmMQ.UserName =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgUsername)
	globalConf.OutMsgAlarmMQ.Password =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgPassword)
	globalConf.OutMsgAlarmMQ.Type =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgType)
	globalConf.OutMsgAlarmMQ.QueueName =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgQueueName)
	globalConf.OutMsgAlarmMQ.BindingKey =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgBindingKey)
	globalConf.OutMsgAlarmMQ.RoutingKey =
		globalConf.apr.Read(ns, constants.KeyMQOutMsgRoutingKey)
	globalConf.OutMsgAlarmMQ.Mode = MQProducer
	globalConf.OutMsgAlarmMQ.ExchangeName = MQProducerExchange
}
// Deinit Config Manager
func (amsconf *AMSConfig) Deinit() {
	amsconf.apr.Close()
}
