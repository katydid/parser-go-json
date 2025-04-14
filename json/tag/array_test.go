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

package tag

import (
	"io"
	"testing"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/parse"
)

func TestTagArrayForEmptyArray(t *testing.T) {
	s := `[]`
	// will be parsed the same as : {"array": []}
	p := NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), WithArrayTag())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect(t, p.Next, parse.ArrayOpenHint)
	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect(t, p.Next, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForNonEmptyArray(t *testing.T) {
	s := `["myelem"]`
	// will be parsed the same as : {"array": ["myelem"]}
	p := NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), WithArrayTag())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect(t, p.Next, parse.ArrayOpenHint)

	// in startState, see "myelem"
	expect(t, p.Next, parse.ValueHint)
	expectStr(t, p, "myelem")

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect(t, p.Next, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayWithEmptyArray(t *testing.T) {
	s := `[[]]`
	// will be parsed the same as : {"array": [{"array": []}]}
	p := NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), WithArrayTag())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect(t, p.Next, parse.ArrayOpenHint)

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect(t, p.Next, parse.ArrayOpenHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect(t, p.Next, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect(t, p.Next, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
