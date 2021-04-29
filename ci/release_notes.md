This project is forked from [vmware-tanzu-labs/vars-to-credhub](https://github.com/vmware-tanzu-labs/vars-to-credhub), with several improvements and bug fixes.

### Bug fixes

- Fix issue with the `json` type being always used after some array value was encountered in the input YAML, instead of properly using the `value` type for plain values
- Properly use the `json` type for arrays, and only for such values
- Handle errors and properly exit with non-zero statuses when an error occurs
- More precise heuristic test for certificate values

### New features

- Add support for the `--version` flag, with `-v` for short
- Add `-h` alias for the `--help` flag
- Keep the order of YAML variables in the output array, just as they appear in the input YAML file, making it possible to predict the output document

### Project updates

- Switch to using go modules
- Simplify code
- Formalize smoke tests
- Add CI/CD pipeline, based on Stark & Wayne template, with Gstack improvements
- Leverage Github workflow to trigger Concourse, using the Gstack Github Action
