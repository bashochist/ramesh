
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

package bufanalysis

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func printAsText(writer io.Writer, fileAnnotations []FileAnnotation) error {
	return printEachAnnotationOnNewLine(
		writer,
		fileAnnotations,
		printFileAnnotationAsText,
	)
}

func printAsMSVS(writer io.Writer, fileAnnotations []FileAnnotation) error {
	return printEachAnnotationOnNewLine(
		writer,
		fileAnnotations,
		printFileAnnotationAsMSVS,
	)
}

func printAsJSON(writer io.Writer, fileAnnotations []FileAnnotation) error {
	return printEachAnnotationOnNewLine(
		writer,
		fileAnnotations,
		printFileAnnotationAsJSON,
	)
}

func printAsJUnit(writer io.Writer, fileAnnotations []FileAnnotation) error {
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ")
	testsuites := xml.StartElement{Name: xml.Name{Local: "testsuites"}}
	err := encoder.EncodeToken(testsuites)
	if err != nil {
		return err
	}
	annotationsByPath := groupAnnotationsByPath(fileAnnotations)
	for _, annotations := range annotationsByPath {
		path := "<input>"
		if fileInfo := annotations[0].FileInfo(); fileInfo != nil {
			path = fileInfo.ExternalPath()
		}
		path = strings.TrimSuffix(path, ".proto")
		testsuite := xml.StartElement{
			Name: xml.Name{Local: "testsuite"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: path},
				{Name: xml.Name{Local: "tests"}, Value: strconv.Itoa(len(annotations))},
				{Name: xml.Name{Local: "failures"}, Value: strconv.Itoa(len(annotations))},
				{Name: xml.Name{Local: "errors"}, Value: "0"},
			},
		}
		if err := encoder.EncodeToken(testsuite); err != nil {
			return err
		}
		for _, annotation := range annotations {
			if err := printFileAnnotationAsJUnit(encoder, annotation); err != nil {
				return err
			}
		}
		if err := encoder.EncodeToken(xml.EndElement{Name: testsuite.Name}); err != nil {
			return err
		}
	}
	if err := encoder.EncodeToken(xml.EndElement{Name: testsuites.Name}); err != nil {
		return err
	}