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

func TestTagMixObjectWithIndexes(t *testing.T) {
	s := `[{"mykey1":[{"mykey2":[]}]}]`
	// will be parsed the same as : {"array":[0: {"object":{"mykey1":{"array":[0: {"object":{"mykey2":{"array":[]}}}]}}}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// 1: first array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 2: first object
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
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
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
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

func TestTagMixObjectWithoutIndexes(t *testing.T) {
	s := `[{"mykey1":[{"mykey2":[]}]}]`
	// will be parsed the same as : {"array":[{"object":{"mykey1":{"array":[{"object":{"mykey2":{"array":[]}}}]}}}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags())

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

func TestSimpleTagMixObjectWithIndexes(t *testing.T) {
	s := `[{"mykey1":true}]`
	// will be parsed the same as : {"array":[0: {"object":{"mykey1":true}}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// 1: first array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 2: first object
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey1")
	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

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

func TestSimpleTagMixTwoObjectsWithIndexes(t *testing.T) {
	s := `[{"mykey1":true}, {"mykey2":false}]`
	// will be parsed the same as : {"array":[0: {"object":{"mykey1":true}}, 1: {"object":{"mykey2":false}}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// 1: first array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 2: first object
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey1")
	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)
	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.Hint(t, p, parse.ObjectCloseHint)

	// 2: second object
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey2")
	expect.Hint(t, p, parse.ValueHint)
	expect.False(t, p)
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

func TestSimpleTagMixObjectArrayWithIndexes(t *testing.T) {
	s := `[{"mykey1":[]}]`
	// will be parsed the same as : {"array":[0: {"object":{"mykey1":{"array": []}}}]}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// 1: first array
	expect.Hint(t, p, parse.ObjectOpenHint)
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "array")
	expect.Hint(t, p, parse.ArrayOpenHint)

	// 2: first object
	expect.Hint(t, p, parse.KeyHint)
	expect.Int(t, p, 0)
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
