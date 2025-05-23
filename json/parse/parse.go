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

	"github.com/katydid/parser-go-json/json/scan"
	"github.com/katydid/parser-go-json/json/tag"
	"github.com/katydid/parser-go-json/json/token"
	"github.com/katydid/parser-go/parse"
)

type Parser interface {
	parse.Parser

	// Init restarts the parser with a new byte buffer, without allocating a new parser.
	Init([]byte)
}

type parser struct {
	// state
	state state
	stack []state

	// initialized via options
	tokenizer token.Tokenizer
}

func NewParser(opts ...Option) Parser {
	p := &parser{
		state: startState,
		stack: make([]state, 0, 10),
	}
	options := newOptions(opts...)
	p.tokenizer = token.NewTokenizerWithCustomAllocator(options.buf, options.alloc)
	if !options.tagArrays && !options.tagObjects {
		return p
	}
	tagOptions := []tag.Option{}
	if options.tagArrays {
		tagOptions = append(tagOptions, tag.WithArrayTag())
	}
	if options.tagObjects {
		tagOptions = append(tagOptions, tag.WithObjectTag())
	}
	return tag.NewTagger(p, tagOptions...)
}

func (p *parser) Init(buf []byte) {
	// Reset the tokenizer with the new buffer.
	p.tokenizer.Init(buf)
	// Reset the state.
	p.state = startState
	// Shrink the stack's length, but keep it's capacity,
	// so we can reuse it on the next parse.
	p.stack = p.stack[:0]
}

func (p *parser) nextToken() (scan.Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err == nil {
		return scanKind, nil
	}
	if err == io.EOF {
		return scanKind, io.ErrShortBuffer
	}
	return scanKind, err
}

func (p *parser) assertValue(scanKind scan.Kind) (parse.Hint, error) {
	switch scanKind {
	case scan.NullKind, scan.FalseKind, scan.TrueKind, scan.NumberKind, scan.StringKind:
		return parse.ValueHint, nil
	case scan.ArrayOpenKind:
		return parse.ArrayOpenHint, nil
	case scan.ObjectOpenKind:
		return parse.ObjectOpenHint, nil
	}
	return parse.UnknownHint, errExpectedValue
}

func (p *parser) nextStart() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return parse.UnknownHint, err
	}
	if hint == parse.ObjectOpenHint {
		p.state = objectOpenState
		p.down(objectOpenState)
	} else if hint == parse.ArrayOpenHint {
		p.state = arrayOpenState
		p.down(arrayOpenState)
	} else {
		p.state = leafState
	}
	return hint, nil
}

func (p *parser) maybeDown(hint parse.Hint) {
	if hint == parse.ObjectOpenHint {
		p.down(objectOpenState)
	}
	if hint == parse.ArrayOpenHint {
		p.down(arrayOpenState)
	}
}

func (p *parser) nextValue() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return hint, err
	}
	p.maybeDown(hint)
	return hint, nil
}

func (p *parser) firstArrayElement() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.ArrayCloseHint, nil
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return parse.UnknownHint, err
	}
	p.state = arrayElementState
	p.maybeDown(hint)
	return hint, nil
}

func (p *parser) nextArrayElement() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.ArrayCloseHint, nil
	}
	if scanKind == scan.CommaKind {
		return p.nextValue()
	}
	return parse.UnknownHint, errExpectedCommaOrCloseBracket
}

func (p *parser) firstObjectKey() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.ObjectCloseHint, nil
	}
	if scanKind == scan.StringKind {
		p.state = objectValueState
		return parse.KeyHint, nil
	}
	return parse.UnknownHint, errExpectedStringOrCloseCurly
}

func (p *parser) nextObjectKey() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return parse.UnknownHint, err
		}
		return parse.ObjectCloseHint, nil
	}
	if scanKind == scan.CommaKind {
		nextScanKind, err := p.nextToken()
		if err != nil {
			return parse.UnknownHint, err
		}
		if nextScanKind == scan.StringKind {
			p.state = objectValueState
			return parse.KeyHint, nil
		} else {
			return parse.UnknownHint, errExpectedCommaOrCloseBracket
		}
	}
	return parse.UnknownHint, errExpectedCommaOrCloseBracket
}

func (p *parser) nextObjectValue() (parse.Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return parse.UnknownHint, err
	}
	if scanKind != scan.ColonKind {
		return parse.UnknownHint, errExpectedColon
	}
	p.state = objectKeyState
	return p.nextValue()
}

func (p *parser) eof() error {
	if _, err := p.tokenizer.Next(); err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return err
	}
	// There is still bytes left in the buffer, but the object or array is closed.
	return io.ErrShortBuffer
}

func (p *parser) Next() (parse.Hint, error) {
	switch p.state {
	case startState:
		return p.nextStart()
	case leafState:
		// leaf was already parsed, so there should be nothing left to parse
		return parse.UnknownHint, p.eof()
	case arrayOpenState:
		return p.firstArrayElement()
	case arrayElementState:
		return p.nextArrayElement()
	case objectOpenState:
		return p.firstObjectKey()
	case objectKeyState:
		return p.nextObjectKey()
	case objectValueState:
		return p.nextObjectValue()
	case endState:
		return parse.UnknownHint, p.eof()
	default:
		panic("unreachable")
	}
}

func (p *parser) Token() (parse.Kind, []byte, error) {
	return p.tokenizer.Token()
}

func (p *parser) Skip() error {
	switch p.state {
	case arrayOpenState, arrayElementState:
		// '[' has been parsed or
		// '['"e1",...,"en" has been parsed.
		// call Next until ']' is parsed,
		// which will result in the stack being popped,
		// which will result in the stack size being smaller.
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	case objectOpenState, objectKeyState:
		// '{' has been parsed or
		// '{'"k1":"v1",...,"kn":"vn" has been parsed.
		// call Next until '}' is parsed,
		// which will result in the stack being popped,
		// which will result in the stack size being smaller.
		currentStackSize := len(p.stack)
		for len(p.stack) >= currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	case objectValueState:
		currentStackSize := len(p.stack)
		_, err := p.Next()
		if err != nil {
			return err
		}
		// If Next parsed down into an array or object,
		// then keep on parsing until we reach our current level.
		// If Next parsed a string, number, boolean or null,
		// then the level would be the same.
		for len(p.stack) > currentStackSize {
			_, err := p.Next()
			if err != nil {
				return err
			}
		}
	default:
		_, err := p.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) down(state state) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state
}

func (p *parser) up() error {
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
	if len(p.stack) == 0 {
		p.state = endState
	}
	return nil
}
