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

func Next(s Scanner) (Kind, []byte, error) {
	k, _, err := s.NextStart()
	if err != nil {
		return k, nil, err
	}
	bs, err := s.ScanToEnd()
	return k, bs, err
}

type Scanner interface {
	// Init restarts the scanner with a new byte buffer, without allocating a new scanner.
	Init([]byte)

	// NextStart skips to the start of the next token and returns a byte slice that start at that token.
	// Following that the client needs to either:
	//   * call ScanToEnd to automatically Scan to the end of the that token and get the slice that contains only that token.
	//   * call Skip to pass back the offset of the end of that token, since the client was able to scan it themselves.
	NextStart() (Kind, []byte, error)
	ScanToEnd() ([]byte, error)
	Skip(offset int) error
}

type scanner struct {
	buf    []byte
	offset int
}

// NewScanner returns a Scanner which keeps track of the buffer and the offset.
func NewScanner(buf []byte) Scanner {
	return &scanner{
		buf:    buf,
		offset: 0,
	}
}

// Init restarts the scanner with a new byte buffer, without allocating a new scanner.
func (s *scanner) Init(buf []byte) {
	s.buf = buf
	s.offset = 0
}

func (s *scanner) NextStart() (Kind, []byte, error) {
	kind, start, err := NextStart(s.buf, s.offset)
	if err != nil {
		return kind, nil, err
	}
	buf := s.buf[start:]
	s.offset = start
	return kind, buf, nil
}

func (s *scanner) Skip(offset int) error {
	s.offset += offset
	if s.offset > len(s.buf) {
		return io.ErrShortBuffer
	}
	return nil
}

func (s *scanner) ScanToEnd() ([]byte, error) {
	start := s.offset
	end, err := NextEnd(s.buf, s.offset)
	if err != nil {
		return nil, err
	}
	s.offset = end
	return s.buf[start:s.offset], nil
}
