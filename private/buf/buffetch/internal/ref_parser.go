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

package internal

import (
	"context"
	"strconv"

	"github.com/bufbuild/buf/private/buf/bufref"
	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/git"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"go.uber.org/zap"
)

type refParser struct {
	logger                *zap.Logger
	rawRefProcessor       func(*RawRef) error
	singleFormatToInfo    map[string]*singleFormatInfo
	archiveFormatToInfo   map[string]*archiveFormatInfo
	dirFormatToInfo       map[string]*dirFormatInfo
	gitFormatToInfo       map[string]*gitFormatInfo
	moduleFormatToInfo    map[string]*moduleFormatInfo
	protoFileFormatToInfo map[string]*protoFileFormatInfo
}

func newRefParser(logger *zap.Logger, options ...RefParserOption) *refParser {
	refParser := &refParser{
		logger:                logger,
		singleFormatToInfo:    make(map[string]*singleFormatInfo),
		archiveFormatToInfo:   make(map[string]*archiveFormatInfo),
		dirFormatToInfo:       make(map[string]*dirFormatInfo),
		gitFormatToInfo:       make(map[string]*gitFormatInfo),
		moduleFormatToInfo:    make(map[string]*moduleFormatInfo),
		protoFileFormatToInfo: make(map[string]*protoFileFormatInfo),
	}
	for _, option := range options {
		option(refParser)
	}
	return refParser
}

func (a *refParser) GetParsedRef(
	ctx context.Context,
	value string,
	options ...GetParsedRefOption,
) (ParsedRef, error) {
	getParsedRefOptions := newGetParsedRefOptions()
	for _, option := range options {
		option(getParsedRefOptions)
	}
	return a.getParsedRef(ctx, value, getParsedRefOptions.allowedFormats)
}

func (a *refParser) getParsedRef(
	ctx context.Context,
	value string,
	allowedFormats map[string]struct{},
) (ParsedRef, error) {
	rawRef, err := a.getRawRef(value)
	if err != nil {
		return nil, err
	}
	singleFormatInfo, singleOK := a.singleFormatToInfo[rawRef.Format]
	archiveFormatInfo, archiveOK := a.archiveFormatToInfo[rawRef.Format]
	_, dirOK := a.dirFormatToInfo[rawRef.Format]
	_, gitOK := a.gitFormatToInfo[rawRef.Format]
	_, moduleOK := a.moduleFormatToInfo[rawRef.Format]
	_, protoFileOK := a.protoFileFormatToInfo[rawRef.Format]
	if !(singleOK || archiveOK || dirOK || gitOK || moduleOK || protoFileOK) {
		return nil, NewFormatUnknownError(rawRef.Format)
	}
	if len(allowedFormats) > 0 {
		if _, ok := allowedFormats[rawRef.Format]; !ok {
			return nil, NewFormatNotAllowedError(rawRef.Format, allowedFormats)
		}
	}
	if singleOK {
		return getSingleRef(rawRef, singleFormatInfo.defaultCompressionType)
	}
	if archiveOK {
		return getArchiveRef(rawRef, archiveFormatInfo.archiveType, archiveFormatInfo.defaultCompressionType)
	}
	if protoFileOK {
		return getProtoFileRef(rawRef), nil
	}
	if dirOK {
		return getDirRef(rawRef)
	}
	if gitOK {
		return getGitRef(rawRef)
	}
	if moduleOK {
		return getModuleRef(rawRef)
	}
	return nil, NewFormatUnknownError(rawRef.Format)
}

// validated per rules on rawRef
func (a *refParser) getRawRef(value string) (*RawRef, error) {
	// path is never empty after returning from this function
	path, options, err := bufref.GetRawPathAndOptions(value)
	if err != nil {
		return nil, err
	}
	rawRef := &RawRef{
		Path: path,
	}
	if a.rawRefProcessor != nil {
		if err := a.rawRefProcessor(rawRef); err != nil {
			return nil, er