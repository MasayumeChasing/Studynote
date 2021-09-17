#!/bin/bash
ipt="/usr/sbin/iptables"
log_file="/log/holo/nginx/access.log"
#当前时间减一分钟时间
d1=`date -d "-1 minute" +%H:%M`
#当前时间的分钟段
d2=`date +%M`
ips="/tmp/ips.txt"

block() {
    grep "$d1:" $log_file |awk '{print $1}' |sort -n |uniq -c |sort -n > $ips
    for ip in `awk '$1 >100 {print $2}' $ips`
    do
        $ipt -I INPUT -p -tcp --dport 9010 -s $ip -j REJECT
        echo "`date +"%F %T"` $ip" >> /tmp/badip.txt
    done 
}
unblock() {
  #将流量小于15的规则索引过滤出来
  for i in `$ipt -nvL --line-number |awk '$2 < 15 {print $1}' |sort -nr`
  do 
  #通过索引来删除规则
     $ipt -D INPUT $i
  done
  $ipt -Z
}

#每半小时清除一次
if [ $d2 == "00" ] || [ $d2 == "30" ]
then
  unblock
  block
else
  block
fi
