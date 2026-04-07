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

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)

type Parser interface {
	parse.Parser
	Init([]byte)
}

type tagger struct {
	p     jsonparse.Parser
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
func NewTagger(p jsonparse.Parser, opts ...Option) Parser {
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
	if !t.tag && t.index {
		panic("unsupported options: WithIndexes requires WithTags")
	}
	return t
}

func (t *tagger) Next() (parse.Hint, error) {
	switch t.state.kind {
	case startState:
		h, err := t.p.Next()
		if err != nil {
			return translateHint(h), err
		}
		// helps to skip over object values
		t.state.hint = h
		if t.tag {
			switch h {
			case jsonparse.ObjectOpenHint:
				t.down(objectTagOpenState)
				return parse.EnterHint, nil
			case jsonparse.ObjectCloseHint:
				if err := t.up(); err != nil {
					return parse.UnknownHint, err
				}
				return parse.LeaveHint, nil
			case jsonparse.ArrayOpenHint:
				t.down(arrayTagOpenState)
				return parse.EnterHint, nil
			case jsonparse.ArrayCloseHint:
				if err := t.up(); err != nil {
					return parse.UnknownHint, err
				}
				return parse.LeaveHint, nil
			}
		}
		return translateHint(h), nil
	case objectTagOpenState:
		t.state.kind = objectTagKeyOpenState
		return parse.FieldHint, nil
	case objectTagKeyOpenState:
		t.state.kind = objectTagKeyCloseState
		t.down(startState)
		return parse.EnterHint, nil
	case objectTagKeyCloseState:
		if err := t.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.LeaveHint, nil
	case objectTagCloseState:
		if err := t.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.LeaveHint, nil
	case arrayTagOpenState:
		t.state.kind = arrayTagKeyOpenState
		return parse.FieldHint, nil
	case arrayTagKeyOpenState:
		t.state.kind = arrayTagKeyCloseState
		if t.index {
			t.down(arrayTagIndexState)
		} else {
			t.down(startState)
		}
		return parse.EnterHint, nil
	case arrayTagKeyCloseState:
		if err := t.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.LeaveHint, nil
	case arrayTagIndexState:
		h, err := t.p.Next()
		if err != nil {
			return translateHint(h), err
		}
		t.state.hint = h
		if t.state.hint == jsonparse.ArrayCloseHint {
			if err := t.up(); err != nil {
				return parse.UnknownHint, err
			}
			return parse.LeaveHint, nil
		}
		t.state.arrayIndex++
		t.state.kind = arrayTagElemState
		return parse.FieldHint, nil
	case arrayTagElemState:
		t.state.kind = arrayTagIndexState
		h := t.state.hint
		if t.tag {
			switch h {
			case jsonparse.ObjectOpenHint:
				t.down(objectTagOpenState)
				return parse.EnterHint, nil
			case jsonparse.ObjectCloseHint:
				if err := t.up(); err != nil {
					return parse.UnknownHint, err
				}
				return parse.LeaveHint, nil
			case jsonparse.ArrayOpenHint:
				t.down(arrayTagOpenState)
				return parse.EnterHint, nil
			case jsonparse.ArrayCloseHint:
				if err := t.up(); err != nil {
					return parse.UnknownHint, err
				}
				return parse.LeaveHint, nil
			}
		}
		return translateHint(h), nil
	case endState:
		return parse.UnknownHint, io.EOF
	}
	panic(fmt.Sprintf("unreachable: unknown state = %v", t.state))
}

func (t *tagger) Skip() error {
	switch t.state.kind {
	case startState:
		if !t.tag {
			return t.p.Skip()
		}
		if len(t.stack) == 0 {
			_, err := t.Next()
			return err
		}
		if t.state.hint != jsonparse.KeyHint {
			// do not go up when it is an object value that needs to be skipped over
			if err := t.up(); err != nil {
				return err
			}
		}
		return t.p.Skip()
	case objectTagOpenState:
		if err := t.up(); err != nil {
			return err
		}
		return t.p.Skip()
	case objectTagKeyOpenState:
		t.state.kind = arrayTagKeyCloseState
		return nil
	case objectTagKeyCloseState:
		_, err := t.Next()
		return err
	case arrayTagOpenState:
		if err := t.up(); err != nil {
			return err
		}
		return t.p.Skip()
	case arrayTagKeyOpenState:
		t.state.kind = arrayTagKeyCloseState
		return nil
	case arrayTagKeyCloseState:
		_, err := t.Next()
		return err
	case arrayTagIndexState:
		if err := t.up(); err != nil {
			return err
		}
		return t.p.Skip()
	case arrayTagElemState:
		t.state.kind = arrayTagIndexState
		if t.state.hint == jsonparse.ValueHint {
			// values do not need to be skipped, Next will take care of it.
			return nil
		}
		return t.p.Skip()
	case endState:
		return t.p.Skip()
	}
	panic(fmt.Sprintf("unreachable: unknown state = %v", t.state))
}

func (t *tagger) Token() (parse.Kind, []byte, error) {
	switch t.state.kind {
	case objectTagKeyOpenState:
		return parse.TagKind, objectTagToken, nil
	case arrayTagKeyOpenState:
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
	t.state.arrayIndex = -1
}

func (t *tagger) up() error {
	if len(t.stack) == 0 {
		return errUnexpectedClose
	}
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
