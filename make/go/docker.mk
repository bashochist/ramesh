# Managed by makego. DO NOT EDIT.

# Must be set
$(call _assert_var,MAKEGO)
$(call _conditional_include,$(MAKEGO)/base.mk)
# Must be set
$(call _assert_var,PROJECT)
# Must be set
$(call _assert_var,GO_MODULE)
# Must be set
$(call _assert_var,DOCKER_ORG)
# Must be set
$(call _assert_var,DOCKER_PROJECT)

DOCKER_WORKSPACE_IMAGE := $(DOCKER_ORG)/$(DOCKER_PROJECT)-workspace
DOCKER_WORKSPACE_FILE := Dockerfile.workspace
DOCKER_WORKSPACE_DIR := /workspace

# Settable
DOCKER_BINS ?=
# Settable
DOCKER_BUILD_EXTRA_FLAGS ?=

# Runtime
DOCKERMAKETARGET ?= all

.PHONY: dockerbuildworkspace
dockerbuildworkspace:
	docker build \
		$(DOCKER_BUILD_EXTRA_FLAGS) \
		--build-arg PROJECT=$(PROJECT) \
		--build-arg GO_MODULE=$(GO_MODULE) \
		-t $(DOCKER_WORKSPACE_IMAGE) \
		-f $(DOCKER_WORKSPACE_FILE) \
		.

.PHONY: dockermakeworkspace
dockermakeworkspace: dockerbuildworkspace
	docker run -v "$(CURDIR):$(DOCKER_WORKSPACE_DIR)" $(DOCKER_WORKSPACE_IMAGE) make -j 8 $(DOCKERMAKETARGET)

.PHONY: dockerbuild
dockerbuild::

define dockerbinfunc
.PHONY: dockerbuilddeps$(1)
dockerbuilddeps$(1)::

.PHONY: dockerbuild$(1)
dockerbuild$(1): dockerbuilddeps$(1)
	docker build $(DOCKER_BUILD_EXT