#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
if [ -x "$(command -v cygpath)" ]; then
  DIR="$(cygpath -w "${DIR}")"
fi
docker-compose -f "${DIR}"/docker-compose.yml up -d

i=0
while [ $i -lt 10 ]
do
	(( i++ ))

  check_brokers="$(kafkacat -b localhost:9092 -L | grep brokers)"
	if [ -z "${check_brokers}" ]
  then
    sleep 10
	  continue
	else
	  break
	fi
done
