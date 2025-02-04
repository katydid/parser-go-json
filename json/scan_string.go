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

func scanString(buf []byte) (int, error) {
	escaped := false
	udigits := -1
	if len(buf) == 0 || buf[0] != '"' {
		return 0, errScanString
	}
	for i, c := range buf[1:] {
		if escaped {
			switch c {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				escaped = false
				continue
			case 'u':
				udigits = 0
				escaped = false
				continue
			}
			return 0, errScanString
		}
		if udigits >= 0 {
			if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
				udigits++
			} else {
				return 0, errScanString
			}
			if udigits == 4 {
				udigits = -1
			}
			continue
		}
		if c == '"' {
			return i + 2, nil
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c < 0x20 {
			return 0, errScanString
		}
	}
	return 0, errScanString
}

func (s *jsonParser) scanString() error {
	n, err := scanString(s.buf[s.offset:])
	if err != nil {
		return err
	}
	if err := s.incOffset(n); err != nil {
		return err
	}
	return s.skipSpace()
}
