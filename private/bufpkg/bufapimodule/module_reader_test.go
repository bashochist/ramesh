// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufapimodule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bufbuild/buf/private/bufpkg/bufmanifest"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/gen/proto/connect/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	modulev1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/module/v1alpha1"
	registryv1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	"github.com/bufbuild/buf/private/pkg/manifest"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	testDownload(
		t,
		"does-not-exist error",
		true,
		newMockDownloadService(
			t,
			withError(connect.NewError(connect.CodeNotFound, nil)),
		),
		"does not exist",
	)
	testDownload(
		t,
		"unexpected download service error",
		true,
		newMockDownloadService(
			t,
			withError(errors.New("internal")),
		),
		"internal",
	)
	testDownload(
		t,
		"success but response has all empty fields",
		true,
		newMockDownloadService(t),
		"expected non-nil manifest with tamper proofing enabled",
	)
	testDownload(
		t,
		"success with empty manifest module",
		true,
		newMockDownloadService(
			t,
			withBlobsFromMap(map[string][]byte{}),
		),
		"",
	)
	testDownload(
		t,
		"manifest module with invalid lock file",
		true,
		newMockDownloadService(
			t,
			withBlobsFromMap(map[string][]byte{
				"buf.lock": []byte("invalid lock file"),
			}),
		),
		"failed to decode lock file",
	)
	testDownload(
		t,
		"tamper proofing enabled no manifest",
		true,
		newMockDownloadService(
			t,
			withModule(&modulev1alpha1.Module{
				Files: []*modulev1alpha1.ModuleFile{
					{
						Path: "foo.proto",
					},
				},
			}),
		),
		"expected non-nil manifest with tamper proofing enabled",
	)
	testDownload(
		t,
		"tamper proofing disabled",
		false,
		newMockDownloadService(
			t,
			withModule(&modulev1alpha1.Module{
				Files: []*modulev1alpha1.ModuleFile{
					{
						Path: "foo.proto",
					},
				},
			}),
		),
		"",
	)
	testDownload(
		t,
		"tamper proofing disabled no image",
		false,
		newMockDownloadService(t),
		"no module in response",
	)
}

func testDownload(
	t *testing.T,
	desc string,
	tamperProofingEnabled bool,
	mock *mockDownloadService,
	erro