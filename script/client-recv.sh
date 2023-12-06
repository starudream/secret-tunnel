#!/usr/bin/env bash

set -e

CUR_DIR="$( cd "$( dirname "$0" )" && pwd )"
source "${CUR_DIR}/env.sh"

export APP__TASK1__NAME="test"
export APP__TASK1__ADDRESS="127.0.0.1:8081"
export APP__TASK1__SECRET="3f6f694825e749ac92ced6f6688f9d9c"

make run-client ARGS="--key 6df6dc28bf804d18acc2d43e9f429e3c"
