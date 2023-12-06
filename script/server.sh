#!/usr/bin/env bash

set -e

CUR_DIR="$( cd "$( dirname "$0" )" && pwd )"
source "${CUR_DIR}/env.sh"

make run-server ARGS=""
