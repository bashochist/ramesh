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
- Add `Types` filter to `buf generate` command to specify types (message, enum,
  service) that should be included in the image. When specified, the resulting
  image will only include descriptors to describe the requested types.

## [v1.9.0] - 2022-10-19

- New compiler that is faster and uses less memory than the outgoing one.
  - When generating source code info, the new compiler is 20% faster, and allocates
    13% less memory.
  - If _not_ generating source code info, the new compiler is 50% faster and
    allocates 35% less memory.
  - In addition to allocating less memory through the course of a compilation, the
    new compiler releases some memory much earlier, allowing it to be garbage
    collected much sooner. This means that by the end of a very large compilation
    process, less than half as much memory is live/pinned to the heap, decreasing
    overall memory pressure.

  The new compiler also addresses a few bugs where Buf would accept proto sources
  that protoc would reject:
  - In proto3 files, field and enum names undergo a validation that they are
    sufficiently different so that there will be no conflicts in JSON names.
  - Fully-qualified names of elements (like a message, enum, or service) may not
    conflict with package names.
  - A oneof or extend block may not contain empty statements.
  - Package names may not be >= 512 characters in length or contain > 100 dots.
  - Nesting depth of messages may not be > 32.
  - Field types and method input/output types may not refer to synthetic
    map entry messages.
- Push lint and breaking configuration to the registry.
- Include `LICENSE` file in the module on `buf push`.
- Formatter better edits/preserves whitespace around inline comments.
- Formatter correctly indents multi-line block (C-style) comments.
- Formatter now indents trailing comments at the end of an indented block body
  (including contents of message and array literals and elements in compact options)
  the same as the rest of the body (instead of out one level, like the closing
  punctuation).
- Formatter uses a compact, single-line representation for array and message literals
  in option values that are sufficiently simple (single scalar element or field).
- `buf beta convert` flags have changed from `--input` to `--from` and `--output`/`-o` to `--to`
- fully qualified type names now must be parsed to the `input` argument and `--type` flag separately

## [v1.8.0] - 2022-09-14

- Change default for `--origin` flag of `buf beta studio-agent` to `https://studio.buf.build`
- Change default for `--timeout` flag of `buf beta studio-agent` to `0` (no timeout). Before it was
  `2m` (the default for all the other `buf` commands).
- Add support for experimental code generation with the `plugin:` key in `buf.gen.yaml`.
- Preserve single quotes with `buf format`.
- Support `junit` format errors with `--error-format`.

## [v1.7.0] - 2022-06-27

- Support protocol and encoding client options based on content-type in Studio Agent.
- Add `--draft` flag to `buf push`.
- Add `buf beta registry draft {list,delete}` commands.

## [v1.6.0] - 2022-06-21

- Fix issue where `// buf:lint:ignore` comment ignores did not work for the
  `ENUM_FIRST_VALUE_ZERO` rule.
- Add `buf beta studio-agent` command to support the up