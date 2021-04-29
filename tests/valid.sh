#!/bin/sh

set -ex

./vars-to-credhub --vars-file="fixtures/input.yml" --prefix="/plop" > result.yml
spruce diff "fixtures/expected.yml" "result.yml"
