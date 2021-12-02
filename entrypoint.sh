#!/bin/bash

set -e

echo $PWD
ls 

# Show user all variables for debugging
printf "repository: ${INPUT_REPOSITORY}\n"
printf "org: ${INPUT_ORG}\n"
printf "config: ${INPUT_CONFIG}\n"
printf "outfile: ${INPUT_OUTFILE}\n"
printf "metric: ${INPUT_METRIC}\n"
printf "skip: ${INPUT_SKIP}\n"
printf "pattern: ${INPUT_PATTERN}\n"


if [ -z "${INPUT_REPOSITORY}" ] && [ -z "${INPUT_ORG}" ]; then
    printf "You must define either repository or org.\n"
    exit
fi

# An org is defined
if [ -z "${INPUT_REPOSITORY}" ]; then
    COMMAND="codestats org ${INPUT_ORG}"
    OUTFILE="org.json"

    # org specific matching commands
    if [ ! -z "${INPUT_PATTERN}" ]; then
        COMMAND="${COMMAND} --pattern ${INPUT_PATTERN}"
    fi
    if [ ! -z "${INPUT_SKIP}" ]; then
        COMMAND="${COMMAND} --skip ${INPUT_SKIP}"
    fi

else
    COMMAND="codestats repo ${INPUT_REPOSITORY}"
    OUTFILE="repo.json"

fi

if [ ! -z "${INPUT_CONFIG}" ]; then
    COMMAND="${COMMAND} --config ${INPUT_CONFIG}"
fi

if [ ! -z "${INPUT_METRIC}" ]; then
    COMMAND="${COMMAND} --metric ${INPUT_METRIC}"
fi

if [ ! -z "${INPUT_OUTFILE}" ]; then
    OUTFILE="${INPUT_OUTFILE}"
fi

COMMAND="${COMMAND} --outfile ${OUTFILE}"
printf "$COMMAND\n"
$COMMAND
echo $?
