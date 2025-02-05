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

// kind of the token that is tokenized.
// This is represented by one for following bytes: {:}[,]"0tfn>-+.
type kind byte

const UnknownKind = kind(0)

func (k kind) IsUnknown() bool {
	return k == UnknownKind
}

const ObjectOpenKind = kind('{')

func (k kind) IsObjectOpen() bool {
	return k == ObjectOpenKind
}

const ObjectCloseKind = kind('}')

func (k kind) IsObjectClose() bool {
	return k == ObjectCloseKind
}

const ColonKind = kind(':')

func (k kind) IsColon() bool {
	return k == ColonKind
}

const ArrayOpenKind = kind('[')

func (k kind) IsArrayOpen() bool {
	return k == ArrayOpenKind
}

const ArrayCloseKind = kind(']')

func (k kind) IsArrayClose() bool {
	return k == ArrayCloseKind
}

const CommaKind = kind(',')

func (k kind) IsComma() bool {
	return k == CommaKind
}

const StringKind = kind('"')

func (k kind) IsString() bool {
	return k == StringKind
}

const TrueKind = kind('t')

func (k kind) IsTrue() bool {
	return k == TrueKind
}

const FalseKind = kind('f')

func (k kind) IsFalse() bool {
	return k == FalseKind
}

const NullKind = kind('n')

func (k kind) IsNull() bool {
	return k == NullKind
}

// a number can be:
// * unknown, since it is not parsed yet.
// * int, uint and double (any number, for example 123) represented by '0'
// * int and double, but not uint (negative number) represented by '-'
// * uint and double, but not int (a large positive number) represented by '+'
// * double, but not int or uint (a fraction) represented by '.'
// * none, since it is a number too large to fit even in double, represented by '>'

const TooLargeNumberKind = kind('>')

func (k kind) IsTooLargeNumber() bool {
	return k == TooLargeNumberKind
}

const NegativeNumberKind = kind('-')

func (k kind) IsInt() bool {
	return k == NumberKind || k == NegativeNumberKind
}

const LargePositiveNumberKind = kind('+')

func (k kind) IsUint() bool {
	return k == NumberKind || k == LargePositiveNumberKind
}

const FractionNumberKind = kind('.')

func (k kind) IsDouble() bool {
	return k == NumberKind || k == LargePositiveNumberKind || k == NegativeNumberKind || k == FractionNumberKind
}

const NumberKind = kind('0')

func (k kind) IsNumber() bool {
	return k == NumberKind || k == LargePositiveNumberKind || k == NegativeNumberKind || k == FractionNumberKind || k == TooLargeNumberKind
}

func (k kind) String() string {
	switch k {
	case UnknownKind:
		return "unknown"
	case NullKind:
		return "null"
	case FalseKind:
		return "false"
	case TrueKind:
		return "true"
	case NumberKind:
		return "number"
	case StringKind:
		return "string"
	case ArrayOpenKind:
		return "arrayOpen"
	case CommaKind:
		return "comma"
	case ArrayCloseKind:
		return "arrayClose"
	case ObjectOpenKind:
		return "objectOpen"
	case ColonKind:
		return "colon"
	case ObjectCloseKind:
		return "objectClose"
	case TooLargeNumberKind:
		return "tooLargeNumberKind"
	case NegativeNumberKind:
		return "negativeNumberKind"
	case LargePositiveNumberKind:
		return "largePositiveNumberKind"
	case FractionNumberKind:
		return "fractionNumberKind"
	}
	return "other"
}
