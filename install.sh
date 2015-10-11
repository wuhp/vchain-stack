#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))

mkdir -p /etc/vstack
mkdir -p /var/log/vstack

install ${ROOT_DIR}/bin/vstack /usr/sbin/vstack
