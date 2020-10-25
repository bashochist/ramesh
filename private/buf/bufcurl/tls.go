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

package bufcurl

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bufbuild/buf/private/pkg/verbose"
)

// TLSSettings contains settings related to creating a TLS client.
type TLSSettings struct {
	// Filenames for a private key, certificate, and CA certificate pool.
	KeyFile, CertFile, CACertFile string
	// Override server name, for SNI.
	ServerName string
	// If true, the server's certificate is not verified.
	Insecure bool
}

// MakeVerboseTLSConfig constructs a *tls.Config that logs information to the
// given printer as a TLS connection is negotiated.
func MakeVerboseTLSConfig(settings *TLSSettings, authority string, printer verbose.Printer) (*tls.Config, error) {
	var conf tls.Config
	// we verify manually so that we can emit verbose output while doing so
	conf.InsecureSkipVerify = true
	conf.VerifyConnection = func(state tls.ConnectionState) error {
		printer.Printf("* TLS connection using %s / %s", versionName(state.Version), tls.CipherSuiteName(state.CipherSuite))
		if state.DidResume {
			printer.Printf("* (TLS session resumed)")
		}
		if st