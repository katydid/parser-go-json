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

package token

import (
	"github.com/katydid/parser-go-json/json/internal/fork/strconv"
	"github.com/katydid/parser-go-json/json/internal/fork/unquote"
	"github.com/katydid/parser-go-json/json/scan"
)

// Tokenizer is a scanner that provides the ability to buffers returned by the scanner into native Go types.
type Tokenizer interface {
	// Next returns the Kind of the token or an error.
	Next() (scan.Kind, error)
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
	// Init restarts the tokenizer with a new byte buffer, without allocating a new tokenizer.
	Init([]byte)
}

type tokenizer struct {
	scanner scan.Scanner
	alloc   func(size int) []byte

	scanToken []byte
	scanKind  scan.Kind

	tokenized   bool
	tokenKind   kind
	tokenErr    error
	tokenDouble float64
	tokenInt    int64
	tokenUint   uint64
	tokenString string
}

func NewTokenizer(buf []byte) Tokenizer {
	alloc := func(size int) []byte {
		return make([]byte, size)
	}
	return NewTokenizerWithCustomAllocator(buf, alloc)
}

func NewTokenizerWithCustomAllocator(buf []byte, alloc func(int) []byte) Tokenizer {
	return &tokenizer{
		scanner: scan.NewScanner(buf),
		alloc:   alloc,
	}
}

// Init restarts the tokenizer with a new byte buffer, without allocating a new tokenizer.
func (t *tokenizer) Init(buf []byte) {
	t.scanner.Init(buf)
}

// Next returns the Kind of the token or an error.
func (t *tokenizer) Next() (scan.Kind, error) {
	t.tokenized = false
	kind, token, err := t.scanner.Next()
	if err != nil {
		return scan.UnknownKind, err
	}
	t.scanKind = kind
	t.scanToken = token
	return kind, nil
}

// Bool attempts to convert the current token to a bool.
func (t *tokenizer) Bool() (bool, error) {
	if !t.scanKind.IsTrue() && !t.scanKind.IsFalse() {
		return false, ErrNotBool
	}
	if err := t.tokenize(); err != nil {
		return false, err
	}
	if t.tokenKind.IsBool() && t.scanKind.IsTrue() {
		return true, nil
	}
	if t.tokenKind.IsBool() && t.scanKind.IsFalse() {
		return false, nil
	}
	return false, ErrNotBool
}

// Int attempts to convert the current token to an int64.
func (t *tokenizer) Int() (int64, error) {
	if !t.scanKind.IsNumber() {
		return 0, ErrNotInt
	}
	if err := t.tokenize(); err != nil {
		return 0, err
	}
	if !t.tokenKind.IsInt() {
		return 0, ErrNotInt
	}
	return t.tokenInt, nil
}

// Uint attempts to convert the current token to an uint64.
func (t *tokenizer) Uint() (uint64, error) {
	if !t.scanKind.IsNumber() {
		return 0, ErrNotUint
	}
	if err := t.tokenize(); err != nil {
		return 0, err
	}
	if !t.tokenKind.IsUint() {
		return 0, ErrNotUint
	}
	return t.tokenUint, nil
}

// Double attempts to convert the current token to a float64.
func (t *tokenizer) Double() (float64, error) {
	if !t.scanKind.IsNumber() {
		return 0, ErrNotDouble
	}
	if err := t.tokenize(); err != nil {
		return 0, err
	}
	if !t.tokenKind.IsDouble() {
		return 0, ErrNotDouble
	}
	return t.tokenDouble, nil
}

// String attempts to convert the current token to a string.
func (t *tokenizer) String() (string, error) {
	if !t.scanKind.IsString() {
		return "", ErrNotString
	}
	if err := t.tokenize(); err != nil {
		return "", err
	}
	if !t.tokenKind.IsString() {
		return "", ErrNotString
	}
	return t.tokenString, nil
}

// Bytes returns the raw current token.
func (t *tokenizer) Bytes() ([]byte, error) {
	return t.scanToken, nil
}

func (t *tokenizer) notInteger() bool {
	for _, b := range t.scanToken {
		if b == '.' || b == 'e' || b == 'E' {
			return true
		}
	}
	return false
}

func (t *tokenizer) tokenizeNumber() error {
	var err error
	if t.notInteger() {
		t.tokenDouble, err = strconv.ParseFloat(t.scanToken)
		if err != nil {
			t.tokenKind = TooLargeNumberKind
			// scan already passed, so we know this is a valid number.
			// The number is just too large represent in 64 bits.
			return nil
		}
		t.tokenKind = FractionNumberKind
		// This can only be a float, so we return and do not try others.
		return nil
	}
	t.tokenInt, err = strconv.ParseInt(t.scanToken)
	if err == nil {
		if t.scanToken[0] == '-' {
			t.tokenKind = NegativeNumberKind
		} else {
			t.tokenKind = NumberKind
			t.tokenUint = uint64(t.tokenInt)
			t.tokenDouble = float64(t.tokenDouble)
		}
		return nil
	}
	// scan already passed, so we know this is a valid number.
	// The number is just too large represent in signed 64 bits.
	t.tokenKind = TooLargeNumberKind
	// This can be overwritten if uint parses it correctly, so we do not return yet.
	// Only if int could not parse the non negative number, tokenKind == TooLargeNumberKind then try uint64 too
	if t.scanToken[0] != '-' {
		t.tokenUint, err = strconv.ParseUint(t.scanToken)
		if err != nil {
			t.tokenKind = TooLargeNumberKind
			// scan already passed, so we know this is a valid number.
			// The number is just too large represent in 64 bits.
		} else {
			t.tokenKind = LargePositiveNumberKind
		}
	}
	return nil
}

func unquoteBytes(alloc func(int) []byte, s []byte) (string, error) {
	var ok bool
	var t string
	s, ok = unquote.Unquote(alloc, s)
	t = castToString(s)
	if !ok {
		return "", errUnquote
	}
	return t, nil
}

func (t *tokenizer) tokenizeString() error {
	res, err := unquoteBytes(t.alloc, t.scanToken)
	if err != nil {
		return err
	}
	t.tokenString = res
	t.tokenKind = StringKind
	return nil
}

func (t *tokenizer) tokenize() error {
	if !t.tokenized {
		var err error
		switch t.scanKind {
		case scan.StringKind:
			err = t.tokenizeString()
		case scan.NumberKind:
			err = t.tokenizeNumber()
		case scan.TrueKind:
			t.tokenKind = BoolKind
		case scan.FalseKind:
			t.tokenKind = BoolKind
		case scan.NullKind:
			t.tokenKind = NullKind
		}
		t.tokenized = true
		if err != nil {
			t.tokenErr = err
			return err
		}
	}
	return t.tokenErr
}
