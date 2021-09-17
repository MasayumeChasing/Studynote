#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"

build_time=`cat ${DIR}/build.version|grep "build_time"|awk -F'time:' '{print $2}'|awk -F'+' '{print $1}'`;
time1=`date -d "${build_time}" +%s`;
time2=`date -d @${time1} +"%Y%m%d%H%M%S"`;

build_branch=`cat ${DIR}/build.version|grep "build_branch"|awk -F'branch:' '{print $2}'`;

build_commit=`cat ${DIR}/build.version|grep "build_commit"|awk -F'commit:' '{print $2}'`;
commit=`echo ${build_commit:0:6}`;

echo $time2.$build_branch.$commit > ${DIR}/tmp_version
