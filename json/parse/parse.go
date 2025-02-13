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
	// Next returns the Kind of the token or an error.
	Next() (Kind, error)
	// Skip skips over some of the json string.
	// If the kind '{' was returned by Next, then the whole object is skipped.
	// If a object key was just parsed, then that key's value is skipped.
	// If an object value was just parsed, then the rest of the object is skipped.
	// If the kind '[' was returned by Next, then the whole array is skipped.
	// If an array element was parsed, then the rest of the array is skipped.
	// In any other case, Skip simply calls Next.
	Skip() error

	// Bool attempts to convert the current token to a bool.
	Bool() (bool, error)
	// Int attempts to convert the current token to an int64.
	Int() (int64, error)
	// Uint attempts to convert the current token to an uint64.
	Uint() (uint64, error)
	// Double attempts to convert the current token to a float64.
	Double() (float64, error)
	// String attempts to convert the current token to a string.
	String() (string, error)
	// Bytes returns the raw current token.
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

func (p *parser) assertValue(scanKind scan.Kind) (Kind, error) {
	switch scanKind {
	case scan.NullKind, scan.NumberKind, scan.StringKind:
		return Kind(scanKind), nil
	case scan.FalseKind, scan.TrueKind:
		return BoolKind, nil
	case scan.ArrayOpenKind:
		return Kind(scanKind), nil
	case scan.ObjectOpenKind:
		return Kind(scanKind), nil
	}
	return UnknownKind, errExpectedValue
}

func (p *parser) nextStart() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	kind, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownKind, err
	}
	if kind == ObjectOpenKind {
		p.state = objectOpenState
		p.down(objectOpenState)
	} else if kind == ArrayOpenKind {
		p.state = arrayOpenState
		p.down(arrayOpenState)
	} else {
		p.state = leafState
	}
	return kind, nil
}

func (p *parser) maybeDown(kind Kind) {
	if kind == ObjectOpenKind {
		p.down(objectOpenState)
	}
	if kind == ArrayOpenKind {
		p.down(arrayOpenState)
	}
}

func (p *parser) nextValue() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	kind, err := p.assertValue(scanKind)
	if err != nil {
		return kind, err
	}
	p.maybeDown(kind)
	return kind, nil
}

func (p *parser) firstArrayElement() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, err
		}
		return ArrayCloseKind, nil
	}
	kind, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownKind, err
	}
	p.state = arrayElementState
	p.maybeDown(kind)
	return kind, nil
}

func (p *parser) nextArrayElement() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, err
		}
		return ArrayCloseKind, nil
	}
	if scanKind == scan.CommaKind {
		return p.nextValue()
	}
	return UnknownKind, errExpectedCommaOrCloseBracket
}

func (p *parser) firstObjectKey() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, err
		}
		return ObjectCloseKind, nil
	}
	if scanKind == scan.StringKind {
		p.state = objectValueState
		return StringKind, nil
	}
	return UnknownKind, errExpectedStringOrCloseCurly
}

func (p *parser) nextObjectKey() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, err
		}
		return ObjectCloseKind, nil
	}
	if scanKind == scan.CommaKind {
		nextScanKind, err := p.nextToken()
		if err != nil {
			return UnknownKind, err
		}
		if nextScanKind == scan.StringKind {
			p.state = objectValueState
			return StringKind, nil
		} else {
			return UnknownKind, errExpectedCommaOrCloseBracket
		}
	}
	return UnknownKind, errExpectedCommaOrCloseBracket
}

func (p *parser) nextObjectValue() (Kind, error) {
	scanKind, err := p.nextToken()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind != scan.ColonKind {
		return UnknownKind, errExpectedColon
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

func (p *parser) Next() (Kind, error) {
	switch p.state {
	case startState:
		return p.nextStart()
	case leafState:
		// leaf was already parsed, so there should be nothing left to parse
		return UnknownKind, p.eof()
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
		return UnknownKind, p.eof()
	default:
		panic("unreachable")
	}
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

func (p *parser) Bool() (bool, error) {
	return p.tokenizer.Bool()
}

func (p *parser) Int() (int64, error) {
	return p.tokenizer.Int()
}

func (p *parser) Uint() (uint64, error) {
	return p.tokenizer.Uint()
}

func (p *parser) Double() (float64, error) {
	return p.tokenizer.Double()
}

func (p *parser) String() (string, error) {
	return p.tokenizer.String()
}

func (p *parser) Bytes() ([]byte, error) {
	return p.tokenizer.Bytes()
}
