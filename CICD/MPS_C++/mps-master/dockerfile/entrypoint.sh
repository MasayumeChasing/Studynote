#!/bin/bash                                                                                                                                                                                           
#set -e
umask 0077

wan_ip=`cat /opt/eip.txt`

cat >/opt/scc/scc.conf <<EOF
[CRYPTO]
primaryKeyStoreFile=/opt/scc/pk/primary.ks
standbyKeyStoreFile=/opt/scc/stb/standby.ks
backupFolderName=/opt/scc/bak_k
logCfgFile=/opt/scc/logger.conf
logCategory=SCC
domainCount=8
EOF

cat >/opt/scc/logger.conf <<EOF
[global]
file perms = 600
[formats]
default  = "%d[%-5V][%p:%t][%f:%U:%L]%m%n"
[rules]
SCC.INFO    "/opt/log/$app/scc/scc.log",5MB * 20 ~ "/opt/log/$app/scc/scc.log.#2r"; default
EOF

cat >/opt/$app/conf/config.yaml <<EOF
---
system:
  daemon: on

mps:
  #1-Ordinary service 2-Quality service
  node_type: 1 
  region: center
  #1:domain  2:ipv4  3:ipv6
  net_type: 2
  jwt_timeout: 120

scc: 
  # scc配置文件地址
  config-path: /opt/scc/scc.conf

log:
  path: /opt/log/mps/run/
  # debug info warn error nolog
  level: info
  # console file all
  mode: all
  file_keep_days: 2
  file_keep_num: 10
  #1:hls  2:rtsp
  protocol_type: 1
server:
  live:
    device_timeout: 30
    client_timeout: 30
    data_timeout: 11
modules:
  hls:
    wan_ip: ${wan_ip}
    lan_ip: ${listen_ip:-127.0.0.1}
    listen: 7081
  rtsp:
    wan_ip: ${wan_ip}
    lan_ip: ${listen_ip:-127.0.0.1}
    listen: 7082

cert_key:
  ivm_cert_path: /opt/certs/user_cert.pem
  local_key_path: /opt/certs/user_cert.key
  
unix:
  addr: /opt/holo/amah_agent/amah_agent.sock

apollo:
  appid: IVMHoloMediaProxyService
  cluster: ${apollo_cluster:-default}
  accesskey_secret: ${apollo_secret_mps:-xxx}
  namespace:
    - "global.jwt"
    - "global.register.center"
    - "global.certs"
    - "global.elb.inside"
    - "billing.rabbitmq"
    - "monitor.rabbitmq"
    - "collector.mps"
    - "application"
  cache_dir: ./
  meta_addr: ${apollo_addr:-https://127.0.0.1:8080}
EOF

cat >/opt/$app/conf/amah_agent.yaml <<EOF
---
log:
  level: ${log_level:-1} # 0: DEBUG 1: INFO 2: WARN 3: ERROR 4: CRITICAL
  path: ${log_path:-/opt/log/$app}

apollo:
  app_id: "service.paas.amah_agent"
  cluster: ${apollo_cluster:-default}
  default_namespace: "namespaces"
  cache_dir: ${apollo_cache:-/opt}
  meta_addr: ${apollo_addr:-https://127.0.0.1:8080}

startup_mode: online # offline

root_key:
  component1: "rY2KogQnKTFGzCD3edALUTjM1XqskW8qtcr78IoK1flbpRgyVy8p0fljWpqkKNpAN6bXCNxPJHlx8VKGWx0GseONsQC2yJQ0iZ6EQ3VgFL8Ude3CtluA6mpqRvaKdpCv"
  salt: "81YauAURApoFKXpjE9l0oeDBT7oLmAd6caV3dDzrr18AazDPihVnXG1xCe1YQ47hUTxooF9sPXWZLHBJtF7ZEf11Haj6ZHTWBSuabuJONZ7vuEtggkNM50WgnKe2vVau"
EOF

# 启动amah_agent
nohup /opt/$app/holo_amah_agent -c /opt/$app/conf/amah_agent.yaml &
ps -ef|grep amah_agent
echo "start amah_agent success."
sleep 1

# 判断amah_agent是否成功启动
ps -ef | grep -w "/opt/$app/holo_amah_agent" | grep -v "grep" >/dev/null 2>&1
if [ $? -eq 0 ];then
    echo "程序已启动"
fi

/opt/$app/holo_$app -c /opt/$app/conf/config.yaml

# tail -f /dev/null

sleep 5
echo "exit"
