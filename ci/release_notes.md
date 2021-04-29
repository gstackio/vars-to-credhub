This project os forked from [vmware-tanzu-labs/vars-to-credhub](https://github.com/vmware-tanzu-labs/vars-to-credhub), with bug fixes

### Bug fixes

- Fix issue with `json` type always used after some array value was encountered, instead of using the `value` type
- Properly use the `json` type for arrays, and only for such values
- Handle errors and properly exit with non-zero statuses when an error occurs
- More precise heuristic test for certificate values

### New features

- Add support for the `--version` flag, with `-v` for short
- Add `-h` alias for the `--help` flag
- Keep the order of YAML variables from the input file in the output array, making it possible to predict the order of variables in output
- Formalized smoke tests

### Project updates

- Use go modules
- Simplify code
- Added CI/CD pipeline, based on Stark & Wayne, with Gstack improvements
- Leverage Github workflow to trigger Concourse, using the Gstack Github Action
