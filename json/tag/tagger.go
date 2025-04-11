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

	"github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/token"
)

type tagger struct {
	p parse.Parser
	// options
	tagObjects bool
	tagArrays  bool
	// state
	state state
	stack []state
}

func NewTagger(p parse.Parser, opts ...Option) parse.Parser {
	t := &tagger{
		p:     p,
		state: startState,
		stack: make([]state, 0, 10),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *tagger) Next() (parse.Hint, error) {
	switch t.state {
	case startState:
		h, err := t.p.Next()
		if err != nil {
			return h, err
		}
		if !t.tagObjects {
			return h, nil
		}
		switch h {
		case parse.ObjectOpenHint:
			t.down(objectTagOpenState)
			return parse.ObjectOpenHint, nil
		case parse.ObjectCloseHint:
			t.state = objectTagCloseState
			return parse.ObjectCloseHint, nil
		default:
			return h, nil
		}
	case objectTagOpenState:
		t.state = objectTagKeyState
		return parse.KeyHint, nil
	case objectTagKeyState:
		t.state = startState
		return parse.ObjectOpenHint, nil
	case objectTagCloseState:
		t.up()
		return parse.ObjectCloseHint, nil
	case endState:
		return parse.UnknownHint, io.EOF
	}
	panic(fmt.Sprintf("unreachable: unknown state = %c", t.state))
}

func (t *tagger) Skip() error {
	return t.p.Skip()
}

func (t *tagger) Token() (token.Kind, []byte, error) {
	switch t.state {
	case startState:
	case objectTagOpenState:
	case objectTagKeyState:
		return token.TagKind, []byte("object"), nil
	case objectTagCloseState:
	case endState:
	}
	return t.p.Token()
}

func (t *tagger) Init(buf []byte) {
	// Reset the state.
	t.state = startState
	// Shrink the stack's length, but keep it's capacity,
	// so we can reuse it on the next parse.
	t.stack = t.stack[:0]
	// Reset the parser too.
	t.p.Init(buf)
}

func (t *tagger) down(state state) {
	// Append the current state to the stack.
	t.stack = append(t.stack, t.state)
	// Create a new state.
	t.state = state
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
		t.state = endState
	}
	return nil
}
