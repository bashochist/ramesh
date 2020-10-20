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
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
)

// headersBlockList contains disallowed headers. These are headers that are part
// of the Connect or gRPC protocol and set by the protocol implementations, so
// should not be set otherwise. It also includes "transfer-encoding", which is
// not part of either protocol, but is unsafe for users to set as it handled
// by the user agent.
//
// In addition to these headers, header names that start with "Connect-" and
// "Grpc-" are also reserved for use by protocol implementations.
var headerBlockList = map[string]struct{}{
	"accept":            {},
	"accept-encoding":   {},
	"content-type":      {},
	"content-encoding":  {},
	"te":                {},
	"transfer-encoding": {},
}

// GetAuthority determines the authority for a request with the given URL and
// request headers. If headers include a "Host" header, that is used. (If the
// request contains more than one, that is usually not valid or acceptable to
// servers, but this function will look at only the first.) If there is no
// such header, the authority is the host portion of the URL (both the domain
// name/IP address and port).
func GetAuthority(url *url.URL, headers http.Header) string {
	header := headers.Get("host")
	if header != "" {
		return header
	}
	return url.Host
}

// LoadHeaders computes the set of request headers from the given flag values,
// loading from file(s) if so instructed. A header flag is