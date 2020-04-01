#!/usr/bin/env bash

cygwin_mount_dir=$(mount | grep "cygwin64 on" | awk '{print $1}') # fix for cygwin64 users
tmp_dir=${cygwin_mount_dir}$(mktemp -d -t ci-XXXXXXXXXX)

printf "CI Testing ...\n"
mkdir -p "${tmp_dir}/.out"
OUTPUT=$(./tests/env/start.sh && go test -v -count=1 ./tests/ && ./tests/env/stop.sh)
printf "%s\n" "$OUTPUT"
if echo "$OUTPUT" | grep -q "FAIL"; then
    echo 'Tests failed' >/dev/stderr
    exit
fi