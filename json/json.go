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
	"bytes"
	"io"

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/internal/strconv"
	"github.com/katydid/parser-go/parser"
)

func scanString(buf []byte) (int, error) {
	escaped := false
	udigits := -1
	if buf[0] != '"' {
		return 0, errScanString
	}
	for i, c := range buf[1:] {
		if escaped {
			switch c {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				escaped = false
				continue
			case 'u':
				udigits = 0
				escaped = false
				continue
			}
			return 0, errScanString
		}
		if udigits >= 0 {
			if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
				udigits++
			} else {
				return 0, errScanString
			}
			if udigits == 4 {
				udigits = -1
			}
			continue
		}
		if c == '"' {
			return i + 2, nil
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c < 0x20 {
			return 0, errScanString
		}
	}
	return 0, errScanString
}

func isSpace(c byte) bool {
	return (c == ' ') || (c == '\n') || (c == '\r') || (c == '\t')
}

func skipSpace(buf []byte) int {
	for i, c := range buf {
		if !isSpace(c) {
			return i
		}
	}
	return len(buf)
}

func unquote(pool pool.Pool, s []byte) (string, error) {
	var ok bool
	var t string
	s, ok = unquoteBytes(pool.Alloc, s)
	t = castToString(s)
	if !ok {
		return "", errUnquote
	}
	return t, nil
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

func (s *jsonParser) isNextDigit() bool {
	return s.offset < len(s.buf) && s.buf[s.offset] >= '0' && s.buf[s.offset] <= '9'
}

func (s *jsonParser) isNextDigit19() bool {
	return s.offset < len(s.buf) && s.buf[s.offset] >= '1' && s.buf[s.offset] <= '9'
}

func (s *jsonParser) isNextOneOf(cs ...byte) bool {
	if s.offset >= len(s.buf) {
		return false
	}
	for _, c := range cs {
		if s.buf[s.offset] == c {
			return true
		}
	}
	return false
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

func (s *jsonParser) skipSpace() error {
	if s.offset >= len(s.buf) {
		return nil
	}
	n := skipSpace(s.buf[s.offset:])
	if err := s.incOffset(n); err != nil {
		return err
	}
	return nil
}

func (s *jsonParser) eof() error {
	if err := s.skipSpace(); err != nil {
		return err
	}
	s.nextErr = io.EOF
	return io.EOF
}

func (s *jsonParser) scanOpenObject() error {
	if !s.isNext('{') {
		return errExpectedOpenCurly
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanOpenArray() error {
	if !s.isNext('[') {
		return errExpectedOpenBracket
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanString() error {
	n, err := scanString(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanConst(valBytes []byte, err error) error {
	start := s.offset
	if orr := s.incOffset(len(valBytes)); orr != nil {
		return err
	}
	end := s.offset
	if !bytes.Equal(s.buf[start:end], valBytes) {
		return err
	}
	return s.skipSpace()
}

var trueBytes = []byte{'t', 'r', 'u', 'e'}

func (s *jsonParser) scanTrue() error {
	return s.scanConst(trueBytes, errExpectedTrue)
}

var falseBytes = []byte{'f', 'a', 'l', 's', 'e'}

func (s *jsonParser) scanFalse() error {
	return s.scanConst(falseBytes, errExpectedFalse)
}

var nullBytes = []byte{'n', 'u', 'l', 'l'}

func (s *jsonParser) scanNull() error {
	return s.scanConst(nullBytes, errExpectedNull)
}

func (s *jsonParser) scanDigits() error {
	if s.offset >= len(s.buf) {
		return io.ErrShortBuffer
	}
	for s.isNextDigit() {
		if err := s.incOffset(1); err != nil {
			return err
		}
	}
	return nil
}

func (s *jsonParser) scanNumber() error {
	if s.isNext('-') {
		if err := s.incOffset(1); err != nil {
			return err
		}
	}
	if s.isNext('0') {
		if err := s.incOffset(1); err != nil {
			return err
		}
	} else if s.isNextDigit19() {
		if err := s.scanDigits(); err != nil {
			return err
		}
	}
	if s.isNext('.') {
		if err := s.incOffset(1); err != nil {
			return err
		}
		if err := s.scanDigits(); err != nil {
			return err
		}
	}
	if s.isNextOneOf('e', 'E') {
		if err := s.incOffset(1); err != nil {
			return err
		}
		if s.isNextOneOf('+', '-') {
			if err := s.incOffset(1); err != nil {
				return err
			}
		}
		if err := s.scanDigits(); err != nil {
			return err
		}
	}
	return nil
}

func (s *jsonParser) scanColon() error {
	if !s.isNext(':') {
		return errExpectedColon
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanCloseObject() error {
	if !s.isNext('}') {
		return errExpectedCloseCurly
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.eof()
}

func (s *jsonParser) scanCloseArray() error {
	if !s.isNext(']') {
		return errExpectedCloseBracket
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.eof()
}

func (s *jsonParser) scanComma() error {
	if !s.isNext(',') {
		return errExpectedComma
	}
	if err := s.incOffset(1); err != nil {
		return err
	}
	return s.skipSpace()
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

func (s *jsonParser) IsLeaf() bool {
	switch s.kind {
	case stringKind, numberKind, trueKind, falseKind, nullKind:
		return true
	}
	return false
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

func (s *jsonParser) String() (string, error) {
	if !s.kind.isString() && !s.kind.isObject() {
		return "", parser.ErrNotString
	}
	if err := s.parse(); err != nil {
		return "", err
	}
	return s.parsedString, nil
}

func (s *jsonParser) parseString(buf []byte) error {
	res, err := unquote(s.pool, buf)
	if err != nil {
		return err
	}
	s.parsedString = res
	return nil
}

func (s *jsonParser) parseArrayIndex() error {
	s.parsedKindOfNumber = intOfNumber
	s.parsedInt = int64(s.arrayIndex)
	return nil
}

func (s *jsonParser) parseNumber(buf []byte) error {
	var err error
	s.parsedDouble, err = strconv.ParseFloat(buf)
	if err != nil {
		s.parsedKindOfNumber = noneOfNumber
		// scan already passed, so we know this is a valid number.
		// The number is just too large represent in a float.
		return nil
	}
	s.parsedUint = uint64(s.parsedDouble)
	isUint := float64(s.parsedUint) == s.parsedDouble
	s.parsedInt = int64(s.parsedDouble)
	isInt := float64(s.parsedInt) == s.parsedDouble
	if isUint && isInt {
		s.parsedKindOfNumber = anyOfNumber
		return nil
	}
	if isInt {
		s.parsedKindOfNumber = intOfNumber
		return nil
	}
	if isUint {
		s.parsedKindOfNumber = uintOfNumber
		return nil
	}
	s.parsedKindOfNumber = doubleOfNumber
	return nil
}

func (s *jsonParser) parse() error {
	scanned, err := s.scan()
	if err != nil {
		return err
	}
	if !s.parsed {
		var err error
		switch s.kind {
		case objectKind:
			err = s.parseString(scanned)
		case arrayKind:
			err = s.parseArrayIndex()
		case stringKind:
			err = s.parseString(scanned)
		case numberKind:
			err = s.parseNumber(scanned)
		case trueKind, falseKind, nullKind:
			// do nothing, scanning is enough
		}
		s.parsed = true
		if err != nil {
			s.parsedErr = err
			return err
		}
	}
	return s.parsedErr
}

func (s *jsonParser) skip() {
	s.Down()
	s.Up()
}

func (s *jsonParser) scan() ([]byte, error) {
	if s.scanned == nil && s.scannedErr == nil {
		if err := s.skipSpace(); err != nil {
			return nil, err
		}
		start := s.offset
		switch s.kind {
		case objectKind:
			if err := s.scanString(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		case arrayKind:
			return nil, nil
		case stringKind:
			if err := s.scanString(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		case numberKind:
			if err := s.scanNumber(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		case trueKind:
			if err := s.scanTrue(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		case falseKind:
			if err := s.scanFalse(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		case nullKind:
			if err := s.scanNull(); err != nil {
				s.scannedErr = err
				return nil, err
			}
		}
		end := s.offset
		s.scanned = s.buf[start:end]
		if err := s.skipSpace(); err != nil {
			return nil, err
		}
	}
	return s.scanned, s.scannedErr
}

func (s *jsonParser) Bytes() ([]byte, error) {
	scanned, err := s.scan()
	if err != nil {
		return nil, err
	}
	return scanned, nil
}

// JsonParser is a parser for JSON
type JsonParser interface {
	parser.Interface
	//Init initialises the parser with a byte buffer containing JSON.
	Init(buf []byte) error
	Kind() Kind
}

// NewJsonParser returns a new JSON parser.
func NewJsonParser() JsonParser {
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
