/* Copyright (c) Huawei Technologies Co., Ltd. 2019-2020. All rights reserved. */

// Package constants for constants
package constants

import "time"

const (
	// JWTToken jwt token 的键名
	JWTToken = "Authorization"
	// JWTUserId jwt token体解析出来user_id
	JWTUserId = "user_id"
	// MaxConnectRetries  最大网络请求重试次数
	MaxConnectRetries = 3
	// OnErrorRetryInterval  请求重试间隔
	OnErrorRetryInterval = 1 * time.Second // 2s
	// HTTPConnectDefaultTimeOut 默认 connect timeout
	HTTPConnectDefaultTimeOut = time.Duration(3) * time.Second
	// HTTPIODefaultTimeOut 默认IO timeout
	HTTPIODefaultTimeOut = time.Duration(10) * time.Second
)

const (
	UacGetPicPath = "%s://%s/getPicPath?alarm_uuid=%s&user_id=%s"
)

// 下载文件接口 target参数
const (
	// 文件用于内部接口访问
	TargetInner = "inner"
	// 文件用于app访问
	TargetApp = "app"
	// 文件用于北向访问
	TargetNorth = "north"
)

// 下载文件接口 type参数
const (
	TypeImage = "image"
)

// 下载文件接口 source参数
const (
	SourceUAC = "uac"
)
