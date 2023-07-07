#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)

cd ${SCRIPT_DIR} && ./awsiot_media_device > /dev/null 2>&1
