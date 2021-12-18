#!/usr/bin/env bash

set -eo pipefail

USAGE="$(basename "$0") [-h] [-l lines]

where:
    -h  Show this help text
    -l  Number of lines to generate, defaults to 100000 lines
    "

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

FILE="${SCRIPT_DIR}/test-data.txt"
NUM_LINES=100000

#
# Start by reading command-line arguments using the getopts function
#
while getopts ":hl:" OPTION
do
  case ${OPTION} in
    h     ) echo "${USAGE}"; exit ;;
    l     ) NUM_LINES="${OPTARG}";;
    *     ) echo "Unknown option, terminating."; echo "${USAGE}"; exit 1;;
  esac
done


if [[ -f "${FILE}" ]]; then
  rm "${FILE}"
fi

for (( i=1; i<=${NUM_LINES}; i++ ))
do
 rand=$((1 + ${RANDOM} ))
 echo "http://api.tech.com/item/${i} ${rand}" >> "${FILE}"
done


exit 0
