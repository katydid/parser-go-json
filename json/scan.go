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

import "github.com/katydid/parser-go-json/json/scan"

func (s *jsonParser) skipSpace() error {
	if s.offset >= len(s.buf) {
		return nil
	}
	n := scan.Space(s.buf[s.offset:])
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

func (s *jsonParser) scanNull() error {
	n, err := scan.Null(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanFalse() error {
	n, err := scan.False(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanTrue() error {
	n, err := scan.True(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanNumber() error {
	n, err := scan.Number(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

func (s *jsonParser) scanString() error {
	n, err := scan.String(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}
