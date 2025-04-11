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

package token

// Kind of the token that is parsed.
// This is represented by one for following bytes:
// * '_': Null (Null)
// * 't': True (Bool)
// * 'f': False (Bool)
// * 'x': Bytes (Bytes)
// * '"': String (String)
// * '-': Int64 (Int64)
// * '.': Float64 (Float64)
// * '/': Decimal (String)
// * '9': Nanoseconds (Int64)
// * 'T': DateTime ISO 8601 (String)
// * '#': User defined tags (String)
type Kind byte

const UnknownKind = Kind(0)

func (k Kind) IsUnknown() bool {
	return k == UnknownKind
}

const NullKind = Kind('_')

func (k Kind) IsNull() bool {
	return k == NullKind
}

const FalseKind = Kind('f')

func (k Kind) IsFalse() bool {
	return k == FalseKind
}

const TrueKind = Kind('t')

func (k Kind) IsTrue() bool {
	return k == TrueKind
}

const BytesKind = Kind('x')

func (k Kind) IsBytes() bool {
	return k == BytesKind
}

const StringKind = Kind('"')

func (k Kind) IsString() bool {
	return k == StringKind
}

const Int64Kind = Kind('-')

func (k Kind) IsInt64() bool {
	return k == Int64Kind
}

const Float64Kind = Kind('.')

func (k Kind) IsFloat64() bool {
	return k == Float64Kind
}

const DecimalKind = Kind('/')

func (k Kind) IsDecimal() bool {
	return k == DecimalKind
}

const NanosecondsKind = Kind('9')

func (k Kind) IsNanoseconds() bool {
	return k == NanosecondsKind
}

const DateTimeKind = Kind('T')

func (k Kind) IsDateTimeKind() bool {
	return k == DateTimeKind
}

const TagKind = Kind('#')

func (k Kind) IsTag() bool {
	return k == TagKind
}

func (k Kind) String() string {
	switch k {
	case UnknownKind:
		return "unknown"
	case NullKind:
		return "null"
	case FalseKind:
		return "false"
	case TrueKind:
		return "true"
	case BytesKind:
		return "bytes"
	case StringKind:
		return "string"
	case Int64Kind:
		return "int64"
	case Float64Kind:
		return "float64"
	case DecimalKind:
		return "decimal"
	case NanosecondsKind:
		return "nanoseconds"
	case DateTimeKind:
		return "dateTime"
	case TagKind:
		return "tag"
	}
	panic("unreachable")
}
