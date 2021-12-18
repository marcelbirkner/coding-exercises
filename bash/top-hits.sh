#!/usr/bin/env bash

set -eo pipefail

USAGE="$(basename "$0") [-r RESULTS] FILE

where:
    -h  Show this help text
    -r  Number of results to print (Default 10 lines)
    "

#SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) # fetches path of script, but not needed here so commented

# Default number of results
RESULTS=10

# Try to determine number of processors / cores, for parallellizing the `sort` command.
CORES=$(nproc 2>/dev/null || echo "")
if [[ "x${CORES}" = "x" ]]; then
  echo -e "Number of cores for parallel execution not detected. Running as single thread.\n"
  CORES=1
fi

#
# Start by reading command-line arguments using the getopts function
#
while getopts ":hr:" OPTION
do
  case ${OPTION} in
    h     ) echo "${USAGE}"; exit ;;
    r	  ) RESULTS="${OPTARG}";;
    *     ) echo "Unknown option, terminating."; echo "${USAGE}"; exit 1;;
  esac
done

# Shift all arguments so to get the FILE as last argument
shift $((OPTIND - 1))
TEST_DATA_FILE="$1"

if [[ "x${TEST_DATA_FILE}" = "x" ]]; then
  echo -e "No test data file provided.\n"; echo "${USAGE}"
  exit 1
fi

if [[ ! -f "${TEST_DATA_FILE}" ]]; then
  echo "Test data file '${TEST_DATA_FILE}' does not exist."
  exit 1
fi


# First sort results based on the correct column. Then take the requested number of top results. Finally just fetch the
# first column from the resulting lines.
sort --parallel=${CORES} -n -r -k 2,2 "${TEST_DATA_FILE}" | head -n "${RESULTS}" | awk '{print $1}'
#      ^^^^^^^^                       Run the sort command with parallel execution to speed up
#                         ^           sort based on numerical value, not text
#                            ^        revert sorting, so highest results first
#                               ^ ^^^ sort on the second column containing the numeric value (<start>,<end> column)


exit 0
