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

package tag_test

import (
	"io"
	"testing"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/tag"
	"github.com/katydid/parser-go/parse"
)

func TestTagObjectForEmptyObject(t *testing.T) {
	s := `{}`
	// will be parsed the same as : {"object": {}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), tag.WithObjectTag())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)
	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagObjectForNonEmptyObject(t *testing.T) {
	s := `{"mykey": "myvalue"}`
	// will be parsed the same as : {"object": {"mykey": "myvalue"}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), tag.WithObjectTag())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "mykey")

	// in startState, see "myvalue"
	expect(t, p.Next, parse.ValueHint)
	expectStr(t, p, "myvalue")

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagObjectForNonEmptyObjectWithEmptyObjectValue(t *testing.T) {
	s := `{"mykey": {}}`
	// will be parsed the same as : {"object": {"mykey": {"object": {}}}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), tag.WithObjectTag())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "mykey")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in objectTagCloseState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagObjectForNonEmptyObjectWithNonEmptyObjectValue(t *testing.T) {
	s := `{"mykey": {"mykey2": {}}}`
	// will be parsed the same as : {"object": {"mykey": {"object": {"mykey2": {"object": {}}}}}}
	p := tag.NewTagger(jsonparse.NewParser(jsonparse.WithBuffer([]byte(s))), tag.WithObjectTag())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "mykey")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect(t, p.Next, parse.KeyHint)
	expectStr(t, p, "mykey2")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyState and return fake key "object"
	expect(t, p.Next, parse.KeyHint)
	expectTag(t, p, "object")

	// in objectTagKeyState, go to startState return real "{"
	expect(t, p.Next, parse.ObjectOpenHint)

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in startState, see "}", go to objectTagCloseState return real "}"
	expect(t, p.Next, parse.ObjectCloseHint)

	// in objectTagCloseState, go up and return fake "}"
	expect(t, p.Next, parse.ObjectCloseHint)
	// in objectTagCloseState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
