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
- Add `buf beta studio-agent` command to support the upcoming Buf Studio.

## [v1.5.0] - 2022-05-30

- Upgrade to `protoc` 3.20.1 support.
- Fix an issue where `buf` would fail if two or more roots contained
  a file with the same name, but with different file types (i.e. a
  regular file vs. a directory).
- Fix check for `PACKAGE_SERVICE_NO_DELETE` to detect deleted services.
- Remove `buf beta registry track`.
- Remove `buf beta registry branch`.

## [v1.4.0] - 2022-04-21

- Fix issue where duplicate synthetic oneofs (such as with proto3 maps or
  optional fields) did not result in a properly formed error.
- Add `buf beta registry repository update` command which supports updating
  repository visibility (public vs private). As with all beta commands, this
  is likely to change in the future.

## [v1.3.1] - 2022-03-30

- Allow `--config` flag to be set when targeting a module within a workspace.
- Update `buf format`'s file option order so that default file options are
  sorted before custom options.
- Update `buf format` to write adjacent string literals across multiple lines.
- Fix `buf format` so that the output directory (if any) is created if and only
  if the input is successfully formatted.

## [v1.3.0] - 2022-03-25

- Add `--exit-code` flag to `buf format` to exit with a non-zero exit code if
  the files were not already formatted.

## [v1.2.1] - 2022-03-24

- Fix a few formatting edge cases.

## [v1.2.0] - 2022-03-24

- Add `buf format` command to format `.proto` files.
- Fix build scripts to avoid using the `command-line-arguments` pseudo-package
  when building binaries and re-introduce checking for proper usage of private
  packages.

## [v1.1.1] - 2022-03-21

- Remove check for proper usage of private packages due to a breaking change made in the Golang standard library in 1.18.

## [v1.1.0] - 2022-03-01
- Add `--type` flag to the `build` command to create filtered images containing
  only the specified types and their required dependencies.
- Trim spaces and new lines from user-supplied token for `buf registry login`.
- Add support for conversion between JSON and binary serialized message for `buf beta convert`.

## [v1.0.0] - 2022-02-17

- Check that the user provided a valid token when running `buf registry login`.
- Add `buf mod open` that opens a module's homepage in a browser.
- Add `buf completion` command to generate auto-completion scripts in commonly used shells.
- Add `--disable-symlinks` flag to the `breaking, build, export, generate, lint, ls-files, push`
  commands. By default, the CLI will follow symlinks except on Windows, and this disables following
  symlinks.
- Add `--include-wkt` flag to `buf generate`. When this flag is specified alongside
  `--include-imports`, this will result in the [Well-Known Types](https://github.com/bufbuild/wellknowntypes/tree/11ea259bf71c4d386131c1986ffe103cb1edb3d6/v3.19.4/google/protobuf)
  being generated as well. Most language runtimes have the Well-Known Types included as part
  of the core library, making generating the Well-Known Types separately undesirable.
- Remove `buf protoc`. This was a pre-v1.0 demonstration to show that `buf` compilation
  produces equivalent results to mainline `protoc`, however `buf` is working on building
  a better Protobuf future that provides easier mechanics than our former `protoc`-based
  world. `buf protoc` itself added no benefit over mainline `protoc` beyond being considerably
  faster and allowing parallel compilation. If `protoc` is required, move back to mainline `protoc`
  until you can upgrade to `buf`. See [#915](https://github.com/bufbuild/buf/pull/915) for more
  details.
- Context modifier no longer overrides an existing token on the context. This allows `buf registry login`
  to properly check the user provided token without the token being overridden by the CLI interceptor.
- Removed the `buf config init` command in favor of `buf mod init`.
- Removed the `buf config ls-breaking-rules` command in favor of `buf mod ls-breaking-rules`.
- Removed the `buf config ls-lint-rules` command in favor of `buf mod ls-lint-rules`.
- Removed the `buf config migrate-v1beta1` command in favor of `buf beta migrate-v1beta1`.
- Add `buf beta decode` command to decode message with provided image source and message type.
- Disable `--config` flag for workspaces.
- Move default config version from `v1beta1` to `v1`.

## [v1.0.0-rc12] - 2022-02-01

- Add `default`, `except` and `override` to `java_package_prefix`.
- Add dependency commits as a part of the `b3` digest.
- Upgrade to `protoc` 3.19.4 support.
- Remove `branch` field from `buf.lock`.

## [v1.0.0-rc11] - 2022-01-18

- Upgrade to `protoc` 3.19.3 support.
- Add `PACKAGE_NO_IMPORT_CYCLE` lint rule to detect package import cycles.
- Add `buf beta registry {plugin,template} {deprecate,undeprecate}`.
- Add warning when using enterprise dependencies without specifying a enterprise
  remote in the module's identity.
- Remove `digest`, and `created_at` fields from the `buf.lock`. This will temporarily create a new commit
  when pushing the same contents to an existing repository, since the `ModulePin` has been reduced down.
- Add `--track` flag to `buf push`
- Update `buf beta registry commit list` to allow a track to be specified.
- Add `buf beta registry track {list,delete}` commands.
- Add manpages for `buf`.

## [v1.0.0-rc10] - 2021-12-16

- Fix issue where remote references were not correctly cached.

## [v1.0.0-rc9] - 2021-12-15

- Always set `compiler_version` parameter in the `CodeGeneratorRequest` to "(unknown)".
- Fix issue where `buf mod update` was unable to resolve dependencies from different remotes.
- Display the user-provided Buf Schema Registry remote, if specified, instead of the default within the `buf login` message.
- Fix issue where `buf generate` fails when the same plugin was specified more than once in a single invocation.
- Update the digest algorithm so that it encodes the `name`, `lint`, and `breaking` configuration encoded in the `buf.yaml`.
  When this change is deployed, users will observe the following:
  - Users on `v0.43.0` or before will notice mismatched digest errors similar to the one described in https://github.com/bufbuild/buf/issues/661.
  - Users on `v0.44.0` or after will have their module cache invalidated, but it will repair itself automatically.
  - The `buf.lock` (across all versions) will reflect the new `b3-` digest values for new commits.

## [v1.0.0-rc8] - 2021-11-10

- Add new endpoints to the recommendation service to make it configurable.
- Add `--exclude-path` flag to `buf breaking`, `buf build`, `buf export`, `buf generate`, and `buf lint` commands. This allows users to exclude specific paths when running commands.
- Change `GetModulePackages` endpoint to return a repeated `ModulePackage` message that now includes package description with the package name.
- Add `Oneof` to the `Message` structure for documentation.

## [v1.0.0-rc7] - 2021-11-08

- Upgrade to `protoc` 3.19.1 support.
- Fix issue with `buf generate` where multiple insertion points are defined in the same file.

## [v1.0.0-rc6] - 2021-10-20

- Fix issue with `buf ls-files` when given an image as an input, imports were being printed,
  even without the `--include-imports` flag.
- Add the ability for users to provide individual protobuf files as inputs to CLI commands. This allows users to run `buf` commands against and file input based on their current working directory, for example, `buf lint foo/bar.proto`, where `foo/bar.proto` is a path to protobuf file on disk.

## [v1.0.0-rc5] - 2021-10-12

- Add `buf beta registry repository deprecate` and `buf beta registry repository undeprecate`.
- Support `--include-imports` for remote plugins.
- Fix issue where `buf config migrate-v1beta1 fails` when files cannot be renamed.
- Fix issue where `buf registry login` panics when an existing .netrc entry exists.

## [v1.0.0-rc4] - 2021-10-07

- Fix issue where `buf generate` could fail when used with large numbers of plugins and files on
  systems with low file limits.
- Add `buf protoc --version` flag back. This was accidentally removed.
- Upgrade to `protoc` 3.18.1 support.

## [v1.0.0-rc3] - 2021-10-04

- Add `--as-import-paths` flag to `ls-files` that strips local directory paths and prints file
  paths as they are imported.
- Fix issue where groups used in custom options did not result in the same behavior as `protoc`.
- Fix issue where insertion points were not applied with respect to the configured output directory.

## [v1.0.0-rc2] - 2021-09-23

- Add `--include-imports` flag to `ls-files`.
- Upgrade to `protoc` 3.18.0 support.
- Fix regression with git inputs using `recurse_submodules=true`.

## [v1.0.0-rc1] - 2021-09-15

This is our first v1.0 release candidate. This release largely concentrates on erroring for
already-deprecated commands and flags.

At Buf, we take compatibility very seriously. When we say v1.0, we mean it - we hope `buf` will be
stable on v1 for the next decade, and if there is something we want to change, it is our responsibility to
make sure that we don't break you, not your responsibility to change because of us. We have learned
a lot about `buf` usage in the last two years of our beta, and have deprecated flags and commands as
we go, but for v1.0, we are removing the deprecated items to make sure we have a clean setup going forward.

All commands and flags have been printing warnings for a long time, and have an easy migration path.
Simply update the command or flag, and you'll be good to go:

- Removed the `buf login` command in favor of `buf registry login`.
- Removed the `buf logout` command in favor of `buf registry logout`.
- Removed the `buf mod init` command in favor of `buf config init`.
- Removed the `--name` and `--dep` flags in `buf config init`.
- Removed the `--log-level` global flag.
- Moved the output of `--version` from stderr to stdout.
- Moved the output of `--help` and `help` from stderr to stdout.
- [From v0.55.0](https://github.com/bufbuild/buf/releases/tag/v0.55.0): The version key in all configuration files (`buf.yaml`, `buf.gen.yaml`, `buf.work.yaml`) is now required.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta config init` command in favor of `buf config init`.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta mod export` command in favor of `buf export`.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta mod init` command in favor of `buf config init`.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta mod update` command in favor of `buf mod update`.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta mod clear-cache` command in favor of `buf mod clear-cache`.
- [From v0.45.0](https://github.com/bufbuild/buf/releases/tag/v0.45.0): Removed the `buf beta push` command in favor of `buf push`.
- [From v0.34.0](https://github.com/bufbuild/buf/releases/tag/v0.34.0): Removed the `buf check breaking` command in favor of `buf breaking`.
- [From v0.34.0](https://github.com/bufbuild/buf/releases/tag/v0.34.0): Removed the `buf check lint` command in favor of `buf lint`.
- [From v0.34.0](https://github.com/bufbuild/buf/releases/tag/v0.34.0): Removed the `buf check ls-lint-checkers` command in favor of `buf config ls-lint-rules`.
- [From v0.34.0](https://github.com/bufbuild/buf/releases/tag/v0.34.0): Removed the `buf check ls-breaking-checkers` command in favor of `buf config ls-breaking-rules`.
- [From v0.31.0](https://github.com/bufbuild/buf/releases/tag/v0.31.0): Removed the `--file` flag on `buf build` in favor of the `--path` flag.
- [From v0.31.0](https://github.com/bufbuild/buf/releases/tag/v0.31.0): Removed the `--file` flag on `buf lint` in favor of the `--path` flag.
- [From v0.31.0](https://github.com/bufbuild/buf/releases/tag/v0.31.0): Removed the `--file` flag on `buf breaking` in favor of the `--path` flag.
- [From v0.31.0](https://github.com/bufbuild/buf/releases/tag/v0.31.0): Removed the `--file` flag on `buf generate` in favor of the `--path` flag.
- [From v0.31.0](https://github.com/bufbuild/buf/releases/tag/v0.31.0): Removed the `--file` flag on `buf export` in favor of the `--path` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--source` flag on `buf build` in favor of the first positional parameter.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--source-config` flag on `buf build` in favor of the `--config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input` flag on `buf lint` in favor of the first positional parameter.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input-config` flag on `buf lint` in favor of the `--config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input` flag on `buf breaking` in favor of the first positional parameter.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input-config` flag on `buf breaking` in favor of the `--config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--against-input` flag on `buf breaking` in favor of the `--against` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--against-input-config` flag on `buf breaking` in favor of the `--against-config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input` flag on `buf generate` in favor of the first positional parameter.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input-config` flag on `buf generate` in favor of the `--config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input` flag on `buf ls-files` in favor of the first positional parameter.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `--input-config` flag on `buf ls-files` in favor of the `--config` flag.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `buf image build` command in favor of `buf build`.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `buf image convert` command.
- [From v0.29.0](https://github.com/bufbuild/buf/releases/tag/v0.29.0): Removed the `buf beta image convert` command.
- [From v0.23.0](https://github.com/bufbuild/buf/releases/tag/v0.23.0): Removed the `buf experimental image convert` command.
- [From v0.52.0](https://github.com/bufbuild/buf/releases/tag/v0.52.0) [and v0.34.0](https://github.com/bufbuild/buf/releases/tag/v0.34.0): Complete deletion `protoc-gen-buf-check-breaking` and `protoc-gen-buf-check-lint`, which have been moved to `protoc-gen-buf-breaking` and `protoc-gen-buf-lint`.

In January 2021 (v0.34.0), `protoc-gen-buf-check-breaking` and `protoc-gen-buf-check-lint` were deprecated and scheduled for removal for v1.0. In August 2021 (v0.52.0), we began returning error for every invocation of `protoc-gen-buf-check-breaking` and `protoc-gen-buf-check-lint`. This release completes the deletion process.

The only migration necessary is to change your installation and invocation from `protoc-gen-buf-check-breaking` to `protoc-gen-buf-breaking` and `protoc-gen-buf-check-lint` to `protoc-gen-buf-lint`. These can be installed in the exact same manner, whether from GitHub Releases, Homebrew, AUR, or direct Go installation:

```
# instead of go get github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking
go get github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking
# instead of curl -sSL https://github.com/bufbuild/buf/releases/download/v0.57.0/protoc-gen-buf-check-breaking-Linux-x86_64
curl -sSL https://github.com/bufbuild/buf/releases/download/v0.57.0/protoc-gen-buf-breaking-Linux-x86_64
```

## [v0.56.0] - 2021-09-08

- Cascade `ENUM_ZERO_VALUE_SUFFIX` comment ignores from the enum level.
- Fix issue where `buf genarate --output` was not being respected in 0.55.0.

## [v0.55.0] - 2021-09-07

- Error if `version:` is not set in `buf.yaml`. This is one of the few breaking changes we must make before v1.0 to guarantee stability for the future. If you do not have a version set, simply add `version: v1beta1` to the top of your `buf.yaml`.
- Support `BUF_TOKEN` for authentication. `buf` will now look for a token in the `BUF_TOKEN` environment variable, falling back to `.netrc` as set via `buf login`.
- Add support for using remote plugins with local source files.
- Add per-file overrides for managed mode.
- Fix issue with the module cache where multiple simulataneous downloads would result in a temporarily-corrupted cache.
- Hide verbose messaing behind the `--verbose` (`-v`) flag.
- Add `--debug` flag to print out debug logging.

## [v0.54.1] - 2021-08-30

- Fix docker build.

## [v0.54.0] - 2021-08-30

- Add windows support.
- Add `java_package_prefix` support to managed mode.
- Fix issue with C# namespaces in managed mode.
- Fix issue where `:main` was appended for errors containing references to modules.

## [v0.53.0] - 2021-08-25

- Fix issue where `buf generate --include-imports` would end up generating files for certain imports twice.
- Error when both a `buf.mod` and `buf.yaml` are present. `buf.mod` was briefly used as the new default name for `buf.yaml`, but we've reverted back to `buf.yaml`.

## [v0.52.0] - 2021-08-19

Return error for all invocations of `protoc-gen-buf-check-breaking` and `protoc-gen-buf-check-lint`.

As one of the few changes buf will ever make, `protoc-gen-buf-check-breaking` and `protoc-gen-buf-check-lint` were deprecated and scheduled for removal for v1.0 in January 2021. In preparation for v1.0, instead of just printing out a message notifying users of this, these commands now return an error for every invocation and will be completely removed when v1.0 is released.

The only migration necessary is to change your installation and invocation from `protoc-gen-buf-check-breaking` to `protoc-gen-buf-breaking` and `protoc-gen-buf-check-lint` to `protoc-gen-buf-lint`. These can be installed in the exact same manner, whether from GitHub Releases, Homebrew, AUR, or direct Go installation:

```
# instead of go get github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking
go get github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking
# instead of curl -sSL https://github.com/bufbuild/buf/releases/download/v0.52.0/protoc-gen-buf-check-breaking-Linux-x86_64
curl -sSL https://github.com/bufbuild/buf/releases/download/v0.52.0/protoc-gen-buf-breaking-Linux-x86_64
```

There is no change in functionality.

## [v0.51.1] - 2021-08-16

- Fix issue with git LFS where a remote must be set for fetch.

## [v0.51.0] - 2021-08-13

- Accept packages of the form `v\d+alpha` and `v\d+beta` as packages with valid versions. These will be considered unstable packages for the purposes of linting and breaking change detection if `ignore_unstable_packages` is set.
- Fix issue with git clones that occurred when using a previous reference of the current branch.

## [v0.50.0] - 2021-08-12

- Add `buf generate --include-imports` that also generates all imports except for the Well-Known Types.
- Fix issue where a deleted file within an unstable package that contained messages, enums, or services resulted in a breaking change failure if the `PACKAGE` category was used and `ignore_unstable_packages` was set.

## [v0.49.0] - 2021-08-10

- Split `FIELD_SAME_TYPE` breaking change rule into `FIELD_SAME_TYPE, FIELD_WIRE_COMPATIBLE_TYPE, FIELD_WIRE_JSON_COMPATIBLE_TYPE` in `v1`. See https://github.com/bufbuild/buf/pull/400 for details.
- Only export imported dependencies from `buf export`.

## [v0.48.2] - 2021-07-30

- Fix git args for http auth with git lfs.

## [v0.48.1] - 2021-07-30

- Fix: use `-c` on `git` parent command instead of `--config` on `git fetch`.
- Add `ruby_package` to managed mode.

## [v0.48.0] - 2021-07-29

- Add `buf export`. `buf export` will export the files from the specified input (default `"."`) to the given directory in a manner that is buildable by `protoc` without any `-I` flags. It also has options `--exclude-imports`, which excludes imports (and won't result in a buildable set of files), and `--path`, which filters to the specific paths.

## [v0.47.0] - 2021-07-29

- Rewrite the git cloner to use `git init && git fetch` rather than `git clone`. `git clone` is limited to local branches on the remote, whereas `git fetch` we can fetch any references on the remote including remote branches.
- Add `php_namespace` managed mode handling.
- Add `java_string_check_utf8` managed mode handling.

## [v0.46.0] - 2021-07-27

- Add `buf login` and `buf logout` to login and logout from the Buf Schema Registry.
- Fix cache, configuration, and data environment variables for Windows. Note that while Windows is still not officially supported, `buf` largely works on Windows.

## [v0.45.0] - 2021-07-26

- Revert default configuration file location back from `buf.mod` to `buf.yaml`. Note that both continue to work.
- Move default workspace configuration file location from `buf.work` to `buf.work.yaml`. Note that both continue to work.
- Move `buf beta push` to `buf push`. Note that `buf beta push` continues to work.
- Move most `buf beta mod` commands to `buf mod`. Note that all `buf beta mod` commands continue to work.
- Add `--only` flag to `buf mod update`.
- Warn if `buf.yaml` contains dependencies that are not represented in the `buf.lock` file.
- Add `--version` flag to `buf config ls-{breaking,lint}-rules`.
- Add `SYNTAX_SPECIFIED` lint rule to `BASIC, DEFAULT` categories for v1 configuration.
- Add `IMPORT_USED` lint rule to `BASIC, DEFAULT` categories for v1 configuration.
- Bring v1 configuration out of beta.
- Add managed mode for `objc_class_prefix`, `csharp_namespace`.

## [v0.44.0] - 2021-07-08

- Fix issue where C++ scoping rules were not properly enforced.
- Add support for splitting directory paths passed to `buf protoc -I` by a directory separator.
- Fix Windows support for builtin `protoc` plugins when using `buf generate` or `buf protoc`. Note that Windows remains officially unsupported as we have not set up testing, but largely works.
- Upgrade to `protoc` 3.17.3 support.
- Change the default module configuration location from `buf.yaml` to `buf.mod`. Note that `buf.yaml` continues to work.
- Continued work on the workspaces beta, including the `v1` configuration specification.
- Continued work on the managed mode beta, including the `v1` configuration specification.
- Add `v1` module configuration specification in beta - please continue to use `v1beta1` until the `v1` configuration specification is rolled out.
- Add `buf config migrate-v1beta1`.

## [v0.43.2] - 2021-05-31

- Fix namespace resolution diff with protoc.

## [v0.43.1] - 2021-05-28

- Revert `protoc` namespace resolution diff change.

## [v0.43.0] - 2021-