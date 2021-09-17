#!/bin/bash
for t in `df --local -P | awk {'if (NR!=1) print $6'}`
do 
#  find $t -perm /177 \( -name '*.pem' -o -name '*.key' -o -name '*.p12' -o -name '*.pfx' \) 2>/dev/null | head
  find $t  \( -name '*.pem' -o -name '*.key' -o -name '*.p12' -o -name '*.pfx' \) 2>/dev/null | head
done
