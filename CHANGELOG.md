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

## [v1.14.0] - 2023-02-09

- Replace `buf generate --include-types` with `buf generate --type` for consistency. `--include-types`
  is now deprecated but continues to work, consistent with our compability guarantee.
- Include type references in `google.protobuf.Any` messages in option values
  when filtering on type, e.g. with `buf build --type` or `buf generate --type`.
- Allow specifying a specific `protoc` path in `buf.gen.yaml` when using `protoc`'s built-in plugins
  via the new `protoc_path` option.
- Allow specifying arguments for local plugins in `buf.gen.yaml`. You can now do e.g.
  `path: ["go, "run", ./cmd/protoc-gen-foo]` in addition to `path: protoc-gen-foo`.
- Add optional name parameter to `buf mod init`, e.g. `buf mod init buf.build/owner/foobar`.
- Fix issue with `php_metadata_namespace` file option in [managed mode](https://docs.buf.build/generate/managed-mode).
- Make all help documentation much clearer. If you notice any inconsistencies, let us know.

## [v1.13.1] - 2023-01-27

- Fix race condition with `buf generate` when remote plugins from multiple
  BSR instances are being used at once.

## [v1.13.0] - 2023-01-26

- Extend the `BUF_TOKEN` environment variable to accept tokens for multiple
  BSR instances. Both `TOKEN` and `TOKEN1@BSRHOSTNAME1,TOKEN2@BSRHOSTNAME2,...`
  are now valid values for `BUF_TOKEN`.
- Remove `buf beta convert` in favor of the now-stable `buf convert`.

## [v1.12.0] - 2023-01-12
- Add `buf curl` command to invoke RPCs via [Connect](https://connect-build),
  [gRPC](https://grpc.io/), or [gRPC-Web](https://github.com/grpc/grpc-web.)
- Introduce `objc_class_prefix` option in managed mode, allowing a `default` value
  for `objc_class_prefix` for all files, `except` and `override`, which both behave
  similarly to other `except` and `override` options. Specifying an empty `default`
  value is equivalent to having managed mode on in previous versions.
- Introduce `ruby_package` option in managed mode, allowing `except` and `override`,
  in the same style as `objc_class_prefix`. Leaving `ruby_package` unspecified has
  the same effect as having mananged mode enabled in previous versions.

## [v1.11.0] - 2022-12-19
- `buf generate` now batches remote plugin generation calls for improved performance.
- Update `optimize_for` option in managed mode, allowing a `default` value for `optimize_for`
  for all files, `except` and `override`, which both behave similarly to other `except`
  and `override` options. Specifying an `optimize_for` value in the earlier versions is
  equivalent to having a `optimize_for` with that value as default.

## [v1.10.0] - 2022-12-07

- When using managed mode, setting `enabled: false` now no longer fails `buf generate`
  and instead prints a warning log and ignores managed mode options.
- Add `csharp_namespace` option to managed mode, allowing `except`, which excludes
  modules from managed mode, and `override`, which specifies `csharp_namespace` values
  per module, overriding the default value. By default, when managed mode is enabled,
  `csharp_namespace` is set to the package name with each package sub-name capitalized.
- Promote `buf convert` to stable, keep `buf beta convert` aliased in the beta command.
- Add `Types` 