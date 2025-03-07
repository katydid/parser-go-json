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

// Kind of the token that is scanned.
// This is represented by one for following bytes: {:}[,]"0tfn
type Kind byte

const UnknownKind = Kind(0)

func (k Kind) IsUnknown() bool {
	return k == UnknownKind
}

const ObjectOpenKind = Kind('{')

func (k Kind) IsObjectOpen() bool {
	return k == ObjectOpenKind
}

const ObjectCloseKind = Kind('}')

func (k Kind) IsObjectClose() bool {
	return k == ObjectCloseKind
}

const ColonKind = Kind(':')

func (k Kind) IsColon() bool {
	return k == ColonKind
}

const ArrayOpenKind = Kind('[')

func (k Kind) IsArrayOpen() bool {
	return k == ArrayOpenKind
}

const ArrayCloseKind = Kind(']')

func (k Kind) IsArrayClose() bool {
	return k == ArrayCloseKind
}

const CommaKind = Kind(',')

func (k Kind) IsComma() bool {
	return k == CommaKind
}

const StringKind = Kind('"')

func (k Kind) IsString() bool {
	return k == StringKind
}

const NumberKind = Kind('0')

func (k Kind) IsNumber() bool {
	return k == NumberKind
}

const TrueKind = Kind('t')

func (k Kind) IsTrue() bool {
	return k == TrueKind
}

const FalseKind = Kind('f')

func (k Kind) IsFalse() bool {
	return k == FalseKind
}

const NullKind = Kind('n')

func (k Kind) IsNull() bool {
	return k == NullKind
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
		return "colonKind"
	case ObjectCloseKind:
		return "objectClose"
	}
	return "other"
}
