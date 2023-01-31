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

// Matching the unix-like build tags in the Golang source i.e. https://github.com/golang/go/blob/912f0750472dd4f674b69ca1616bfaf377af1805/src/os/file_unix.go#L6

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package netrc

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMachineForName(t *testing.T) {
	t.Parallel()
	testGetMachineForNameSuccess(
		t,
		"foo.com",
		"testdata/unix/home1",
		"foo.com",
		"bar",
		"baz",
	)
	testGetMachineForNameNil(
		t,
		"api.foo.com",
		"testdata/unix/home1",
	)
	testGetMachineForNameNil(
		t,
		"bar.com",
		"testdata/unix/home1",
	)
	testGetMachineForNameSuccess(
		t,
		"foo.com",
		"testdata/unix/home2",
		"",
		"bar",
		"baz",
	)
	testGetMachineForNameSuccess(
		t,
		"api.foo.com",
		"testdata/unix/home2",
		"",
		"bar",
		"baz",
	)
	testGetMachineForNameSuccess(
		t,
		"bar.com",
		"testdata/unix/home2",
		"",
		"bar",
		"baz",
	)
	testGetMachineForNameNil(
		t,
		"foo.com",
		"testdata/unix/home3",
	)
	testGetMachineForNameNil(
		t,
		"api.foo.com",
		"testdata/unix/home3",
	)
	testGetMachineForNameNil(
		t,
		"bar.com",
		"testdata/unix/home1",
	)
}

func TestPutMachines(t *testing.T) {
	t.Parallel()
	testPutMachinesSuccess(
		t,
		false,
		NewMachine(
			"foo.com",
			"test@foo.com",
			"password",
		),
	)
	testPutMachinesSuccess(
		t,
		true,
		NewMachine(
			"foo.com",
			"test@foo.com",
			"password",
		),
	)
	testPutMachinesSuccess(
		t,
		false,
		NewMachine(
			"bar.com",
			"test@bar.com",
			"password",
		),
		NewMachine(
			"baz.com",
			"test@baz.com",
			"password",
		),
	)
}

// https://github.com/bufbuild/buf/issues/611
func TestPutLotsOfBigMachinesSingleLineFiles(t *tes