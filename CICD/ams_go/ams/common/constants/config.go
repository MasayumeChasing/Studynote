// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: config.go
// Description: Contians config functions

package constants

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type apolloConfig struct {
	Addr    string `yaml:"addr"`
	Cluster string `yaml:"cluster"`
	AppID   string `yaml:"appId"`
	Secret  string `yaml:"secret"`
	Cache   string `yaml:"cache"`
}

type mongoOpt struct {
	EnableSsl bool `yaml:"enableSsl"`
}

type Config struct {
	ApolloConfig apolloConfig `yaml:"apollo,flow"`
	ListenIP     string       `yaml:"listenIp"`
	SccConfPath  string       `yaml:"sccConfPath"`
	LogPath      string       `yaml:"logPath"`
	Mongo        mongoOpt     `yaml:"mongo"`
}

var EnvConfig Config

const EnvConfigFile = "conf/env.yaml"
const ForciblyExitWaitTime = 10 * time.Second

// Load Env Config Files
func LoadEnvConfig() error {
	fileData, err := ioutil.ReadFile(EnvConfigFile)
	if err != nil {
		newErr := fmt.Errorf("ENV config file read failed")
		zaplog.Error(newErr)
		return newErr
	}
	err = yaml.Unmarshal(fileData, &EnvConfig)
	if err != nil {
		newErr := fmt.Errorf("Env config file parsing failed")
		zaplog.Error(newErr)
		return newErr
	}
	zaplog.Info("Loaded Env Config")
	return nil
}
