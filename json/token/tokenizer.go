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

	scanTokenStart []byte
	skipped        bool
	scanKind       scan.Kind

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
		skipped: true,
	}
}

// Init restarts the tokenizer with a new byte buffer, without allocating a new tokenizer.
func (t *tokenizer) Init(buf []byte) {
	t.skipped = true
	t.scanner.Init(buf)
}

// Next returns the Kind of the token or an error.
func (t *tokenizer) Next() (scan.Kind, error) {
	if !t.skipped {
		if _, err := t.scanner.ScanToEnd(t.scanKind); err != nil {
			return scan.UnknownKind, nil
		}
	}
	t.tokenized = false
	kind, token, err := t.scanner.NextStart()
	if err != nil {
		return scan.UnknownKind, err
	}
	t.skipped = false
	t.scanKind = kind
	t.scanTokenStart = token
	return kind, nil
}

func (t *tokenizer) tokenizeNumber() error {
	offset, intval, intok, floatval, floatok, decimalok := scan.ParseNumber(t.scanTokenStart)
	if err := t.scanner.Skip(offset); err != nil {
		return err
	}
	t.skipped = true
	if intok {
		t.tokenKind = parse.Int64Kind
		t.tokenInt = intval
		return nil
	}
	if floatok {
		t.tokenKind = parse.Float64Kind
		t.tokenDouble = floatval
		return nil
	}
	if decimalok {
		t.tokenKind = parse.DecimalKind
		t.tokenBytes = t.scanTokenStart[:offset]
		return nil
	}
	return ErrNotNumber
}

func unquoteBytes(alloc func(int) []byte, s []byte) ([]byte, int, error) {
	u, offset, ok := unquote.Unquote(alloc, s)
	if !ok {
		return nil, 0, errUnquote
	}
	return u, offset, nil
}

func (t *tokenizer) tokenizeString() error {
	res, offset, err := unquoteBytes(t.alloc, t.scanTokenStart)
	if err != nil {
		return err
	}
	if err := t.scanner.Skip(offset); err != nil {
		return err
	}
	t.skipped = true
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
