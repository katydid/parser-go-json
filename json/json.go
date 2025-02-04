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
	"io"

	"github.com/katydid/parser-go-json/json/internal/pool"
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
	state
	stack []state
	pool  pool.Pool
}

type state struct {
	buf    []byte
	offset int

	nextErr error

	kind Kind

	arrayIndex int

	scannedValue bool

	scanned    []byte
	scannedErr error

	parsed             bool
	parsedErr          error
	parsedKindOfNumber kindOfNumber
	parsedDouble       float64
	parsedInt          int64
	parsedUint         uint64
	parsedString       string
}

// NewParser returns a new JSON parser.
func NewParser() Interface {
	return &jsonParser{
		pool:  pool.New(),
		state: state{},
		stack: make([]state, 0, 10),
	}
}

func (s *jsonParser) Init(buf []byte) error {
	s.state = state{buf: buf}
	s.stack = s.stack[:0]
	s.pool.FreeAll()
	return nil
}

func (s *jsonParser) Kind() Kind {
	return s.kind
}

func (s *jsonParser) Next() error {
	if s.nextErr != nil {
		return s.nextErr
	}
	switch s.kind {
	case objectKind:
		return s.scanKeyValue()
	case arrayKind:
		return s.scanElement()
	case stringKind, numberKind, trueKind, falseKind, nullKind:
		if s.scanned != nil || s.scannedErr != nil {
			return io.EOF
		}
		return nil
	default:
		c, err := s.look()
		if err != nil {
			return err
		}
		s.kind = getKind(c)
		switch s.kind {
		case objectKind:
			return s.scanKeyValue()
		case arrayKind:
			return s.scanElement()
		}
		_, err = s.scan()
		return err
	}
}

func (s *jsonParser) Down() {
	if !s.kind.isArray() && !s.kind.isObject() {
		s.nextErr = errNotLeaf
		return
	}
	if s.kind.isObject() {
		_, err := s.scan()
		if err != nil {
			s.nextErr = err
		}
		if err := s.scanColon(); err != nil {
			s.nextErr = err
		}
	} else if s.kind.isArray() {
		_, err := s.scan()
		if err != nil {
			s.nextErr = err
		}
		if s.isNext(',') {
			if err := s.scanComma(); err != nil {
				s.nextErr = err
			}
		}
	}
	s.stack = append(s.stack, s.state)
	s.state = state{
		buf: s.buf[s.offset:],
	}
}

func (s *jsonParser) Up() {
	err := s.Next()
	for err == nil {
		err = s.Next()
	}
	top := len(s.stack) - 1
	s.stack[top].offset += s.offset
	s.state = s.stack[top]
	s.stack = s.stack[:top]
	if err != io.EOF {
		s.nextErr = err
	}
	s.state.scannedValue = true
}

func (s *jsonParser) IsLeaf() bool {
	switch s.kind {
	case stringKind, numberKind, trueKind, falseKind, nullKind:
		return true
	}
	return false
}

func (s *jsonParser) Bool() (bool, error) {
	if !s.kind.isTrue() || s.kind.isFalse() {
		return false, parser.ErrNotBool
	}
	if err := s.parse(); err != nil {
		return false, err
	}
	if s.kind.isTrue() {
		return true, nil
	}
	if s.kind.isFalse() {
		return false, nil
	}
	return false, parser.ErrNotBool
}

func (s *jsonParser) Int() (int64, error) {
	if !s.kind.isNumber() && !s.kind.isArray() {
		return 0, parser.ErrNotInt
	}
	if err := s.parse(); err != nil {
		return 0, err
	}
	if !s.parsedKindOfNumber.isInt() {
		return 0, parser.ErrNotInt
	}
	return s.parsedInt, nil
}

func (s *jsonParser) Uint() (uint64, error) {
	if !s.kind.isNumber() {
		return 0, parser.ErrNotUint
	}
	if err := s.parse(); err != nil {
		return 0, err
	}
	if !s.parsedKindOfNumber.isUint() {
		return 0, parser.ErrNotUint
	}
	return s.parsedUint, nil
}

func (s *jsonParser) Double() (float64, error) {
	if !s.kind.isNumber() {
		return 0, parser.ErrNotDouble
	}
	if err := s.parse(); err != nil {
		return 0, err
	}
	if !s.parsedKindOfNumber.isDouble() {
		return 0, parser.ErrNotDouble
	}
	return s.parsedDouble, nil
}

func (s *jsonParser) String() (string, error) {
	if !s.kind.isString() && !s.kind.isObject() {
		return "", parser.ErrNotString
	}
	if err := s.parse(); err != nil {
		return "", err
	}
	return s.parsedString, nil
}

func (s *jsonParser) Bytes() ([]byte, error) {
	scanned, err := s.scan()
	if err != nil {
		return nil, err
	}
	return scanned, nil
}

func (s *jsonParser) look() (byte, error) {
	if s.offset < len(s.buf) {
		return s.buf[s.offset], nil
	}
	return 0, io.ErrShortBuffer
}

func (s *jsonParser) isNext(c byte) bool {
	return s.offset < len(s.buf) && s.buf[s.offset] == c
}

func (s *jsonParser) incOffset(o int) error {
	s.offset = s.offset + o
	if s.offset > len(s.buf) {
		return io.ErrShortBuffer
	}
	return nil
}

func (s *jsonParser) startedParsing() bool {
	return s.offset > 0
}

func (s *jsonParser) eof() error {
	if err := s.skipSpace(); err != nil {
		return err
	}
	s.nextErr = io.EOF
	return io.EOF
}

func (s *jsonParser) scanElement() error {
	if !s.startedParsing() {
		s.arrayIndex = 0
		if err := s.skipSpace(); err != nil {
			return err
		}
		if err := s.scanOpenArray(); err != nil {
			return err
		}
		if err := s.skipSpace(); err != nil {
			return err
		}
		if s.isNext(']') {
			return s.scanCloseArray()
		}
		return nil
	}
	s.arrayIndex += 1
	if !s.scannedValue {
		// skips pass the whole element
		s.skip()
	}
	if s.isNext(',') {
		s.scanned = nil // clear scanned cache for next element
		if err := s.scanComma(); err != nil {
			return err
		}
		return nil
	} else if s.isNext(']') {
		return s.scanCloseArray()
	}
	return errExpectedCommaOrCloseBracket
}

func (s *jsonParser) scanKeyValue() error {
	if !s.startedParsing() {
		if err := s.skipSpace(); err != nil {
			return err
		}
		if err := s.scanOpenObject(); err != nil {
			return err
		}
		if err := s.skipSpace(); err != nil {
			return err
		}
		if s.isNext('}') {
			return s.scanCloseObject()
		}
		return nil
	}
	if !s.scannedValue {
		// scans the key, if not already scanned
		if _, err := s.scan(); err != nil {
			return err
		}
		if err := s.scanColon(); err != nil {
			return err
		}
		// skips pass the whole value
		s.skip()
		s.scannedValue = false
	}
	if s.isNext(',') {
		s.scanned = nil // clear scanned cache for next key
		if err := s.scanComma(); err != nil {
			return err
		}
		return nil
	} else if s.isNext('}') {
		return s.scanCloseObject()
	}
	return errExpectedCommaOrCloseCurly
}

func (s *jsonParser) skip() {
	s.Down()
	s.Up()
}
