system:
  daemon: on

mps:
  # 1-Ordinary service 2-Quality service
  node_type: 1 
  region: center
  # 1:domain  2:ipv4  3:ipv6
  net_type: 2
  jwt_timeout: 120
  #1:hls  2:rtsp
  protocol_type: 1

scc: 
  config-path: /opt/holo/mps/config/scc.conf

log:
  path: $LOG_PATH/mps/log/run/
  # debug info warn error nolog
  level: info
  # console file
  mode: file
  # log file size, unit:MB
  size: 100
  file_keep_days: 7
  file_keep_num: 10

bandwidth:
  collection_interval: 5
  report_interval: 60

server:
  live:
    device_timeout: 30
    client_timeout: 30
    data_timeout: 11

modules:
  hls:
    wan_ip: $NETWORK_WANADDR
    lan_ip: $NETWORK_LANADDR
    listen: 7081

  rtsp:
    wan_ip: $NETWORK_WANADDR
    lan_ip: $NETWORK_LANADDR
    listen: 7082

cert_key:
  ivm_cert_path: /opt/ops/certs/media_cert.pem
  local_key_path: /opt/ops/certs/media_cert.key

unix:
  addr: /opt/holo/amah_agent/amah_agent.sock

apollo:
  appid: IVMHoloMediaProxyService
  cluster: $APOLLO_CLUSTER
  namespace:
#    - "global.jwt.interface"
#    - "paas.jwt.stream"
    - "global.jwt"
    - "global.register.center"
    - "global.certs"
    - "global.elb.inside"
    - "billing.rabbitmq"
    - "monitor.rabbitmq"
    - "collector.mps"
    - "application"
  cache_dir: ./
  meta_addr: $APOLLO_ADDR:$APOLLO_PORT
  accesskey_secret: ${APOLLO_SECRET_MPS:-xxx}
