//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package json_test

import (
	"testing"

	sjson "github.com/katydid/parser-go-json/json"
)

func testValue(t *testing.T, input, output string) {
	t.Helper()
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(input)); err != nil {
		t.Errorf("init error: %v", err)
		return
	}
	jout, err := walk(parser)
	if err != nil {
		t.Errorf("walk error: %v", err)
		return
	}
	if len(jout) != 1 {
		t.Errorf("expected one node")
		return
	}
	if len(jout[0].Children) != 0 {
		t.Errorf("did not expected any children")
		return
	}
	if jout[0].Label != output {
		t.Errorf("expected %q got %q", output, jout[0].Label)
	}
}

func testSame(t *testing.T, input string) {
	t.Helper()
	testValue(t, input, input)
}

func testError(t *testing.T, s string) {
	t.Helper()
	parser := sjson.NewJsonParser()
	if err := parser.Init([]byte(s)); err != nil {
		t.Logf("PASS: given <%s> error: %v", s, err)
		return
	}
	parsed, err := walk(parser)
	if err != nil {
		t.Logf("PASS: given <%s> error: %v", s, err)
		return
	}
	t.Errorf("FAIL: expected error given: <%v> got: %v", s, parsed)
}
