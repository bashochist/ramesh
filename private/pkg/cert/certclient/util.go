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

package certclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// NewClientTLSConfigFromRootCertFiles creates a new tls.Config from a root certificate files.
func NewClientTLSConfigFromRootCertFiles(rootCertFilePaths ...string) (*tls.Config, error) {
	rootCertDatas := make([][]byte, len(rootCertFilePaths))
	f