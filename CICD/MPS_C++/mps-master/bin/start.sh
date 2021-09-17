#!/bin/bash
BASE_DIR=$(cd "$(dirname -- $0)"; pwd)
PROJECT_DIR=${BASE_DIR}/..

CONFIG_PATH=${PROJECT_DIR}/config/config.yaml

export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${PROJECT_DIR}/lib
echo "start mps begin"
nohup ${PROJECT_DIR}/bin/holo_mps -c ${CONFIG_PATH} >> /dev/null 2>&1 &
#${PROJECT_DIR}/bin/holo_mps -c ${CONFIG_PATH}
echo "start mps end"
