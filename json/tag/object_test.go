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
	"testing"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestTagObjectForEmptyObject(t *testing.T) {
	s := `{}`
	// will be parsed the same as : {"object": {}}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)
	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	expect.EOF(t, p)
}

func TestTagObjectForNonEmptyObject(t *testing.T) {
	s := `{"mykey": "myvalue"}`
	// will be parsed the same as : {"object": {"mykey": "myvalue"}}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey")

	// in startState, see "myvalue"
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "myvalue")

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in endState, at top of stack return EOF
	expect.EOF(t, p)
}

func TestTagObjectForNonEmptyObjectWithEmptyObjectValue(t *testing.T) {
	s := `{"mykey": {}}`
	// will be parsed the same as : {"object": {"mykey": {"object": {}}}}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in objectTagCloseState, at top of stack return EOF
	expect.EOF(t, p)
}

func TestTagObjectForNonEmptyObjectWithNonEmptyObjectValue(t *testing.T) {
	s := `{"mykey": {"mykey2": {}}}`
	// will be parsed the same as : {"object": {"mykey": {"object": {"mykey2": {"object": {}}}}}}
	p := jsonparse.NewParser(jsonparse.WithBuffer([]byte(s)), jsonparse.WithTags(), jsonparse.WithIndexes())

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "mykey"
	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "mykey2")

	// in startState, see "{", go down to objectTagOpenState and return fake "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in objectTagOpenState, go to objectTagKeyOpenState and return fake key "object"
	expect.Hint(t, p, parse.KeyHint)
	expect.Tag(t, p, "object")

	// in objectTagKeyOpenState, go to objectTagKeyCloseState and down to startState and return real "{"
	expect.Hint(t, p, parse.ObjectOpenHint)

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in startState, see "}", go up to objectTagKeyCloseState return real "}"
	expect.Hint(t, p, parse.ObjectCloseHint)

	// in objectTagKeyCloseState, go up and return fake "}"
	expect.Hint(t, p, parse.ObjectCloseHint)
	// in objectTagCloseState, at top of stack return EOF
	expect.EOF(t, p)
}
