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

package parse

import (
	"io"
	"testing"

	"github.com/katydid/parser-go-json/json/internal/expect"
	"github.com/katydid/parser-go/parse"
)

func TestTagMixObject(t *testing.T) {
	s := `[{"mykey1":[{"mykey2":[]}]}]`
	// will be parsed the same as : {"array":[{"object":{"mykey":{"array":[{"object":{"mykey2":{"array":[]}}}]}}}]}
	p := NewParser(WithBuffer([]byte(s)), WithObjectTag(), WithArrayTag())

	// 1: first array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 2: first object
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey1")

	// 3: second array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 4: second object
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey2")

	// 5: third array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	// 4: second object
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	// 3: second array
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	// 2: first object
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	// 1: first array
	expect.Hint(t, p, parse.ArrayCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
