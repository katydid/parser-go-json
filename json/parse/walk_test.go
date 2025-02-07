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

package parse

import (
	"errors"
	"io"
)

var errUnknownToken = errors.New("unknown token")

var errExpectedBool = errors.New("expected bool")

var errExpectedString = errors.New("expected string")

func walkValue(t Parser, kind Kind) error {
	if _, err := t.Bool(); err == nil {
		return nil
	}
	if kind == BoolKind {
		return errExpectedBool
	}
	if _, err := t.Int(); err == nil {
		return nil
	}
	if _, err := t.Uint(); err == nil {
		return nil
	}
	if _, err := t.Double(); err == nil {
		return nil
	}
	if _, err := t.String(); err == nil {
		return nil
	}
	if kind == StringKind {
		return errExpectedString
	}
	if _, err := t.Bytes(); err == nil {
		return nil
	}
	return errUnknownToken
}

func walk(t Parser) error {
	kind, err := t.Next()
	for err == nil {
		switch kind {
		case NullKind, BoolKind, NumberKind, StringKind:
			if err := walkValue(t, kind); err != nil {
				return err
			}
		}
		kind, err = t.Next()
	}
	if err != io.EOF {
		return err
	}
	return nil
}
