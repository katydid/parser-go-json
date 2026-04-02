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
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestTagArrayForEmptyArrayWithIndex(t *testing.T) {
	s := `[]`
	// will be parsed the same as : {"array": []}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)
	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForEmptyArrayWithoutIndex(t *testing.T) {
	s := `[]`
	// will be parsed the same as : {"array": []}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)
	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForNonEmptyArrayWithIndex(t *testing.T) {
	s := `["myelem"]`
	// will be parsed the same as : {"array": ["myelem"]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)

	// in startState, see "myelem"
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "myelem")

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForNonEmptyArrayWithoutIndex(t *testing.T) {
	s := `["myelem"]`
	// will be parsed the same as : {"array": ["myelem"]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "myelem"
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "myelem")

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayWithEmptyArrayWithIndex(t *testing.T) {
	s := `[[]]`
	// will be parsed the same as : {"array": [{"array": []}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayWithEmptyArrayWithoutIndex(t *testing.T) {
	s := `[[]]`
	// will be parsed the same as : {"array": [{"array": []}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForThreeElementsWithIndex(t *testing.T) {
	s := `["myelem", 789, true]`
	// will be parsed the same as : {"array": [0: "myelem", 1: 789, 2: true]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "myelem"
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "myelem")

	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 789)

	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 2)
	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}

func TestTagArrayForThreeElementsWithoutIndex(t *testing.T) {
	s := `["myelem", 789, true]`
	// will be parsed the same as : {"array": [0: "myelem", 1: 789, 2: true]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags())

	// in startState, see "[", go down to arrayTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in arrayTagOpenState, go to arrayTagKeyState and return fake key "array"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")

	// in arrayTagKeyState, go to startState return real "["
	expect.Hint(t, p, parse.ArrayOpenHint)

	// in startState, see "myelem"
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "myelem")

	expect.Hint(t, p, parse.ValueHint)
	expect.Int(t, p, 789)

	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

	// in startState, see "]", go to arrayTagCloseState return real "]"
	expect.Hint(t, p, parse.ArrayCloseHint)

	// in arrayTagCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
