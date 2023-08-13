#!/usr/bin/env bash


set -o errexit
set -o nounset
set -o pipefail


unset CDPATH

ROOT_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
OUTPUT_SUBPATH=${OUTPUT_SUBPATH:-"_output"}
OUTPUT_PATH="${ROOT_PATH}/${OUTPUT_SUBPATH}"
OUTPUT_BINPATH="${OUTPUT_PATH}/bin"
GO_MODULE_URL=$( grep module < go.mod | cut -d " " -f2)


export OUTPUT_PATH
export OUTPUT_BINPATH
export GO_MODULE_URL





source "${ROOT_PATH}/scripts/lib/golang.sh"
source "${ROOT_PATH}/scripts/lib/version.sh"
