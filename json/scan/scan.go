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
	"errors"
	"io"
)

// NextStart returns the Kind and the start offset of the token or an error.
// The start is not always zero, since spaces can be skipped.
func NextStart(buf []byte, offset int) (Kind, int, error) {
	var err error
	offset, err = skipSpace(buf, offset)
	if err != nil {
		return UnknownKind, offset, err
	}
	start := offset
	if offset == len(buf) {
		return UnknownKind, start, io.EOF
	}
	c, err := look(buf, offset)
	if err != nil {
		return UnknownKind, start, err
	}
	kind := getKind(c)
	return kind, start, nil
}

// NextBegin returns the end offset of the token or an error.
func NextEnd(buf []byte, offset int) (int, error) {
	c, err := look(buf, offset)
	if err != nil {
		return offset, err
	}
	kind := getKind(c)
	switch kind {
	case ObjectOpenKind, ObjectCloseKind, ArrayOpenKind, ArrayCloseKind, ColonKind, CommaKind:
		offset, err = incOffset(buf, offset, 1)
		if err != nil {
			return offset, err
		}
	case StringKind:
		offset, err = scanString(buf, offset)
		if err != nil {
			return offset, err
		}
	case NumberKind:
		offset, err = scanNumber(buf, offset)
		if err != nil {
			return offset, err
		}
	case TrueKind:
		offset, err = scanTrue(buf, offset)
		if err != nil {
			return offset, err
		}
	case FalseKind:
		offset, err = scanFalse(buf, offset)
		if err != nil {
			return offset, err
		}
	case NullKind:
		offset, err = scanNull(buf, offset)
		if err != nil {
			return offset, err
		}
	}
	return offset, nil
}

var errUnknownKind = errors.New("unknown kind")

func getKind(b byte) Kind {
	switch b {
	case '{':
		return ObjectOpenKind
	case '}':
		return ObjectCloseKind
	case ':':
		return ColonKind
	case '[':
		return ArrayOpenKind
	case ']':
		return ArrayCloseKind
	case ',':
		return CommaKind
	case '"':
		return StringKind
	case 't':
		return TrueKind
	case 'f':
		return FalseKind
	case 'n':
		return NullKind
	case '-':
		return NumberKind
	case '0':
		return NumberKind
	case '1':
		return NumberKind
	case '2':
		return NumberKind
	case '3':
		return NumberKind
	case '4':
		return NumberKind
	case '5':
		return NumberKind
	case '6':
		return NumberKind
	case '7':
		return NumberKind
	case '8':
		return NumberKind
	case '9':
		return NumberKind
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
