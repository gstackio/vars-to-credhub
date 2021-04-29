#!/bin/sh

export PATH=".:${PATH}"

set -x
vars-to-credhub --vars-file="fixtures/input_invalid.yml" 2> /dev/null
error=$?
set +x

if [ "${error}" -eq 0 ]; then
    echo >&2 "Expected non-zero exit status, but got '${error}'"
    exit 1
fi
