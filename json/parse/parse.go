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

//	{"num":3.14,"arr":[null,false,true],"obj":{"k":"v", "boring": [1,2,3]}}
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
	kind Kind
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

func (p *parser) Next() (Kind, error) {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return UnknownKind, err
	}
	switch p.kind {
	case UnknownKind:
		switch scanKind {
		case scan.UnknownKind:
			return UnknownKind, nil
		case scan.NullKind:
			p.kind = NullKind
			return p.kind, nil
		case scan.FalseKind:
			p.kind = FalseKind
			return p.kind, nil
		case scan.TrueKind:
			p.kind = TrueKind
			return p.kind, nil
		case scan.NumberKind:
			p.kind = NumberKind
			return p.kind, nil
		case scan.StringKind:
			p.kind = StringKind
			return p.kind, nil
		case scan.ArrayOpenKind:
			p.kind = ArrayKind
			return p.kind, nil
		case scan.CommaKind:
			return UnknownKind, errUnexpectedComma
		case scan.ArrayCloseKind:
			return UnknownKind, errUnexpectedCloseBracket
		case scan.ObjectOpenKind:
			p.kind = ObjectKind
			return p.kind, nil
		case scan.ColonKind:
			return UnknownKind, errUnexpectedColon
		case scan.ObjectCloseKind:
			return UnknownKind, errUnexpectedCloseCurly
		default:
			panic("unreachable")
		}
	case NullKind, FalseKind, TrueKind, StringKind, NumberKind:
		return p.kind, nil
	case ArrayKind:
		switch scanKind {
		case scan.UnknownKind:
			return UnknownKind, nil
		case scan.NullKind:

		case scan.FalseKind:

		case scan.TrueKind:

		case scan.NumberKind:

		case scan.StringKind:

		case scan.ArrayOpenKind:

		case scan.CommaKind:

		case scan.ArrayCloseKind:
			return UnknownKind, io.EOF
		case scan.ObjectOpenKind:

		case scan.ColonKind:
			return UnknownKind, errUnexpectedColon
		case scan.ObjectCloseKind:
			return UnknownKind, errUnexpectedCloseCurly
		default:
			panic("unreachable")
		}
	case ObjectKind:
		switch scanKind {
		case scan.UnknownKind:
			return UnknownKind, nil
		case scan.NullKind:

		case scan.FalseKind:

		case scan.TrueKind:

		case scan.NumberKind:

		case scan.StringKind:

		case scan.ArrayOpenKind:

		case scan.CommaKind:

		case scan.ArrayCloseKind:
			return UnknownKind, errUnexpectedCloseBracket
		case scan.ObjectOpenKind:

		case scan.ColonKind:

		case scan.ObjectCloseKind:
			return UnknownKind, io.EOF
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

func (p *parser) next() error {
	scanKind, err := p.tokenizer.Next()
	if err != nil {
		return err
	}
}

func (p *parser) Down() {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	p.state = state{}
}

func (p *parser) Up() {
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
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
