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

package strconv

import (
	"errors"
	"strings"
)

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat, ParseComplex)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string {
	return "strconv." + e.Func + ": " + "parsing " + Quote(e.Num) + ": " + e.Err.Error()
}

func (e *NumError) Unwrap() error { return e.Err }

func syntaxError(fn, str string) *NumError {
	return &NumError{fn, strings.Clone(str), ErrSyntax}
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, strings.Clone(str), ErrRange}
}

func BaseError(fn, str string, base int) *NumError {
	return &NumError{fn, strings.Clone(str), errors.New("invalid base " + Itoa(base))}
}

func BitSizeError(fn, str string, bitSize int) *NumError {
	return &NumError{fn, strings.Clone(str), errors.New("invalid bit size " + Itoa(bitSize))}
}
