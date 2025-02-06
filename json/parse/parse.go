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

//	{"num":3.14,"arr":[null,false,true],"obj":{"k":"v","boring":[1,2,3]}}
//
// Can be parsed using Next, Skip and tokenizer methods.
//
//	p.Next() // {
//
//	p.Next() // "
//	p.String() // "num"
//
//	p.Next() // 0
//	p.Int() // token.ErrNotInt
//	p.Uint() // token.ErrNotInt
//	p.Double() // 3.14
//
//	p.Next() // "
//	p.String() // "arr"
//
//	p.Next() // [
//
//	p.Next() // n
//
//	p.Next() // f
//	p.Bool() // false
//
//	p.Next() // t
//	p.Bool() // true
//
//	p.Next() // ]
//
//	p.Next() // "
//	p.String() // "obj"
//
//	p.Next() // {
//
//	p.Next() // "
//	p.String() // "k"
//
//	p.Next() // "
//	p.String() // "v"
//
//	p.Next() // "
//	p.String() // "boring"
//
//	p.Skip()
//
//	p.Next() // }
//
//	p.Next() // }
package parse

import (
	"io"

	"github.com/katydid/parser-go-json/json/scan"
	"github.com/katydid/parser-go-json/json/token"
)

type Parser interface {
	// Next returns the Kind of the token or an error.
	Next() (Kind, error)
	// Skip skips over the value of an object.
	Skip() error

	Bool() (bool, error)
	Int() (int64, error)
	Uint() (uint64, error)
	Double() (float64, error)
	String() (string, error)
	Bytes() ([]byte, error)

	Init([]byte)
}

type parser struct {
	state
	stack     []state
	alloc     func(int) []byte
	tokenizer token.Tokenizer
}

type state struct {
	kind stateKind
}

func NewParser(buf []byte) Parser {
	alloc := func(size int) []byte { return make([]byte, size) }
	return NewParserWithCustomAllocator(buf, alloc)
}

func NewParserWithCustomAllocator(buf []byte, alloc func(int) []byte) Parser {
	p := &parser{
		state:     state{},
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
	p.state = state{}
	// Shrink the stack's length, but keep it's capacity,
	// so we can reuse it on the next parse.
	p.stack = p.stack[:0]
}

func (p *parser) assertValue(scanKind scan.Kind) (Kind, error) {
	switch scanKind {
	case scan.NullKind, scan.FalseKind, scan.TrueKind, scan.NumberKind, scan.StringKind:
		return Kind(scanKind), nil
	case scan.ArrayOpenKind:
		return Kind(scanKind), nil
	case scan.ObjectOpenKind:
		return Kind(scanKind), nil
	}
	return UnknownKind, errExpectedValue
}

func (p *parser) nextUnknown() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	kind, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownKind, err
	}
	if kind == ObjectOpenKind {
		p.kind = objectOpenStateKind
		p.down(objectOpenStateKind)
	}
	if kind == ArrayOpenKind {
		p.kind = arrayOpenStateKind
		p.down(arrayOpenStateKind)
	}
	return kind, nil
}

func (p *parser) maybeDown(kind Kind) {
	if kind == ObjectOpenKind {
		p.down(objectOpenStateKind)
	}
	if kind == ArrayOpenKind {
		p.down(arrayOpenStateKind)
	}
}

func (p *parser) nextValue() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
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
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, nil
		}
		return ArrayCloseKind, nil
	}
	kind, err := p.assertValue(scanKind)
	if err != nil {
		return UnknownKind, err
	}
	p.kind = arrayElementStateKind
	p.maybeDown(kind)
	return kind, nil
}

func (p *parser) nextArrayElement() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ArrayCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, nil
		}
		return ArrayCloseKind, nil
	}
	if scanKind == scan.CommaKind {
		return p.nextValue()
	}
	return UnknownKind, errExpectedCommaOrCloseBracket
}

func (p *parser) firstObjectKey() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, nil
		}
		return ObjectCloseKind, nil
	}
	if scanKind == scan.StringKind {
		p.kind = objectValueStateKind
		return StringKind, nil
	}
	return UnknownKind, errExpectedStringOrCloseCurly
}

func (p *parser) nextObjectKey() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind == scan.ObjectCloseKind {
		if err := p.up(); err != nil {
			return UnknownKind, nil
		}
		return ObjectCloseKind, nil
	}
	if scanKind == scan.CommaKind {
		nextScanKind, err := p.tokenizer.Next()
		if err != nil {
			return UnknownKind, err
		}
		if nextScanKind == scan.StringKind {
			p.kind = objectValueStateKind
			return StringKind, nil
		} else {
			return UnknownKind, errExpectedCommaOrCloseBracket
		}
	}
	return UnknownKind, errExpectedCommaOrCloseBracket
}

func (p *parser) nextObjectValue() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	if scanKind != scan.ColonKind {
		return UnknownKind, errExpectedColon
	}
	p.kind = objectKeyStateKind
	return p.nextValue()
}

func (p *parser) Next() (Kind, error) {
	switch p.kind {
	case unknownStateKind:
		return p.nextUnknown()
	case arrayOpenStateKind:
		return p.firstArrayElement()
	case arrayElementStateKind:
		return p.nextArrayElement()
	case objectOpenStateKind:
		return p.firstObjectKey()
	case objectKeyStateKind:
		return p.nextObjectKey()
	case objectValueStateKind:
		return p.nextObjectValue()
	default:
		panic("unreachable")
	}
}

func (p *parser) Skip() error {
	if p.kind != objectValueStateKind {
		_, err := p.Next()
		return err
	}
	currentLevel := len(p.stack)
	_, err := p.nextObjectValue()
	if err != nil {
		return err
	}
	for len(p.stack) > currentLevel {
		_, err := p.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) down(stateKind stateKind) {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state{
		kind: stateKind,
	}
}

func (p *parser) up() error {
	if len(p.stack) == 0 {
		return io.EOF
	}
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
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
