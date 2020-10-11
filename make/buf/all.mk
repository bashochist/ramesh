GO_ALL_REPO_PKGS := ./cmd/... ./private/...
GO_BINS := $(GO_BINS) \
	cmd/buf \
	cmd/protoc-gen-buf-breaking \
	cmd/protoc-gen-buf-lint \
	private/bufpkg/bufstyle/cmd/bufstyle \
	private/bufpkg/bufwkt/cmd/wkt-go-data \
	private/pkg/bandeps/cmd/bandeps \
	private/pkg/git/cmd/git-ls-files-unstaged \
	private/pkg/storage/cmd/ddiff \
	private/pkg/st