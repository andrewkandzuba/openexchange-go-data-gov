#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
if [ -x "$(command -v cygpath)" ]; then
  DIR="$(cygpath -w "${DIR}")"
fi

check_brokers="$(docker run -it --network=host edenhill/kafkacat:1.5.0 -b localhost:9092 -L | grep brokers)"
if [ "${check_brokers}" ]
then
  printf "%s\n" "Stopping kafka brokers..."
  "${DIR}"/stop.sh
fi

docker-compose -f "${DIR}"/docker-compose.yml up -d

i=0
while [ $i -lt 10 ]
do
	(( i++ ))

  check_brokers="$(docker run -it --network=host edenhill/kafkacat:1.5.0 -b localhost:9092 -L | grep brokers)"
	if [ -z "${check_brokers}" ]
  then
    printf "%s\n" "waiting on kafka brokers..."
    sleep 5
	  continue
	else
	  printf "%s\n" "kafka brokers are ready!"
	  break
	fi
done
