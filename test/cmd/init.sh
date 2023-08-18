#!/usr/bin/env bash
#
set -o errexit
set -o nounset
set -o pipefail
set -x


TEST_GO_FILES_DIR="${TEST_DATA_DIR}/go_files"


function runTests() {
    test_help
    test_single_file
    test_single_file_dry_run
    test_single_file_interactive
    test_dir_with_single_file
    test_recursive_dir_with_single_files
    test_recursive_dir_with_single_files_dry_run

}

 function test_help(){
     exspected_output="revar is a tool to replace variables in files"
     test::if_has_string "$("${OUTPUT_HOSTBINPATH}"/revar --help)" "${exspected_output}"
 }

function test_single_file() {
    local temp_dir=$(test::create_temp_dir)
    local file="${temp_dir}/clearscreen.go"
    local expected="${temp_dir}/clearscreen_expected.go"

    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen.go" "${file}"
    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen_expected.go" "${expected}"

    tr -d "\n" < "${file}" > "${file}.tmp"
    mv "${file}.tmp" "${file}"

    tr -d "\n" < "${expected}" > "${expected}.tmp"
    mv "${expected}.tmp" "${expected}"

    cmd=( "${OUTPUT_HOSTBINPATH}"/revar  "fmt" "test" "${file}")
    test::run_command "${cmd[@]}"


    test::no_diff_assert "${expected}" "${file}" "Test single file for variable replacement"
    rm -rf "${temp_dir}"
 }

 function test_single_file_dry_run() {
    local temp_dir=$(test::create_temp_dir)
    local file="${temp_dir}/clearscreen.go"
    local expected="${temp_dir}/clearscreen_expected.go"

    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen.go" "${file}"
    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen_expected.go" "${expected}"

     tr -d "\n" < "${file}" > "${file}.tmp"
     mv "${file}.tmp" "${file}"

     tr -d "\n" < "${expected}" > "${expected}.tmp"
     mv "${expected}.tmp" "${expected}"

    cmd=( "${OUTPUT_HOSTBINPATH}"/revar  "fmt" "test" "${file}" "-n")
    test::run_command "${cmd[@]}"

    test::diff_assert "${expected}" "${file}" "Test single file with dry-run for variable replacement - no replacement should happen"

    rm -rf "${temp_dir}"
}

function test_single_file_interactive() {
    local temp_dir=$(test::create_temp_dir)
    local file="${temp_dir}/clearscreen.go"
    local expected="${temp_dir}/clearscreen_expected.go"

    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen.go" "${file}"
    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen_expected.go" "${expected}"

     tr -d "\n" < "${file}" > "${file}.tmp"
     mv "${file}.tmp" "${file}"

     tr -d "\n" < "${expected}" > "${expected}.tmp"
     mv "${expected}.tmp" "${expected}"

    cmd=( "${OUTPUT_HOSTBINPATH}"/revar  "fmt" "test" "${file}" "-i")
    test::run_command_with_input "y\ny\ny\ny\n" "${cmd[@]}"

    test::no_diff_assert "${expected}" "${file}" "Test single file with interactive mode for variable replacement"

    rm -rf "${temp_dir}"
}

 function test_dir_with_single_file() {
     local temp_dir=$(test::create_temp_dir)
     local file="${temp_dir}/clearscreen.go"
     local expected="${temp_dir}/clearscreen_expected.go"

     cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen.go" "${file}"
     cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen_expected.go" "${expected}"

     tr -d "\n" < "${file}" > "${file}.tmp"
     mv "${file}.tmp" "${file}"

     tr -d "\n" < "${expected}" > "${expected}.tmp"
     mv "${expected}.tmp" "${expected}"

     cmd=( "${OUTPUT_HOSTBINPATH}"/revar  "fmt" "test" "${temp_dir}")
     test::run_command "${cmd[@]}"

     test::no_diff_assert "${expected}" "${file}" "Test single file with given dir for variable replacement"

  }

 function test_tes_dir_with_single_file_dry_run() {
    local temp_dir=$(test::create_temp_dir)
    local file="${temp_dir}/clearscreen.go"
    local expected="${temp_dir}/clearscreen_expected.go"

    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen.go" "${file}"
    cp "${TEST_GO_FILES_DIR}/clearscreen/clearscreen_expected.go" "${expected}"

    cmd=( "${OUTPUT_HOSTBINPATH}"/revar  "fmt" "test" "${file}" "-n")
    test::run_command "${cmd[@]}"

    test::diff_assert "${expected}" "${file}" "Test single file with given dir with dry-run for variable replacement - no replacement should happen"

    rm -rf "${temp_dir}"
}

function test_recursive_dir_with_single_files() {
    local temp_dir=$(test::create_temp_dir)

    cp -r "${TEST_GO_FILES_DIR}" "${temp_dir}"

    for file in $(find "${temp_dir}" -type f -name "*.go"); do
        tr -d "\n" < "${file}" > "${file}.tmp"
        mv "${file}.tmp" "${file}"
    done


    cmd=( "${OUTPUT_HOSTBINPATH}"/revar -r  "fmt" "test" "${temp_dir}")
    test::run_command "${cmd[@]}"


    for expected in $(find "${temp_dir}" -type f -name "*_expected.go"); do
        local actual="${expected%_expected.go}.go"
        test::no_diff_assert "${expected}" "${actual}" "Testing recursive file - ${actual}"
    done
}

function test_recursive_dir_with_single_files_dry_run() {
    local temp_dir=$(test::create_temp_dir)

    cp -r "${TEST_GO_FILES_DIR}" "${temp_dir}"

    for file in $(find "${temp_dir}" -type f -name "*.go"); do
        tr -d "\n" < "${file}" > "${file}.tmp"
        mv "${file}.tmp" "${file}"
    done


    cmd=( "${OUTPUT_HOSTBINPATH}"/revar -r  "fmt" "test" "${temp_dir}" "-n")
    test::run_command "${cmd[@]}"


    for expected in $(find "${temp_dir}" -type f -name "*_expected.go"); do
        local actual="${expected%_expected.go}.go"
        test::diff_assert "${expected}" "${actual}" "Testing recursive file with dry run - ${actual}"
    done
}

