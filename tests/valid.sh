#!/bin/sh

export PATH=".:${PATH}"

set -ex

vars-to-credhub --vars-file="fixtures/input.yml" --prefix="/plop" > result.yml
diff -u "fixtures/expected.yml" "result.yml"
