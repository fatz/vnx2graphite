#!/bin/bash
# ./vnx2graphite.sh <statsname> <server_name> 
DBUG=""
if [ -n "$VNX2GRAPHITE_DEBUG" ];then
    DBUG="-v"
fi
export NAS_DB='/nas';/nas/bin/server_stats $2 -count 1 -interval 1 -terminationsummary no -format csv -monitor $1>/tmp/$1.vnx2graphite.txt 2>/dev/null && $HOME/vnx2graphite/vnx2graphite -c $HOME/vnx2graphite/vnx2graphite.conf -d /tmp/$1.vnx2graphite.txt -m $2 -s $1 $DBUG && rm /tmp/$1.vnx2graphite.txt && sleep 20