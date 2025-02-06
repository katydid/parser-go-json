//  Copyright 2013 Walter Schulze
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

// Package json contains the implementation of a JSON parser.
package json

import (
	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/parser"
)

// Interface is a parser for JSON
type Interface interface {
	parser.Interface
	//Init initialises the parser with a byte buffer containing JSON.
	Init(buf []byte) error
	Kind() Kind
}

type jsonParser struct {
	state  state
	stack  []state
	parser parse.Parser
	pool   pool.Pool
}

// NewParser returns a new JSON parser.
func NewParser() Interface {
	return &jsonParser{
		pool:  pool.New(),
		stack: make([]state, 0, 10),
	}
}

func (p *jsonParser) Init(buf []byte) error {
	p.pool.FreeAll()
	return nil
}

func (p *jsonParser) Kind() Kind {
	return p.kind
}

func (p *jsonParser) firstNext() error {
	parseKind, err := p.parser.Next()
	if err != nil {
		return err
	}
	switch parseKind {
	case parse.ObjectOpenKind:
		p.state = inObjectState
	case parse.ArrayOpenKind:
		p.state = inArrayState
	case parse.NullKind, parse.BoolKind, parse.StringKind, parse.NumberKind:
		p.state = inLeafState
	}
}

func (p *jsonParser) nextKeyValue() error {
	return p.parser.Skip()
}

func (p *jsonParser) Next() error {
	switch p.state {
	case atStartState:
		return p.firstNext()
	case inObjectState:
		return p.nextKeyValue()
	case goIntoKeyState:
		return p.Next()
	}
}

func (p *jsonParser) Down() {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	// Create a new state.
	switch p.state {
	case inObjectState:
		p.state = goIntoKeyState
	case inArrayState:
		p.state = goIntoElemState
	}
}

func (p *jsonParser) Up() {
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
}

func (p *jsonParser) IsLeaf() bool {
	return p.state == inLeafState
}

func (p *jsonParser) Bool() (bool, error) {
	return p.parser.Bool()
}

func (p *jsonParser) Int() (int64, error) {
	return p.parser.Int()
}

func (p *jsonParser) Uint() (uint64, error) {
	return p.parser.Uint()
}

func (p *jsonParser) Double() (float64, error) {
	return p.parser.Double()
}

func (p *jsonParser) String() (string, error) {
	return p.parser.String()
}

func (p *jsonParser) Bytes() ([]byte, error) {
	return p.parser.Bytes()
}
