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
	"github.com/katydid/parser-go-json/json/token"
)

type Parser interface {
	// Next returns the Hint of the token or an error.
	Next() (Hint, error)

	// Skip allows the user to skip over uninteresting parts of the parse tree.
	// Based on the Hint skip has different intuitive behaviours.
	// If the Hint was:
	// * '{': the whole Map is skipped.
	// * 'k': the key's value is skipped.
	// * '[': the whole List is skipped.
	// * 'v': the rest of the Map or List is skipped.
	// * ']': same as calling Next and ignoring the Hint.
	// * '}': same as calling Next and ignoring the Hint.
	Skip() error

	// Tokenize parses the current token.
	Tokenize() (token.Kind, error)

	// Int attempts to convert the current token to an int64.
	Int() (int64, error)
	// Double attempts to convert the current token to a float64.
	Double() (float64, error)
	// Bytes returns the bytes token or a unquoted string or decimal.
	Bytes() ([]byte, error)
	// Init restarts the parser with a new byte buffer, without allocating a new parser.
	Init([]byte)
}

type parser struct {
	state     state
	stack     []state
	alloc     func(int) []byte
	tokenizer token.Tokenizer
}

func NewParser(buf []byte) Parser {
	alloc := func(size int) []byte { return make([]byte, size) }
	return NewParserWithCustomAllocator(buf, alloc)
}

func NewParserWithCustomAllocator(buf []byte, alloc func(int) []byte) Parser {
	p := &parser{
		state:     startState,
		stack:     make([]state, 0, 10),
		alloc:     alloc,
		tokenizer: token.NewTokenizerWithCustomAllocator(buf, alloc),
	}
	return p
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

func (p *parser) assertValue(scanKind scan.Kind) (Hint, error) {
	switch scanKind {
	case scan.NullKind, scan.FalseKind, scan.TrueKind, scan.NumberKind, scan.StringKind:
		return ValueHint, nil
	case scan.ArrayOpenKind:
		return ArrayOpenHint, nil
	case scan.ObjectOpenKind:
		return ObjectOpenHint, nil
	}
	return UnknownHint, errExpectedValue
}

func (p *parser) nextStart() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownHint, err
	}
	if hint == ObjectOpenHint {
		p.state = objectOpenState
		p.down(objectOpenState)
	} else if hint == ArrayOpenHint {
		p.state = arrayOpenState
		p.down(arrayOpenState)
	} else {
		p.state = leafState
	}
	return hint, nil
}

func (p *parser) maybeDown(hint Hint) {
	if hint == ObjectOpenHint {
		p.down(objectOpenState)
	}
	if hint == ArrayOpenHint {
		p.down(arrayOpenState)
	}
}

func (p *parser) nextValue() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return hint, err
	}
	p.maybeDown(hint)
	return hint, nil
}

func (p *parser) firstArrayElement() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownHint, err
		}
		return ArrayCloseHint, nil
	}
	hint, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownHint, err
	}
	p.state = arrayElementState
	p.maybeDown(hint)
	return hint, nil
}

func (p *parser) nextArrayElement() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownHint, err
		}
		return ArrayCloseHint, nil
	}
	if scanKind == scan.CommaKind {
		return p.nextValue()
	}
	return UnknownHint, errExpectedCommaOrCloseBracket
}

func (p *parser) firstObjectKey() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownHint, err
		}
		return ObjectCloseHint, nil
	}
	if scanKind == scan.StringKind {
		p.state = objectValueState
		return KeyHint, nil
	}
	return UnknownHint, errExpectedStringOrCloseCurly
}

func (p *parser) nextObjectKey() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownHint, err
		}
		return ObjectCloseHint, nil
	}
	if scanKind == scan.CommaKind {
		nextScanKind, err := p.nextToken()
		if err != nil {
			return UnknownHint, err
		}
		if nextScanKind == scan.StringKind {
			p.state = objectValueState
			return KeyHint, nil
		} else {
			return UnknownHint, errExpectedCommaOrCloseBracket
		}
	}
	return UnknownHint, errExpectedCommaOrCloseBracket
}

func (p *parser) nextObjectValue() (Hint, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownHint, err
	}
	if scanKind != scan.ColonKind {
		return UnknownHint, errExpectedColon
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

func (p *parser) Next() (Hint, error) {
	switch p.state {
	case startState:
		return p.nextStart()
	case leafState:
		// leaf was already parsed, so there should be nothing left to parse
		return UnknownHint, p.eof()
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
		return UnknownHint, p.eof()
	default:
		panic("unreachable")
	}
}

func (p *parser) Tokenize() (token.Kind, error) {
	return p.tokenizer.Tokenize()
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

func (p *parser) Int() (int64, error) {
	return p.tokenizer.Int()
}

func (p *parser) Double() (float64, error) {
	return p.tokenizer.Double()
}

func (p *parser) Bytes() ([]byte, error) {
	return p.tokenizer.Bytes()
}
