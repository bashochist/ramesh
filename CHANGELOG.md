# Changelog

## [Unreleased]

- No changes yet.

## [v1.15.1] - 2023-03-08

- Fix a bug in `buf generate` with `v1beta1` config files.
- Fix a potential crash when using the `--type` flag with `buf build` or `buf generate`.

## [v1.15.0] - 2023-02-28

- Update built-in Well-Known Types to Protobuf v22.0.
- Fixes a bug in `buf format` where C-style block comments in which every
  line includes a prefix (usually "*") would be incorrectly indented.
- Add `--private-network` flag to `buf beta studio-agent` to support handling CORS requests
  from Studio on private networks that set the `Access-Control-Request-Private-Network` header.

##