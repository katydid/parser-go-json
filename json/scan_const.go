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

import "bytes"

func scanConst(buf []byte, valBytes []byte, err error) (int, error) {
	if len(buf) < len(valBytes) {
		return 0, err
	}
	if !bytes.Equal(buf[0:len(valBytes)], valBytes) {
		return 0, err
	}
	return len(valBytes), nil
}

var trueBytes = []byte{'t', 'r', 'u', 'e'}

func scanTrue(buf []byte) (int, error) {
	return scanConst(buf, trueBytes, errExpectedTrue)
}

func (s *jsonParser) scanTrue() error {
	n, err := scanTrue(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

var falseBytes = []byte{'f', 'a', 'l', 's', 'e'}

func scanFalse(buf []byte) (int, error) {
	return scanConst(buf, falseBytes, errExpectedFalse)
}

func (s *jsonParser) scanFalse() error {
	n, err := scanFalse(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}

var nullBytes = []byte{'n', 'u', 'l', 'l'}

func scanNull(buf []byte) (int, error) {
	return scanConst(buf, nullBytes, errExpectedNull)
}

func (s *jsonParser) scanNull() error {
	n, err := scanNull(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}
