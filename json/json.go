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
	"strconv"

	"github.com/katydid/parser-go/parser"
)

func isString(buf []byte) bool {
	return buf[0] == '"'
}

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

func isNumber(c byte) bool {
	return (c == '-') || ((c >= '0') && (c <= '9'))
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

func unquote(s []byte) (string, error) {
	var ok bool
	var t string
	s, ok = unquoteBytes(s)
	t = string(s)
	if !ok {
		return "", errUnquote
	}
	return t, nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isDigit19(c byte) bool {
	return c >= '1' && c <= '9'
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

func (s *jsonParser) skipSpace() error {
	n := skipSpace(s.buf[s.offset:])
	if err := s.incOffset(n); err != nil {
		return err
	}
	return nil
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
	s.startValueOffset = s.offset
	n, err := scanString(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	s.endValueOffset = s.offset
	return s.skipSpace()
}

func (s *jsonParser) scanName() error {
	startOffset := s.offset
	n, err := scanString(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	s.name, err = unquote(s.buf[startOffset:s.offset])
	if err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanConst(valBytes []byte, err error) error {
	s.startValueOffset = s.offset
	if orr := s.incOffset(len(valBytes)); orr != nil {
		return err
	}
	s.endValueOffset = s.offset
	if !bytes.Equal(s.buf[s.startValueOffset:s.offset], valBytes) {
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

func (s *jsonParser) scanArray() error {
	count := 0
	index := 0
	for i, c := range s.buf[s.offset:] {
		if c == '[' {
			count++
		}
		if c == ']' {
			count--
		}
		if count == 0 {
			index = i
			break
		}
	}
	if count != 0 {
		return errExpectedCloseBracket
	}
	s.startValueOffset = s.offset
	s.endValueOffset = s.offset + index + 1
	if err := s.incOffset(index + 1); err != nil {
		return err
	}
	s.isValueArray = true
	return s.skipSpace()
}

func (s *jsonParser) scanObject() error {
	count := 0
	index := 0
	for i, c := range s.buf[s.offset:] {
		if c == '{' {
			count++
		}
		if c == '}' {
			count--
		}
		if count == 0 {
			index = i
			break
		}
	}
	if count != 0 {
		return errExpectedCloseCurly
	}
	s.startValueOffset = s.offset
	s.endValueOffset = s.offset + index + 1
	if err := s.incOffset(index + 1); err != nil {
		return err
	}
	s.isValueObject = true
	return s.skipSpace()
}

func (s *jsonParser) scanNumber() error {
	s.startValueOffset = s.offset
	if s.isNext('-') {
		if err := s.incOffset(1); err != nil {
			return err
		}
	}
	if s.isNext('0') {
		if err := s.incOffset(1); err != nil {
			return err
		}
	} else if isDigit19(s.buf[s.offset]) {
		if err := s.incOffset(1); err != nil {
			return err
		}
		for s.offset < len(s.buf) && isDigit(s.buf[s.offset]) {
			if err := s.incOffset(1); err != nil {
				return err
			}
		}
	}
	if s.isNext('.') {
		if err := s.incOffset(1); err != nil {
			return err
		}
		for s.offset < len(s.buf) && isDigit(s.buf[s.offset]) {
			if err := s.incOffset(1); err != nil {
				return err
			}
		}
	}
	if s.offset < len(s.buf) &&
		(s.buf[s.offset] == 'e' || s.buf[s.offset] == 'E') {
		if err := s.incOffset(1); err != nil {
			return err
		}
		if s.offset < len(s.buf) {
			if s.buf[s.offset] == '+' || s.buf[s.offset] == '-' {
				if err := s.incOffset(1); err != nil {
					return err
				}
			}
		}
		for s.offset < len(s.buf) && isDigit(s.buf[s.offset]) {
			if err := s.incOffset(1); err != nil {
				return err
			}
		}
	}
	s.endValueOffset = s.offset
	return nil
}

func (s *jsonParser) scanValue() error {
	c := s.buf[s.offset]
	if isNumber(c) {
		return s.scanNumber()
	}
	switch c {
	case '"':
		return s.scanString()
	case '{':
		return s.scanObject()
	case '[':
		return s.scanArray()
	case 't':
		return s.scanTrue()
	case 'f':
		return s.scanFalse()
	case 'n':
		return s.scanNull()
	}
	return errExpectedValue
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
	return io.EOF
}

func (s *jsonParser) scanCloseArray() error {
	if !s.isNext(']') {
		return errExpectedCloseBracket
	}
	return io.EOF
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

func (s *jsonParser) nextValueInArray() error {
	if s.firstArrayValue {
		if err := s.scanOpenArray(); err != nil {
			return err
		}
		s.firstArrayValue = false
	} else {
		if s.isNext(',') {
			if err := s.scanComma(); err != nil {
				return err
			}
		} else {
			return s.scanCloseArray()
		}
	}
	if s.isNext(']') {
		return s.scanCloseArray()
	}
	return s.scanValue()
}

func (s *jsonParser) nextValueInObject() error {
	if s.firstObjectValue {
		if err := s.scanOpenObject(); err != nil {
			return err
		}
		s.firstObjectValue = false
	} else {
		if s.isNext(',') {
			if err := s.scanComma(); err != nil {
				return err
			}
		} else {
			if err := s.scanCloseObject(); err != nil {
				return err
			}
		}
	}
	if s.offset >= len(s.buf) {
		return io.ErrShortBuffer
	}
	if isString(s.buf[s.offset:]) {
		if err := s.scanName(); err != nil {
			return err
		}
		if err := s.scanColon(); err != nil {
			return err
		}
		if err := s.scanValue(); err != nil {
			return err
		}
		return nil
	}
	return s.scanCloseObject()
}

func (s *jsonParser) Next() error {
	if s.isLeaf {
		if s.firstObjectValue {
			s.firstObjectValue = false
			return nil
		}
		return io.EOF
	}
	s.isValueObject = false
	s.isValueArray = false
	if err := s.skipSpace(); err != nil {
		return err
	}
	if s.inArray {
		if !s.firstArrayValue {
			s.arrayIndex++
		}
		return s.nextValueInArray()
	}
	return s.nextValueInObject()
}

func (s *jsonParser) IsLeaf() bool {
	return s.isLeaf
}

func (s *jsonParser) value() []byte {
	return s.buf[s.startValueOffset:s.endValueOffset]
}

func (s *jsonParser) Double() (float64, error) {
	if s.isLeaf {
		v := string(s.value())
		i, err := strconv.ParseFloat(v, 64)
		return i, err
	}
	return 0, parser.ErrNotDouble
}

func (s *jsonParser) Int() (int64, error) {
	if s.isLeaf {
		v := string(s.value())
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			f, ferr := strconv.ParseFloat(v, 64)
			if ferr != nil {
				return i, err
			}
			if float64(int64(f)) == f {
				return int64(f), nil
			}
		}
		return i, err
	}
	if s.inArray {
		return int64(s.arrayIndex), nil
	}
	return 0, parser.ErrNotInt
}

func (s *jsonParser) Uint() (uint64, error) {
	if s.isLeaf {
		v := string(s.value())
		i, err := strconv.ParseUint(v, 10, 64)
		return uint64(i), err
	}
	return 0, parser.ErrNotUint
}

func (s *jsonParser) Bool() (bool, error) {
	if s.isLeaf {
		v := s.value()
		if bytes.Equal(v, trueBytes) {
			return true, nil
		}
		if bytes.Equal(v, falseBytes) {
			return false, nil
		}
	}
	return false, parser.ErrNotBool
}

func (s *jsonParser) String() (string, error) {
	if s.isLeaf {
		v := s.value()
		if v[0] != '"' {
			return "", parser.ErrNotString
		}
		res, err := unquote(v)
		if err != nil {
			return "", err
		}
		return res, nil
	}
	if s.inArray {
		return "", parser.ErrNotString
	}
	return s.name, nil
}

func (s *jsonParser) Bytes() ([]byte, error) {
	return s.value(), nil
}

// JsonParser is a parser for JSON
type JsonParser interface {
	parser.Interface
	//Init initialises the parser with a byte buffer containing JSON.
	Init(buf []byte) error
	Reset() error
}

// NewJsonParser returns a new JSON parser.
func NewJsonParser() JsonParser {
	return &jsonParser{
		state: state{
			firstObjectValue: true,
		},
		stack: make([]state, 0, 10),
	}
}

func (s *jsonParser) Init(buf []byte) error {
	s.state = state{
		firstObjectValue: true,
		buf:              buf,
	}
	s.stack = s.stack[:0]
	if err := s.skipSpace(); err != nil {
		return err
	}
	if s.offset >= len(s.buf) {
		return io.ErrShortBuffer
	}
	if s.isNext('{') {
		//do nothing
	} else if s.isNext('[') {
		if err := s.scanValue(); err != nil {
			return err
		}
		s.inArray = true
		s.firstArrayValue = true
		s.buf = s.buf[s.startValueOffset:s.endValueOffset]
		s.offset = 0
	} else {
		if err := s.scanValue(); err != nil {
			return err
		}
		s.state.isLeaf = true
		s.state.firstObjectValue = true
	}
	return nil
}

func (s *jsonParser) Reset() error {
	if len(s.stack) > 0 {
		return s.Init(s.stack[0].buf)
	}
	return s.Init(s.buf)
}

type jsonParser struct {
	state
	stack []state
}

type state struct {
	buf              []byte
	offset           int
	name             string
	startValueOffset int
	endValueOffset   int
	inArray          bool
	firstObjectValue bool
	firstArrayValue  bool
	isValueObject    bool
	isValueArray     bool
	isLeaf           bool
	arrayIndex       int
}

func (s *jsonParser) Up() {
	top := len(s.stack) - 1
	s.state = s.stack[top]
	s.stack = s.stack[:top]
}

func (s *jsonParser) Down() {
	if s.isValueObject {
		s.stack = append(s.stack, s.state)
		s.state = state{
			buf:              s.buf[s.startValueOffset:s.endValueOffset],
			firstObjectValue: true,
		}
	} else if s.isValueArray {
		s.stack = append(s.stack, s.state)
		s.state = state{
			buf:             s.buf[s.startValueOffset:s.endValueOffset],
			firstArrayValue: true,
			inArray:         true,
		}
	} else {
		s.stack = append(s.stack, s.state)
		s.state.isLeaf = true
		s.state.firstObjectValue = true
	}
}
