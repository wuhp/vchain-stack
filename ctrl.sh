#!/bin/bash

usage() {
    echo "Usage: service vstack start/stop/restart"
    exit 1
}

start() {
    ps -ef | grep -v "grep" | grep "vstack" > /dev/null 2>&1
    [ $? -eq 0 ] && echo "Service vstack is already running" && return

    /usr/sbin/vstack -c /etc/vstack/vstack.yml
    [ $? -eq 0 ] && echo "Starting vstack [   OK   ]" || echo "Starting vstack [ Failed ]"
}

stop() {
    killall vstack > /dev/null 2>&1
    for x in `seq 1 15`
    do
        sleep 1

        if [ $x -eq 15 ]; then
            killall -9 vstack > /dev/null 2>&1
            break
        fi

        ps -ef | grep -v "grep" | grep "vstack" > /dev/null 2>&1
        [ $? -ne 0 ] && break
    done
    echo "Stopping vstack [   OK   ]"
}

restart() {
    stop
    start
}

[ $# -ne 1 ] && usage
$1
