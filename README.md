# Vars to Credhub

This app will take a vars file suitable for use in [Concourse](https://concourse-ci.org)
pipelines or in [BOSH](https://bosh.io) commands.

## To build locally
1. Make sure you have `go>=1.10` installed.
1. Make sure you have `dep` installed.
1. Clone this repository.
1. From the repository directory, run `dep ensure`.
1. Run the tests: `go test ./...`
1. Build the app: `go build`

## Running
```shell
$ ./vars-to-credhub --prefix /concourse/main --vars-file vars.yml > bulk-import.yml
$ credhub bulk-import --file bulk-import.yml
```
