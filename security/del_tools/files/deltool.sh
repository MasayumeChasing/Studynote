#【安全加固项】--make卸载
cnt0=` rpm -q make |grep -v not |wc -l`
if [ $cnt0 -eq 0 ] ;then
echo " cnt0 The Reinforcement items have been solved. set 0"
else
yum remove make -y
echo " cnt0 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--telnet卸载
cnt1=` rpm -q telnet |grep -v not |wc -l`
if [ $cnt1 -eq 0 ] ; then
echo " cnt1 The Reinforcement items have been solved. set 0"
else
yum remove telnet -y
echo " cnt1 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--gcc卸载
cnt2=` rpm -q gcc |grep -v not |wc -l`
if [ $cnt2 -eq 0 ] ; then
echo " cnt2 The Reinforcement items have been solved. set 0"
else
yum remove gcc -y
echo " cnt2 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--cpp卸载
cnt3=` rpm -q cpp |grep -v not |wc -l`
if [ $cnt3 -eq 0 ] ; then
echo " cnt3 The Reinforcement items have been solved. set 0"
else
yum remove cpp -y
echo " cnt3 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--gdb卸载
cnt4=` rpm -q gdb |grep -v not |wc -l`
if [ $cnt4 -eq 0 ] ; then
echo " cnt4 The Reinforcement items have been solved. set 0"
else
yum remove gdb -y
echo " cnt4 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--readelf卸载
cnt5=` rpm -q readelf |grep -v not |wc -l`
if [ $cnt5 -eq 0 ] ; then
echo " cnt5 The Reinforcement items have been solved. set 0"
else
yum remove readelf -y
echo " cnt5 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--ld卸载
cnt6=` rpm -q ld |grep -v not |wc -l`
if [ $cnt6 -eq 0 ] ; then
echo " cnt6 The Reinforcement items have been solved. set 0"
else
yum remove ld -y
echo " cnt6 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--strace卸载
cnt7=` rpm -q strace |grep -v not |wc -l`
if [ $cnt7 -eq 0 ] ; then
echo " cnt7 The Reinforcement items have been solved. set 0"
else
yum remove strace -y
echo " cnt7 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--gcc-c++卸载
cnt8=` rpm -q gcc-c++ |grep -v not |wc -l`
if [ $cnt8 -eq 0 ] ; then
echo " cnt8 The Reinforcement items have been solved. set 0"
else
yum remove gcc-c++ -y
echo " cnt8 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--gcc-go卸载
cnt9=` rpm -q gcc-go |grep -v not |wc -l`
if [ $cnt9 -eq 0 ] ; then
echo " cnt9 The Reinforcement items have been solved. set 0"
else
yum remove gcc-go -y
echo " cnt9 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--gcc-gfortran卸载
cnt10=` rpm -q gcc-gfortran |grep -v not |wc -l`
if [ $cnt10 -eq 0 ] ; then
echo " cnt10 The Reinforcement items have been solved. set 0"
else
yum remove gcc-gfortran -y
echo " cnt10 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--objdump卸载
cnt11=` rpm -q objdump |grep -v not |wc -l`
if [ $cnt11 -eq 0 ] ; then
echo " cnt11 The Reinforcement items have been solved. set 0"
else
yum remove objdump -y
echo " cnt11 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--卸载biosdevname
cnt12=`rpm -qa|grep biosdevname |wc -l`
if [ $cnt12 -eq 0 ] ; then
echo " cnt12 The Reinforcement items have been solved. set 0 "
else
yum remove biosdevname -y
echo " cnt12 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--netcat是否卸载完成
cnt13=`sudo rpm -qa|grep netcat |wc -l`
if [ $cnt13 -eq 0 ] ; then
echo " cnt13 The Reinforcement items have been solved. set 0 "
else
rpm -e netcat --nodeps
echo " cnt13 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--ftp、telnet-server、tcpdump是否卸载
cnt14=`rpm -qa | grep ftp | wc -l`
cnt14_1=`rpm -qa | grep telnet-server | wc -l`
cnt14_2=`rpm -qa | grep tcpdump | wc -l`
if [ $cnt14 -eq 0 ] && [ $cnt14_1 -eq 0 ] && [ $cnt14_2 -eq 0 ] ; then
echo " cnt14 The Reinforcement items have been solved. set 0 "
else
test $cnt14 = "1" && rpm -e ftp
test $cnt14_1 = "1" && rpm -e telnet-server
test $cnt14_2 = "1" && rpm -e tcpdump
echo " cnt14 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--卸载openldap-clients
cnt15=`rpm -qa|grep openldap-clients |wc -l`
if [ $cnt15 -eq 0 ] ; then
echo " cnt15 The Reinforcement items have been solved. set 0 "
else
yum remove openldap-clients -y
echo " cnt15 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi

#【安全加固项】--卸载sssd-ldap
cnt16=`rpm -qa|grep sssd-ldap |wc -l`
if [ $cnt16 -eq 0 ] ; then
echo " cnt16 The Reinforcement items have been solved. set 0 "
else
yum remove sssd-ldap -y
echo " cnt16 The Reinforcement items still exists！set 1"
i=`expr $i + 1`
fi



for i in tcpdump sniffer wireshark netcat winpcap gdb strace readelf cpp gcc dexdump mirror jdk objdump;do which $i 2>/dev/null |xargs rm -f ;done
