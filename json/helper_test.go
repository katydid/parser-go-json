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

package json

import (
	"testing"

	"github.com/katydid/parser-go/parser/debug"
)

func parseJSON(s string) (debug.Nodes, error) {
	parser := NewJsonParser()
	if err := parser.Init([]byte(s)); err != nil {
		return nil, err
	}
	return debug.Parse(parser)
}

func testWalk(t *testing.T, s string) {
	t.Run(s, func(t *testing.T) {
		m, err := parseJSON(s)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Logf("%v", m)
	})
}

func testValue(t *testing.T, input, output string) {
	t.Run(input, func(t *testing.T) {
		jout, err := parseJSON(input)
		if err != nil {
			t.Fatalf("walk error: %v", err)
		}
		if len(jout) != 1 {
			t.Fatalf("expected one node")
		}
		if len(jout[0].Children) != 0 {
			t.Fatalf("did not expected any children")
		}
		if jout[0].Label != output {
			t.Fatalf("expected %q got %q", output, jout[0].Label)
		}
	})
}

func testSame(t *testing.T, input string) {
	testValue(t, input, input)
}

func testError(t *testing.T, s string) {
	t.Run("ExpectError"+s, func(t *testing.T) {
		parsed, err := parseJSON(s)
		if err != nil {
			t.Logf("PASS: given <%s> error: %v", s, err)
			return
		}
		t.Fatalf("FAIL: expected error given: <%v> got: %v", s, parsed)
	})
}
