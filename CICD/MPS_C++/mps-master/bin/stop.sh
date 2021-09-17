#!/bin/bash
process=holo_mps

pid=`ps -ef | grep $process| awk '{print $2}'`

if [ "-$pid" != "-" ]
then
    kill -9 $pid
else
    echo "mps already stoped."
fi
