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
	// Tokenize parses the current token.
	Tokenize() (Kind, error)
	// Int attempts to convert the current token to an int64.
	Int() (int64, error)
	// Double attempts to convert the current token to a float64.
	Double() (float64, error)
	// Bytes returns the bytes token or a unquoted string or decimal.
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
	tokenKind   Kind
	tokenErr    error
	tokenDouble float64
	tokenInt    int64
	tokenBytes  []byte
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

// Int attempts to convert the current token to an int64.
func (t *tokenizer) Int() (int64, error) {
	if !t.scanKind.IsNumber() {
		return 0, ErrNotInt
	}
	bs, err := t.Bytes()
	if err != nil {
		return 0, err
	}
	if t.tokenKind == Int64Kind {
		return castToInt64(bs), nil
	}
	return 0, ErrNotInt
}

// Double attempts to convert the current token to a float64.
func (t *tokenizer) Double() (float64, error) {
	if !t.scanKind.IsNumber() {
		return 0, ErrNotDouble
	}
	if err := t.tokenize(); err != nil {
		return 0, err
	}
	if !t.tokenKind.IsFloat64() {
		return 0, ErrNotDouble
	}
	return t.tokenDouble, nil
}

// Bytes returns the raw current token.
func (t *tokenizer) Bytes() ([]byte, error) {
	if err := t.tokenize(); err != nil {
		return nil, err
	}
	if t.tokenKind == Int64Kind {
		return deprecatedCastFromInt64(t.tokenInt), nil
	}
	if t.tokenKind != BytesKind && t.tokenKind != StringKind && t.tokenKind != DecimalKind {
		return nil, ErrNotBytes
	}
	return t.tokenBytes, nil
}

func (t *tokenizer) notParseableInteger() bool {
	for _, b := range t.scanToken {
		if b == '.' || b == 'e' || b == 'E' {
			return true
		}
	}
	return false
}

func (t *tokenizer) tokenizeNumber() error {
	var err error
	if t.notParseableInteger() {
		t.tokenDouble, err = strconv.ParseFloat(t.scanToken)
		if err != nil {
			t.tokenKind = DecimalKind
			t.tokenBytes = t.scanToken
			// scan already passed, so we know this is a valid number.
			// The number is just too large represent in 64 float bits.
			return nil
		}
		t.tokenKind = Float64Kind
		// This can only be a float, so we return and do not try others.
		return nil
	}
	parsedInt, err := strconv.ParseInt(t.scanToken)
	if err != nil {
		// scan already passed, so we know this is a valid number.
		// The number is just too large represent in signed 64 bits.
		t.tokenKind = DecimalKind
		t.tokenBytes = t.scanToken
		return nil
	}
	t.tokenKind = Int64Kind
	t.tokenInt = parsedInt
	return nil
}

func unquoteBytes(alloc func(int) []byte, s []byte) ([]byte, error) {
	var ok bool
	var u []byte
	u, ok = unquote.Unquote(alloc, s)
	if !ok {
		return nil, errUnquote
	}
	return u, nil
}

func (t *tokenizer) tokenizeString() error {
	res, err := unquoteBytes(t.alloc, t.scanToken)
	if err != nil {
		return err
	}
	t.tokenBytes = res
	t.tokenKind = StringKind
	return nil
}

func (t *tokenizer) Tokenize() (Kind, error) {
	if err := t.tokenize(); err != nil {
		return 0, err
	}
	return t.tokenKind, nil
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
			t.tokenKind = TrueKind
		case scan.FalseKind:
			t.tokenKind = FalseKind
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
