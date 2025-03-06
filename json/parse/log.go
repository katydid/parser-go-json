//  Copyright 2015 Walter Schulze
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
	"github.com/katydid/parser-go-json/json/token"
	"github.com/katydid/parser-go/parser/debug"
)

type l struct {
	name   string
	p      Parser
	l      debug.Logger
	copies int
}

// NewLogger returns a parser that when called returns and logs the value returned.
// This is only to be used for debugging purposes.
func NewLogger(p Parser) Parser {
	return &l{"parser", p, debug.NewLineLogger(), 0}
}

func (l *l) Init(buf []byte) {
	l.p.Init(buf)
	l.l.Printf(l.name + ".Init(...)")
}

func (l *l) Skip() error {
	err := l.p.Skip()
	l.l.Printf(l.name+".Double() (%v)", err)
	return err
}

func (l *l) Next() (Hint, error) {
	v, err := l.p.Next()
	l.l.Printf(l.name+".Next() (%v, %v)", v, err)
	return v, err
}

func (l *l) Tokenize() (token.Kind, error) {
	v, err := l.p.Tokenize()
	l.l.Printf(l.name+".Tokenize() (%v, %v)", v, err)
	return v, err
}

func (l *l) Double() (float64, error) {
	v, err := l.p.Double()
	l.l.Printf(l.name+".Double() (%v, %v)", v, err)
	return v, err
}

func (l *l) Int() (int64, error) {
	v, err := l.p.Int()
	l.l.Printf(l.name+".Int() (%v, %v)", v, err)
	return v, err
}

func (l *l) Bytes() ([]byte, error) {
	v, err := l.p.Bytes()
	l.l.Printf(l.name+".Bytes() (%v, %v)", v, err)
	return v, err
}
