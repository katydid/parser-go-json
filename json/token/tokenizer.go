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
	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)

// Tokenizer is a scanner that provides the ability to buffers returned by the scanner into native Go types.
type Tokenizer interface {
	// Next returns the Kind of the token or an error.
	Next() (scan.Kind, error)
	// Token parses and returns the current token.
	Token() (parse.Kind, []byte, error)
	// Init restarts the tokenizer with a new byte buffer, without allocating a new tokenizer.
	Init([]byte)
}

type tokenizer struct {
	scanner scan.Scanner
	alloc   func(size int) []byte

	scanToken []byte
	scanKind  scan.Kind

	tokenized   bool
	tokenKind   parse.Kind
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
			t.tokenKind = parse.DecimalKind
			t.tokenBytes = t.scanToken
			// scan already passed, so we know this is a valid number.
			// The number is just too large represent in 64 float bits.
			return nil
		}
		t.tokenKind = parse.Float64Kind
		// This can only be a float, so we return and do not try others.
		return nil
	}
	parsedInt, err := strconv.ParseInt(t.scanToken)
	if err != nil {
		// scan already passed, so we know this is a valid number.
		// The number is just too large represent in signed 64 bits.
		t.tokenKind = parse.DecimalKind
		t.tokenBytes = t.scanToken
		return nil
	}
	t.tokenKind = parse.Int64Kind
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
	t.tokenKind = parse.StringKind
	return nil
}

func (t *tokenizer) Token() (parse.Kind, []byte, error) {
	if err := t.tokenize(); err != nil {
		return parse.UnknownKind, nil, err
	}
	if t.tokenKind == parse.Int64Kind {
		return t.tokenKind, cast.FromInt64(t.tokenInt, t.alloc), nil
	}
	if t.tokenKind == parse.Float64Kind {
		return t.tokenKind, cast.FromFloat64(t.tokenDouble, t.alloc), nil
	}
	return t.tokenKind, t.tokenBytes, nil
}

func (t *tokenizer) Tokenize() (parse.Kind, error) {
	if err := t.tokenize(); err != nil {
		return parse.UnknownKind, err
	}
	return t.tokenKind, nil
}

func (t *tokenizer) Int() (int64, error) {
	if t.tokenKind == parse.Int64Kind {
		return t.tokenInt, nil
	}
	return 0, ErrNotInt
}

func (t *tokenizer) Double() (float64, error) {
	if t.tokenKind == parse.Float64Kind {
		return t.tokenDouble, nil
	}
	return 0, ErrNotDouble
}

func (t *tokenizer) Bytes() ([]byte, error) {
	if t.tokenKind == parse.BytesKind || t.tokenKind == parse.StringKind || t.tokenKind == parse.DecimalKind {
		return t.tokenBytes, nil
	}
	return nil, ErrNotBytes
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
			t.tokenKind = parse.TrueKind
		case scan.FalseKind:
			t.tokenKind = parse.FalseKind
		case scan.NullKind:
			t.tokenKind = parse.NullKind
		}
		t.tokenized = true
		if err != nil {
			t.tokenErr = err
			return err
		}
	}
	return t.tokenErr
}
