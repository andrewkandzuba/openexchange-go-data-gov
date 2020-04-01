#!/usr/bin/env bash

cygwin_mount_dir=$(mount | grep "cygwin64 on" | awk '{print $1}') # fix for cygwin64 users
tmp_dir=${cygwin_mount_dir}$(mktemp -d -t ci-XXXXXXXXXX)

printf "Building ...\n"
OUTPUT=$(go build ./...)
if echo "$OUTPUT" | grep -q "FAIL"; then
    echo 'Build failed' >/dev/stderr
    printf "%s\\n" "$OUTPUT"
    exit
fi

# @ToDo: Replace default test coverage tool with https://github.com/grosser/go-testcov
printf "Testing ...\n"
mkdir -p "${tmp_dir}/.out"
OUTPUT=$(go test -v -test.short -coverprofile=${tmp_dir}/.out/c.out -covermode=atomic ./...)
if echo "$OUTPUT" | grep -q "FAIL"; then
    echo 'Tests failed' >/dev/stderr
    printf "%s\\n" "$OUTPUT"
    exit
fi

go tool cover -html=${tmp_dir}/.out/c.out -o ${tmp_dir}/.out/coverage.html
printf "%s\n" "Test coverage results has been persisted into: file://${tmp_dir}/.out/coverage.html"