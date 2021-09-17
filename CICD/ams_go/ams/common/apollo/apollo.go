/* Copyright (c) Huawei Technologies Co., Ltd. 2019-2020. All rights reserved. */

// Package apollo for cloud platform connect apollo
package apollo

import (
	ap "codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/middleware/apollo"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"github.com/philchia/agollo/v4"
	"github.com/pkg/errors"
    "codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
)

type ApolloReader struct {
	cli ap.Reader
}

// Init apollo初始化
func NewApolloReader() (*ApolloReader, error) {
    var apr ApolloReader
	var nss = []string {
		constants.ApolloNsCommon,
		constants.ApolloNsRabbitMQ,
	}

	if len(constants.EnvConfig.ApolloConfig.Secret) == 0 {
		return nil, errors.New("Parameter missing for apollo connection")
	}
	secret, err := sccapi.Decrypt(constants.EnvConfig.ApolloConfig.Secret)
	if err != nil {
		return nil, errors.Wrapf(err, "scc decrypt")
	}

	apr.cli, err = ap.Init(&agollo.Conf{
		AppID:              constants.EnvConfig.ApolloConfig.AppID,
		Cluster:            constants.EnvConfig.ApolloConfig.Cluster,
		NameSpaceNames:     nss,
		CacheDir:           constants.EnvConfig.ApolloConfig.Cache,
		MetaAddr:           constants.EnvConfig.ApolloConfig.Addr,
		AccesskeySecret:    secret,
		InsecureSkipVerify: true,
	}, &logger{})
	if err != nil {
		return nil, errors.Wrap(err, "apollo InitWithConfig")
	}

	return &apr, nil
}

// AppolloReader
func (apr *ApolloReader) Read(ns, k string) string {
	return apr.cli.Read(ns, k)
}

// Add new Listener
func (apr *ApolloReader) AddListener(ns, k string, cb func()) {
	apr.cli.AddListener(ns, k, cb)
}

// Add new NSListener
func (apr *ApolloReader) AddNSListener(ns string, cb func()) {
	apr.cli.AddNSListener(ns, cb)
}

// close
func (apr *ApolloReader) Close() error {
    // TODO Need to do
	return nil
}
