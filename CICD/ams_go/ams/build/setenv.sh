#!/bin/bash

# Download dependencies (like github.com/golang/protobuf) from below URL
export GOPROXY=http://cmc.centralrepo.rnd.huawei.com/cbu-go/

# Disable checksum fetch for all repos
export GONOSUMDB=*

export GOSUM=off

export GO111MODULE=on
