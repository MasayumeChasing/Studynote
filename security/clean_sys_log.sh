#!/bin/bash
truncate -s 0 /var/log/yum.log
truncate -s 0 /var/log/messages
truncate -s 0 /var/log/sudo.log
truncate -s 0 /var/log/secure
truncate -s 0 /var/log/syslog
truncate -s 0 /var/log/lastlog
truncate -s 0 /var/log/cron
truncate -s 0 /var/log/daemon.log
truncate -s 0 /var/log/authlog
/bin/rm -f /var/log/*.gz
/bin/rm -f /var/log/tuned/*.gz
/bin/rm -f /var/log/aide/*.gz
/bin/rm -f /var/log/audit/*.gz
/usr/bin/rm -rf /usr/src/kernels/*


/usr/bin/rm -rf /home/holo_ops/*
/usr/bin/rm -rf /home/holo_ops/.ansible
/usr/bin/rm -rf /root/*
/usr/bin/rm -rf /root/.ansible

/usr/bin/rm -rf /usr/share/.history/*
> /root/.bash_history
> /home/devdeploy/.bash_history
> /home/opsadmin/.bash_history
> /home/devcloud/.bash_history
history -c
/bin/rm -rf /tmp/*
