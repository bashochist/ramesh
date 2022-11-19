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

package bandeps

import (
	"context"
	"sync"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type state struct {
	logger            *zap.Logger
	envStdioContainer app.EnvStdioContainer
	runner            command.Runner
	violationMap      map[string]Violation
	// map from ./foo/bar/... to actual packages
	packageExpressionToPackages     map[string]*packagesResult
	packageExpressionToPackagesLock *keyRWLock
	// map from packages to dependencies
	packageToDeps     map[string]*depsResult
	packageToDepsLock *keyRWLock
	lock              sync.RWMutex
	calls             int
	cacheHits         int
	tracer            trace.Tracer
}

func newState(
	logger *zap.Logger,
	envStdioContainer app.EnvStdioContainer,
	runner command.Runner,
) *state {
	return &state{
		logger:        