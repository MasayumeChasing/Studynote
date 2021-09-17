#!/bin/bash                                                                                                                                                                                           
set -e
umask 0077
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
SCC.INFO    "/opt/log/ams/scc/scc.log",5MB * 20 ~ "/opt/log/ams/scc/scc.log.#2r"; default
EOF

cat >/opt/$app/conf/env.yaml <<EOF
apollo:
  addr: ${apollo_addr:-https://127.0.0.1:8080}
  cluster: ${apollo_cluster:-default}
  appId: ${apollo_appId:-service.paas.dfs}
  secret: ${apollo_secret}
  cache: ${apollo_cache:-/opt}

mongo:
  enableSsl: ${enable_ssl:-true}

listenIp: ${listen_ip:-0.0.0.0}
sccConfPath: ${scc_conf_path:-/opt/scc/scc.conf}
logPath: ${log_path:-/opt/log/$app}
EOF
/opt/$app/$app
