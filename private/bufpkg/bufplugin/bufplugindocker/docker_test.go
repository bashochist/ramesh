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

package bufplugindocker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	dockerVersion = "1.41"
)

var (
	imagePattern = regexp.MustCompile("^(?P<image>[^/]+/[^/]+/[^/:]+)(?::(?P<tag>[^/]+))?(?:/(?P<op>[^/]+))?$")
)

func TestPushSuccess(t *testing.T) {
	t.Parallel()
	server := newDockerServer(t, dockerVersion)
	listenerAddr := server.httpServer.Listener.Addr().String()
	dockerClient := createClient(t, WithHost("tcp://"+listenerAddr), WithVersion(dockerVersion))
	image, err := buildDockerPlugin(t, "testdata/success/Dockerfile", listenerAddr+"/library/go")
	require.Nilf(t, err, "failed to build docker plugin")
	require.NotEmpty(t, image)
	pushResponse, err := dockerClient.Push(context.Background(), image, &RegistryAuthConfig{})
	require.Nilf(t, err, "failed to push docker plugin")
	require.NotNil(t, pushResponse)
	assert.NotEmpty(t, pushResponse.Digest)
}

func TestPushError(t *testing.T) {
	t.Parallel()
	server := newDockerServer(t, dockerVersion)
	// Send back an error on ImagePush (still return 200 OK).
	server.pushErr = errors.New("failed to push image")
	listenerAddr := server.httpServer.Listener.Addr().String()
	dockerClient := createClient(t, WithHost("tcp://"+listenerAddr), WithVersion(dockerVersion))
	image, err := buildDockerPlugin(t, "testdata/success/Dockerfile", listenerAddr+"/library/go")
	require.Nilf(t, err, "failed to build docker plugin")
	require.NotEmpty(