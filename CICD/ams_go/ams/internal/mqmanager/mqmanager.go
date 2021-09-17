// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name:   mqmanager.go
// Description: Msg Publish and Subscribe operation on Rabbit Message Queue

// Package mqmanager implements mq management functions
package mqmanager

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/middleware/rabbitmq"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"fmt"
	"github.com/streadway/amqp"
	"net/url"
	"strconv"
)

type MQManager struct {
	mqinfo   *configmanager.MQInfo
	mqClient rabbitmq.MQ
}

// Create new instance of MQManager structure
func NewMQManager(info *configmanager.MQInfo) (*MQManager, error) {
	zaplog.Info("MQManager initialization for Exchange", info.ExchangeName)
	mqClient, err := openMQConnection(info)
	if err != nil {
		return nil, err
	}
	zaplog.Info("MQ connection established")
	if info.Mode == configmanager.MQProducer {
		err = mqClient.Declare(func(ch *rabbitmq.Chan) error {
			err := ch.DeclareExchange(info.ExchangeName, amqp.ExchangeTopic)
			if err != nil {
				return fmt.Errorf("declare exchange err:%w", err)
			}

			return ch.DeclareQueueAndBind(info.QueueName, info.BindingKey, info.ExchangeName)
		})
	} else {
		err = mqClient.DeclareQueueTTL(info.QueueName, info.BindingKey, info.ExchangeName, 0)
	}
	if err != nil {
		mqClient.Close()
		return nil, err
	}
	return &MQManager{
		mqinfo:   info,
		mqClient: mqClient,
	}, nil
}

func openMQConnection(info *configmanager.MQInfo) (rabbitmq.MQ, error) {
	if len(info.UserName) == 0 || len(info.Password) == 0 {
		newErr := fmt.Errorf("Parameter missing for MQ connection")
		zaplog.Error(newErr)
		return nil, newErr
	}
	var loginCredential = ""
	secret, err := sccapi.Decrypt(info.Password)
	if err != nil {
		newErr := fmt.Errorf("Sccapi decrypt failed in MQ init %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	user := url.QueryEscape(info.UserName)
	sec := url.QueryEscape(secret)
	loginCredential = user + ":" + sec + "@"

	url := "amqps://" + loginCredential + info.IP + ":" +
		strconv.Itoa(int(info.Port)) + "/"
	conn, err := rabbitmq.Open(url)
	if err != nil {
		newErr := fmt.Errorf("Connection to RabbitMQ failed, %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	return conn, nil
}

// Receive msgs by registering asynchronous call to Consume
func (mq *MQManager) ReceiveMsgs() (<-chan amqp.Delivery, error) {
	if mq.mqinfo.Mode == configmanager.MQProducer {
		newErr := fmt.Errorf("MQ Producer cant receive msg")
		zaplog.Error(newErr)
		return nil, newErr
	}
	return mq.mqClient.Consume(mq.mqinfo.QueueName, "")
}

// SendMsg
func (mq *MQManager) SendMsg(msg string) error {
	if mq.mqinfo.Mode == configmanager.MQConsumer {
		newErr := fmt.Errorf("MQ Consumer cant send msg")
		zaplog.Error(newErr)
		return newErr
	}
	mqMsg := &rabbitmq.MQMessage{
		Exchange:   mq.mqinfo.ExchangeName,
		RoutingKey: mq.mqinfo.RoutingKey,
		Body:       nil,
	}
	return mq.mqClient.Publish(mqMsg)
}
