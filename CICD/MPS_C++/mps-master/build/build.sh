#!/bin/bash 

BUILD_DIR=$(cd $(dirname $0);pwd)
WORKSPACE=$(cd ${BUILD_DIR}/..;pwd)

module=mps
module_sym=mps_sym
output_path=${WORKSPACE}/output

function build() {
    cd ${WORKSPACE}/deps
    sh ./gen_deps.sh
    cd ${WORKSPACE}
    #当前绝对路径
    cmake .
    clean
    make -j4
}

function pack() {
    build

    mkdir -p ${output_path}/${module}/config
    mkdir -p ${output_path}/${module}/lib
    mkdir -p ${output_path}/${module}/bin
    mkdir -p ${output_path}/${module_sym}

    cp -rf ${WORKSPACE}/build/dog                 ${output_path}/${module}
    cp -rf ${WORKSPACE}/bin/holo_mps              ${output_path}/${module}
    cp -rf ${WORKSPACE}/bin/start.sh              ${output_path}/${module}/bin
    cp -rf ${WORKSPACE}/bin/get_version.sh        ${output_path}/${module}/bin
    cp -rf ${WORKSPACE}/bin/libmps_common.so      ${output_path}/${module}/lib
    cp -rf ${WORKSPACE}/config/config.yaml        ${output_path}/${module}/config
    cp -rf ${WORKSPACE}/config/config.yaml.tpl    ${output_path}/${module}/config
    cp -rf ${WORKSPACE}/config/logger.conf        ${output_path}/${module}/config
    cp -rf ${WORKSPACE}/config/scc.conf           ${output_path}/${module}/config

    echo -e "build_time:`date --rfc-3339=seconds`\nbuild_branch:${CID_GLOBAL_REPO_BRANCH}\nbuild_commit:`git rev-parse HEAD`\nbuild_user:${CID_BUILD_USER}" > ${output_path}/${module}/build.version
    # 去符号
    objcopy --only-keep-debug ${output_path}/${module}/holo_mps                    ${output_path}/${module_sym}/holo_mps.sym
    objcopy --strip-all ${output_path}/${module}/holo_mps

    objcopy --only-keep-debug ${output_path}/${module}/lib/libmps_common.so        ${output_path}/${module_sym}/libmps_common.so.sym
    objcopy --strip-all ${output_path}/${module}/lib/libmps_common.so

    echo "...build ${module} success"
    cd ${output_path}
    tar -cvzf ${module}.tar.gz ${module}
    tar -cvzf ${module_sym}.tar.gz ${module_sym}
    cp -rf ${module}.tar.gz ${WORKSPACE}
    echo "${module}.tar.gz"
}

function help() {
    echo "$0 build|pack"
}

function clean() {
    rm -rf ${WORKSPACE}/${module}.tar.gz
    rm -rf ${output_path}/*
    rm -rf ${WORKSPACE}/build/target
    rm -f ${WORKSPACE}/bin/holo_mps ${WORKSPACE}/bin/libmps_common.so
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "pack" ];then
    pack
else
    help
fi

