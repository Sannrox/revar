#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/../..
source "${ROOT_PATH}/scripts/lib/init.sh"
source "${ROOT_PATH}/scripts/lib/test.sh"
source "${ROOT_PATH}/test/cmd/init.sh"

function create_revar_command() {
    make -C "${ROOT_PATH}" revar
}


create_revar_command
runTests
