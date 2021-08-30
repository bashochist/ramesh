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

package studioagent

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/bufbuild/buf/private/buf/bufcli"
	"github.com/bufbuild/buf/private/bufpkg/bufstudioagent"
	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/app/appflag"
	"github.com/bufbuild/buf/private/pkg/cert/certclient"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"github.com/bufbuild/buf/private/pkg/transport/http/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	bindFlagName              = "bind"
	portFlagName              = "port"
	originFlagName            = "origin"
	disallowedHeadersFlagName = "disallowed-header"
	forwardHeadersFlagName    = "forward-header"
	caCertFlagName            = "ca-cert"
	clientCertFlagName        = "client-cert"
	clientKeyFlagName         = "client-key"
	serverCertFlagName        = "server-cert"
	serverKeyFlagName         = "server-key"
	privateNetworkFlagName    = "private-network"
)

// NewCommand returns a new Command.
func NewCommand(
	name string,
	builder appflag.Builder,
) *appcmd.Command {
	flags := newFlags()
	return &appcmd.Command{
		Use:   name,
		Short: "Run an HTTP(S) server as the Studio agent",
		Args:  cobra.ExactArgs(0),
		Run: builder.NewRunFunc(
			func(ctx context.Context, container appflag.Container) error {
				return run(ctx, container, flags)
			},
			bufcli.NewErrorInterceptor(),
		),
		BindFlags: flags.Bind,
	}
}

type flags struct {
	BindAddress       string
	Port              string
	Origin            string
	DisallowedHeaders []string
	ForwardHeaders    map[string]string
	CACert            string
	ClientCert        string
	ClientKey         string
	ServerCert