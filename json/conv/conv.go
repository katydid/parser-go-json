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

package conv

import (
	"github.com/katydid/parser-go-json/json/internal/fork/strconv"
	"github.com/katydid/parser-go-json/json/internal/fork/unquote"
	"github.com/katydid/parser-go-json/json/scan"
)

// Converter is a scanner that provides the ability to convert tokens to native Go types.
type Converter interface {
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
	// Restart the converter with a new byte buffer, without allocating a new converter.
	Restart([]byte)
}

type converter struct {
	scanner scan.Scanner
	alloc   func(size int) []byte

	scanToken []byte
	scanKind  scan.Kind

	converted  bool
	convKind   Kind
	convErr    error
	convDouble float64
	convInt    int64
	convUint   uint64
	convString string
}

func NewConverter(scanner scan.Scanner) Converter {
	alloc := func(size int) []byte {
		return make([]byte, size)
	}
	return NewConverterWithCustomAllocator(scanner, alloc)
}

func NewConverterWithCustomAllocator(scanner scan.Scanner, alloc func(int) []byte) Converter {
	return &converter{
		scanner: scanner,
		alloc:   alloc,
	}
}

// Restart the scanner with a new byte buffer, without allocating a new scanner.
func (c *converter) Restart(buf []byte) {
	c.scanner.Restart(buf)
}

// Next returns the Kind of the token or an error.
func (c *converter) Next() (scan.Kind, error) {
	c.converted = false
	kind, token, err := c.scanner.Next()
	if err != nil {
		return scan.UnknownKind, err
	}
	c.scanKind = kind
	c.scanToken = token
	return kind, nil
}

// Bool attempts to convert the current token to a bool.
func (c *converter) Bool() (bool, error) {
	if !c.scanKind.IsTrue() || c.scanKind.IsFalse() {
		return false, ErrNotBool
	}
	if err := c.convert(); err != nil {
		return false, err
	}
	if c.convKind.IsTrue() {
		return true, nil
	}
	if c.convKind.IsFalse() {
		return false, nil
	}
	return false, ErrNotBool
}

// Int attempts to convert the current token to an int64.
func (c *converter) Int() (int64, error) {
	if !c.scanKind.IsNumber() {
		return 0, ErrNotInt
	}
	if err := c.convert(); err != nil {
		return 0, err
	}
	if !c.convKind.IsInt() {
		return 0, ErrNotInt
	}
	return c.convInt, nil
}

// Uint attempts to convert the current token to an uint64.
func (c *converter) Uint() (uint64, error) {
	if !c.scanKind.IsNumber() {
		return 0, ErrNotUint
	}
	if err := c.convert(); err != nil {
		return 0, err
	}
	if !c.convKind.IsUint() {
		return 0, ErrNotUint
	}
	return c.convUint, nil
}

// Double attempts to convert the current token to a float64.
func (c *converter) Double() (float64, error) {
	if !c.scanKind.IsNumber() {
		return 0, ErrNotDouble
	}
	if err := c.convert(); err != nil {
		return 0, err
	}
	if !c.convKind.IsDouble() {
		return 0, ErrNotDouble
	}
	return c.convDouble, nil
}

// String attempts to convert the current token to a string.
func (c *converter) String() (string, error) {
	if !c.scanKind.IsString() {
		return "", ErrNotString
	}
	if err := c.convert(); err != nil {
		return "", err
	}
	if !c.convKind.IsString() {
		return "", ErrNotString
	}
	return c.convString, nil
}

// Bytes returns the raw current token.
func (c *converter) Bytes() ([]byte, error) {
	return c.scanToken, nil
}

func (c *converter) convertNumber() error {
	var err error
	c.convDouble, err = strconv.ParseFloat(c.scanToken)
	if err != nil {
		c.convKind = TooLargeNumberKind
		// scan already passed, so we know this is a valid number.
		// The number is just too large represent in a float.
		return nil
	}
	c.convUint = uint64(c.convDouble)
	isUint := float64(c.convUint) == c.convDouble
	c.convInt = int64(c.convDouble)
	isInt := float64(c.convInt) == c.convDouble
	if isUint && isInt {
		c.convKind = NumberKind
		return nil
	}
	if isInt {
		c.convKind = NegativeNumberKind
		return nil
	}
	if isUint {
		c.convKind = LargePositiveNumberKind
		return nil
	}
	c.convKind = FractionKind
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

func (c *converter) convertString() error {
	res, err := unquoteBytes(c.alloc, c.scanToken)
	if err != nil {
		return err
	}
	c.convString = res
	return nil
}

func (c *converter) convert() error {
	if !c.converted {
		var err error
		switch c.scanKind {
		case scan.StringKind:
			err = c.convertString()
		case scan.NumberKind:
			err = c.convertNumber()
		case scan.TrueKind:
			c.convKind = TrueKind
		case scan.FalseKind:
			c.convKind = FalseKind
		case scan.NullKind:
			c.convKind = NullKind
		}
		c.converted = true
		if err != nil {
			c.convErr = err
			return err
		}
	}
	return c.convErr
}
