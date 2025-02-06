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

type Scanner interface {
	// Next returns the Kind and the slice of the buffer containing the token or an error.
	Next() (Kind, []byte, error)
	// Init restarts the scanner with a new byte buffer, without allocating a new scanner.
	Init([]byte)
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

// Next returns the Kind and the slice of the buffer containing the token or an error.
func (s *scanner) Next() (Kind, []byte, error) {
	kind, start, end, err := Next(s.buf, s.offset)
	if err != nil {
		return kind, nil, err
	}
	buf := s.buf[start:end]
	s.offset = end
	return kind, buf, nil
}
