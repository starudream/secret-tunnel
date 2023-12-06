#!/usr/bin/env bash

set -e

CUR_DIR="$( cd "$( dirname "$0" )" && pwd )"
source "${CUR_DIR}/env.sh"

make run-client ARGS="--key 9559c307f11f48a88caf42ea5b7844a9"
