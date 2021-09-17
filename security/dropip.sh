#!/bin/bash
ipt="/usr/sbin/iptables"
expired_date=`date -d "-30 minute" +%H:%M`

block() {
/usr/bin/netstat -na|grep ESTABLISHED|awk '{print $5}'|awk -F: '{print $1}'|sort|uniq -c|sort -rn|head -10|grep -v -E '192.168|127.0'|awk '{if ($2!=null && $1>100) {print $2}}' >/tmp/dropip
for i in $(cat /tmp/dropip)
do
  $ipt -I INPUT -s $i -j DROP
  echo "$i kill at `date +"%F %T"`">>/tmp/badip.txt
done 
#清除30分钟前的IP
expired_ip=`awk '/'$expired_date'/{print $1}' /tmp/badip.txt`
$ipt -D INPUT -s $expired_ip -j DROP
}
block
