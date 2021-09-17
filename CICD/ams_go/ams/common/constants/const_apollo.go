// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: const_appolo.go
// Description: Contians appolo config functions

package constants

// apollo 信息
const (
	// ************ namespace ************
	ApolloNsCommon   = "common"
	ApolloNsRabbitMQ = "rabbitmq"
)

const (
	defaultTimeout = 20
	// TimeoutRead 读超时
	TimeoutRead  = "read"
	timeoutWrite = "write"
	timeoutDial  = "dial"
	// TimeoutLock 锁超时
	TimeoutLock = "lock"
)

const (
	// apollo

	// KeyApolloSecret 访问秘钥
	KeyApolloSecret = "apollo.secret"
	// KeyApolloAppID app id
	KeyApolloAppID = "apollo.appId"
	// KeyApolloCluster 集群
	KeyApolloCluster = "apollo.cluster"
	// KeyApolloLocalCache 缓存目录
	KeyApolloLocalCache = "apollo.localcache"
	// KeyApolloIP 访问ip
	KeyApolloIP = "apollo.ip"

	KeyJwtServiceSecret    = "common.jwt.service.secret"
	KeyJwtServiceOldSecret = "common.jwt.service.old.secret"
	KeyJwtServiceExpire    = "common.jwt.service.expire"

	// mongodb
	KeyMongodbIP       = "common.mongodb.endpoint"
	KeyMongodbPort     = "common.mongodb.port"
	KeyMongodbUsername = "common.mongodb.userName"
	KeyMongodbPassword = "common.mongodb.passWord"

	// Server config
	KeyGRPCPort                   = "common.grpc.port"
	KeyHttpPort                   = "common.http.port"
	KeyServerForciblyExitWaitTime = "common.http.timeout"

	// Rabbit MQ input msg to AMS for search
	KeyMQInMsgIP         = "rabbitmq.in.endpoint"
	KeyMQInMsgPort       = "rabbitmq.in.port"
	KeyMQInMsgUsername   = "rabbitmq.in.username"
	KeyMQInMsgPassword   = "rabbitmq.in.password"
	KeyMQInMsgType       = "rabbitmq.in.type"
	KeyMQInMsgQueueName  = "rabbitmq.in.queuename"
	KeyMQInMsgBindingKey = "rabbitmq.in.bindingkey"
	KeyMQInMsgRoutingKey = "rabbitmq.in.routingkey"

	// Rabbit MQ ouput alarm msg raised by AMS
	KeyMQOutMsgIP         = "rabbitmq.out.endpoint"
	KeyMQOutMsgPort       = "rabbitmq.out.port"
	KeyMQOutMsgUsername   = "rabbitmq.out.username"
	KeyMQOutMsgPassword   = "rabbitmq.out.password"
	KeyMQOutMsgType       = "rabbitmq.out.type"
	KeyMQOutMsgQueueName  = "rabbitmq.out.queuename"
	KeyMQOutMsgBindingKey = "rabbitmq.out.bindingkey"
	KeyMQOutMsgRoutingKey = "rabbitmq.out.routingkey"

	// SLM server configuration
	KeySlmServerIP   = "common.slms.endpoint"
	KeySlmServerPort = "common.slms.port"

	// JWT info
	KeyAMSJwtSecret = "common.jwt.amssecret"
	KeySLMJwtSecret = "common.jwt.slmsecret"

	// ES Engine INfo
	KeyEngineEndpoints = "common.engine.endpoints"
	KeyEngineUsername = "common.engine.username"
	KeyEnginePassword = "common.engine.password"
)
