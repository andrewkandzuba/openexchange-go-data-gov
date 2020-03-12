#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
if [ -x "$(command -v cygpath)" ]; then
  DIR="$(cygpath -w "${DIR}")"
fi
docker-compose -f "${DIR}"/docker-compose.yml down