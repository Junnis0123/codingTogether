#!/bin/bash

processid=`ps aux | grep main | grep -v "grep" | awk '{print $2}'`

kill -9 ${processid}

go build main.go
nohup ./main config_server server &

