#!/usr/bin/env bash
#
set -o errexit
set -o nounset
set -o pipefail

TEST_DATA_DIR="${ROOT_PATH}/test/data"

reset=$(tput sgr0)
bold=$(tput bold)
black=$(tput setaf 0)
red=$(tput setaf 1)
green=$(tput setaf 2)

function test::get_caller(){
    local levels :=${1:-2}
    local caller_file="${BASH_SOURCE[${levels}]}"
    local caller_line="${BASH_LINENO[${levels}]}"
    echo "$(basename "${caller_file}"):${caller_line}"
}

function test::run_command() {
    command=("$@")
    output=$("${command[@]}" 2>&1) || {
        cat <<EOF >&2
${bold}${red}Failed to run command ${command[@]} ${reset}${red}

${output}

EOF
return 1
}
}

function test::run_command_with_input() {
    local input="${1:-}"
    shift
    command=("$@")
    output=$(echo "${input}" | "${command[@]}" 2>&1) || {
        cat <<EOF >&2
${bold}${red}Failed to run command ${command[@]} ${reset}${red}

${output}

EOF
return 1
}

}
function test::no_diff_assert(){
    local expected="${1}"
    local actual="${2}"
    local message="${3:-}"

    if [[ ! -f "${expected}" ]]; then
        echo "File not found: ${expected}"
        exit 1
    fi

    if [[ ! -f "${actual}" ]]; then
        echo "File not found: ${actual}"
        exit 1
    fi

    if ! diff  -B -q "${expected}" "${actual}" >/dev/null; then
        echo "${bold}${red}File assertion failed: ${message}${reset}${red}"
        echo "Expected: ${expected}"
        echo "Actual: ${actual}"
        exit 1
    fi

    echo -n "${green}"
    echo "Successful"
    echo -n "${reset}"
    echo "message:${message}"
    echo "expected:${expected}"
    echo "actual:${actual}"
}

function test::diff_assert(){
    local expected="${1}"
    local actual="${2}"
    local message="${3:-}"

    if [[ ! -f "${expected}" ]]; then
        echo "File not found: ${expected}"
        exit 1
    fi

    if [[ ! -f "${actual}" ]]; then
        echo "File not found: ${actual}"
        exit 1
    fi

    if ! diff -B  -q "${expected}" "${actual}" >/dev/null; then
        echo -n "${green}"
        echo "Successful"
        echo -n "${reset}"
        echo "message:${message}"
        echo "expected:${expected}"
        echo "actual:${actual}"
    else
        echo "${bold}${red}diff assertion failed: ${message}${reset}${red}"
        echo "Expected: ${expected}"
        echo "Actual: ${actual}"
        exit 1
    fi
}



function test::if_has_string() {
  local message=$1
  local match=$2

  if grep -q "${match}" <<< "${message}"; then
    echo -n "${green}"
    echo "Successful"
    echo -n "${reset}"
    echo "message:${message}"
    echo "has:${match}"
    return 0
  else
    echo -n "${bold}${red}"
    echo "FAIL!"
    echo -n "${reset}"
    echo "message:${message}"
    echo "has not:${match}"
    caller
    return 1
  fi
}

function test::create_temp_dir() {
  local dir=$(mktemp -d)
  echo "${dir}"
}

function test::create_temp_file() {
  local file=$(mktemp)
  echo "${file}"
}

function test::create_dir_in_temp_dir() {
  local dir=$(test::create_temp_dir)
  mkdir -p "${dir}/$1"
  echo "${dir}/$1"
}

function test::create_temp_file_with_content() {
  local content="${1:-}"
  local file=$(test::create_temp_file)
  echo "${content}" > "${file}"
  echo "${file}"
}

function test::create_temp_file_with_content_from_file() {
  local content_file="${1:-}"
  local file=$(test::create_temp_file)
  cat "${content_file}" > "${file}"
  echo "${file}"
}

function test::create_temp_file_with_content_in_tempdir_with_dir(){
local dir = $(test::create_dir_in_temp_dir "$1")
local content="${2:-}"
local file=$(test::create_temp_file)
echo "${content}" > "${dir}/${file}"
echo "${dir}/${file}"
}

function test::create_temp_file_with_content_in_tempdir_with_file(){
local dir = $(test::create_dir_in_temp_dir "$1")
local content_file="${2:-}"
local file=$(test::create_temp_file)
cat "${content_file}" > "${dir}/${file}"
echo "${dir}/${file}"
}

