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

package json

import (
	"github.com/katydid/parser-go-json/json/internal/fork/strconv"
	"github.com/katydid/parser-go-json/json/internal/fork/unquote"
	"github.com/katydid/parser-go-json/json/internal/pool"
)

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

func unquoteBytes(pool pool.Pool, s []byte) (string, error) {
	var ok bool
	var t string
	s, ok = unquote.Unquote(pool.Alloc, s)
	t = castToString(s)
	if !ok {
		return "", errUnquote
	}
	return t, nil
}

func (s *jsonParser) parseString(buf []byte) error {
	res, err := unquoteBytes(s.pool, buf)
	if err != nil {
		return err
	}
	s.parsedString = res
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
