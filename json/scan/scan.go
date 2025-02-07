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

import (
	"io"
)

// Next returns the Kind and the start and end offset of the token or an error.
// The start is not always zero, since spaces can be skipped.
func Next(buf []byte, offset int) (Kind, int, int, error) {
	var err error
	offset, err = skipSpace(buf, offset)
	if err != nil {
		return UnknownKind, offset, offset, err
	}
	start := offset
	if offset == len(buf) {
		return UnknownKind, start, offset, io.EOF
	}
	c, err := look(buf, offset)
	if err != nil {
		return UnknownKind, start, offset, err
	}
	kind := getKind(c)
	switch kind {
	case ObjectOpenKind, ObjectCloseKind, ArrayOpenKind, ArrayCloseKind, ColonKind, CommaKind:
		offset, err = incOffset(buf, offset, 1)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	case StringKind:
		offset, err = scanString(buf, offset)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	case NumberKind:
		offset, err = scanNumber(buf, offset)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	case TrueKind:
		offset, err = scanTrue(buf, offset)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	case FalseKind:
		offset, err = scanFalse(buf, offset)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	case NullKind:
		offset, err = scanNull(buf, offset)
		if err != nil {
			return UnknownKind, start, offset, err
		}
	}
	return kind, start, offset, nil
}

var kindMap = map[byte]Kind{
	'{': ObjectOpenKind,
	'}': ObjectCloseKind,
	':': ColonKind,
	'[': ArrayOpenKind,
	']': ArrayCloseKind,
	',': CommaKind,
	'"': StringKind,
	't': TrueKind,
	'f': FalseKind,
	'n': NullKind,
	'-': NumberKind,
	'0': NumberKind,
	'1': NumberKind,
	'2': NumberKind,
	'3': NumberKind,
	'4': NumberKind,
	'5': NumberKind,
	'6': NumberKind,
	'7': NumberKind,
	'8': NumberKind,
	'9': NumberKind,
}

func getKind(b byte) Kind {
	k, ok := kindMap[b]
	if ok {
		return k
	}
	return UnknownKind
}

func scanNull(buf []byte, offset int) (int, error) {
	n, err := Null(buf[offset:])
	if err != nil {
		return 0, err
	}
	return incOffset(buf, offset, n)
}

func scanFalse(buf []byte, offset int) (int, error) {
	n, err := False(buf[offset:])
	if err != nil {
		return 0, err
	}
	return incOffset(buf, offset, n)
}

func scanTrue(buf []byte, offset int) (int, error) {
	n, err := True(buf[offset:])
	if err != nil {
		return 0, err
	}
	return incOffset(buf, offset, n)
}

func scanNumber(buf []byte, offset int) (int, error) {
	n, err := Number(buf[offset:])
	if err != nil {
		return 0, err
	}
	return incOffset(buf, offset, n)
}

func scanString(buf []byte, offset int) (int, error) {
	n, err := String(buf[offset:])
	if err != nil {
		return 0, err
	}
	return incOffset(buf, offset, n)
}

func skipSpace(buf []byte, offset int) (int, error) {
	if offset >= len(buf) {
		return offset, nil
	}
	n := Space(buf[offset:])
	return incOffset(buf, offset, n)
}

func look(buf []byte, offset int) (byte, error) {
	if offset < len(buf) {
		return buf[offset], nil
	}
	return 0, io.ErrShortBuffer
}

func incOffset(buf []byte, offset int, inc int) (int, error) {
	offset = offset + inc
	if offset > len(buf) {
		return 0, io.ErrShortBuffer
	}
	return offset, nil
}
