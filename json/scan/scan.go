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

package scan

import "io"

type Scanner interface {
	Next() (Kind, []byte, error)
}

type scanner struct {
	buf    []byte
	offset int
}

func NewScanner(buf []byte) Scanner {
	return &scanner{
		buf:    buf,
		offset: 0,
	}
}

func (s *scanner) Next() (Kind, []byte, error) {
	if err := s.skipSpace(); err != nil {
		return unknownKind, nil, err
	}
	if s.offset == len(s.buf) {
		return unknownKind, nil, io.EOF
	}
	c, err := s.look()
	if err != nil {
		return unknownKind, nil, err
	}
	kind := getKind(c)
	start := s.offset
	switch kind {
	case objectOpenKind, objectCloseKind, arrayOpenKind, arrayCloseKind, colonKind, commaKind:
		if err := s.incOffset(1); err != nil {
			return unknownKind, nil, err
		}
	case stringKind:
		if err := s.scanString(); err != nil {
			return unknownKind, nil, err
		}
	case numberKind:
		if err := s.scanNumber(); err != nil {
			return unknownKind, nil, err
		}
	case trueKind:
		if err := s.scanTrue(); err != nil {
			return unknownKind, nil, err
		}
	case falseKind:
		if err := s.scanFalse(); err != nil {
			return unknownKind, nil, err
		}
	case nullKind:
		if err := s.scanNull(); err != nil {
			return unknownKind, nil, err
		}
	}
	end := s.offset
	token := s.buf[start:end]
	return kind, token, nil
}

func (s *scanner) scanNull() error {
	n, err := Null(s.buf[s.offset:])
	if err != nil {
		return err
	}
	return s.incOffset(n)
}

func (s *scanner) scanFalse() error {
	n, err := False(s.buf[s.offset:])
	if err != nil {
		return err
	}
	return s.incOffset(n)
}

func (s *scanner) scanTrue() error {
	n, err := True(s.buf[s.offset:])
	if err != nil {
		return err
	}
	return s.incOffset(n)
}

func (s *scanner) scanNumber() error {
	n, err := Number(s.buf[s.offset:])
	if err != nil {
		return err
	}
	return s.incOffset(n)
}

func (s *scanner) scanString() error {
	n, err := String(s.buf[s.offset:])
	if err != nil {
		return err
	}
	return s.incOffset(n)
}

func (s *scanner) skipSpace() error {
	if s.offset >= len(s.buf) {
		return nil
	}
	n := Space(s.buf[s.offset:])
	if err := s.incOffset(n); err != nil {
		return err
	}
	return nil
}

func (s *scanner) look() (byte, error) {
	if s.offset < len(s.buf) {
		return s.buf[s.offset], nil
	}
	return 0, io.ErrShortBuffer
}

func (s *scanner) incOffset(o int) error {
	s.offset = s.offset + o
	if s.offset > len(s.buf) {
		return io.ErrShortBuffer
	}
	return nil
}
