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
	"fmt"
	"io"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)

// Parser is a copy of the json.Parse interface.
type Parser interface {
	Next() (parse.Hint, error)
	Skip() error
	Token() (parse.Kind, []byte, error)
	Init([]byte)
}

type tagger struct {
	p     Parser
	tag   bool
	index bool
	alloc func(size int) []byte
	// state
	state state
	stack []state
}

var objectTagToken = []byte("object")
var arrayTagToken = []byte("array")

// NewTagger can tag objects and arrays.
// The following json: `{"a": []}`
// is parsed as: `{"object": {"a": {"array": []}}}`.
// The kind returned from the Token method for
// "object" and "array" will be parse.TagKind.
func NewTagger(p Parser, opts ...Option) Parser {
	t := &tagger{
		p:     p,
		tag:   false,
		index: false,
		alloc: func(size int) []byte {
			return make([]byte, size)
		},
		state: state{},
		stack: make([]state, 0, 10),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *tagger) Next() (parse.Hint, error) {
	switch t.state.kind {
	case startState:
		h, err := t.p.Next()
		if err != nil {
			return h, err
		}
		if t.tag {
			switch h {
			case parse.ObjectOpenHint:
				t.down(objectTagOpenState)
				return parse.ObjectOpenHint, nil
			case parse.ObjectCloseHint:
				t.state.kind = objectTagCloseState
				return parse.ObjectCloseHint, nil
			case parse.ArrayOpenHint:
				t.down(arrayTagOpenState)
				return parse.ObjectOpenHint, nil
			case parse.ArrayCloseHint:
				t.state.kind = arrayTagCloseState
				return parse.ArrayCloseHint, nil
			}
		}
		return h, nil
	case objectTagOpenState:
		t.state.kind = objectTagKeyState
		return parse.KeyHint, nil
	case objectTagKeyState:
		t.state.kind = startState
		return parse.ObjectOpenHint, nil
	case objectTagCloseState:
		t.up()
		return parse.ObjectCloseHint, nil
	case arrayTagOpenState:
		t.state.kind = arrayTagKeyState
		return parse.KeyHint, nil
	case arrayTagKeyState:
		if t.index {
			t.state.kind = arrayTagIndexState
			t.state.arrayIndex = -1
			return parse.ArrayOpenHint, nil
		}
		t.state.kind = startState
		return parse.ArrayOpenHint, nil
	case arrayTagIndexState:
		h, err := t.p.Next()
		if err != nil {
			return h, err
		}
		t.state.arrayElemHint = h
		if t.state.arrayElemHint == parse.ArrayCloseHint {
			t.state.kind = arrayTagCloseState
			return parse.ArrayCloseHint, nil
		}
		t.state.arrayIndex++
		t.state.kind = arrayTagElemState
		return parse.KeyHint, nil
	case arrayTagElemState:
		t.state.kind = arrayTagIndexState
		h := t.state.arrayElemHint
		if t.tag {
			switch h {
			case parse.ObjectOpenHint:
				t.down(objectTagOpenState)
				return parse.ObjectOpenHint, nil
			case parse.ObjectCloseHint:
				t.state.kind = objectTagCloseState
				return parse.ObjectCloseHint, nil
			case parse.ArrayOpenHint:
				t.down(arrayTagOpenState)
				return parse.ObjectOpenHint, nil
			case parse.ArrayCloseHint:
				t.state.kind = arrayTagCloseState
				return parse.ArrayCloseHint, nil
			}
		}
		return h, nil
	case arrayTagCloseState:
		t.up()
		return parse.ObjectCloseHint, nil
	case endState:
		return parse.UnknownHint, io.EOF
	}
	panic(fmt.Sprintf("unreachable: unknown state = %v", t.state))
}

func (t *tagger) Skip() error {
	return t.p.Skip()
}

func (t *tagger) Token() (parse.Kind, []byte, error) {
	switch t.state.kind {
	case objectTagKeyState:
		return parse.TagKind, objectTagToken, nil
	case arrayTagKeyState:
		return parse.TagKind, arrayTagToken, nil
	case arrayTagElemState:
		return parse.Int64Kind, cast.FromInt64(t.state.arrayIndex, t.alloc), nil
	}
	return t.p.Token()
}

func (t *tagger) Init(buf []byte) {
	// Reset the state.
	t.state.kind = startState
	// Shrink the stack's length, but keep it's capacity,
	// so we can reuse it on the next parse.
	t.stack = t.stack[:0]
	// Reset the parser too.
	t.p.Init(buf)
}

func (t *tagger) down(stateKind stateKind) {
	// Append the current state to the stack.
	t.stack = append(t.stack, t.state)
	// Create a new state.
	t.state.kind = stateKind
}

func (t *tagger) up() error {
	top := len(t.stack) - 1
	// Set the current state to the state on top of the stack.
	t.state = t.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	t.stack = t.stack[:top]
	if len(t.stack) == 0 {
		t.state.kind = endState
	}
	return nil
}
