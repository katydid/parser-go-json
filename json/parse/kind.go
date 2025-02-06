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

package parse

// Kind of the token that is parsed.
// This is represented by one for following bytes: {}[]"0tfn
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

const ArrayOpenKind = Kind('[')

func (k Kind) IsArrayOpen() bool {
	return k == ArrayOpenKind
}

const ArrayCloseKind = Kind(']')

func (k Kind) IsArrayClose() bool {
	return k == ArrayCloseKind
}

const StringKind = Kind('"')

func (k Kind) IsString() bool {
	return k == StringKind
}

const NumberKind = Kind('0')

func (k Kind) IsNumber() bool {
	return k == NumberKind
}

const BoolKind = Kind('?')

func (k Kind) IsBool() bool {
	return k == BoolKind
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
	case BoolKind:
		return "bool"
	case NumberKind:
		return "number"
	case StringKind:
		return "string"
	case ArrayOpenKind:
		return "arrayOpen"
	case ArrayCloseKind:
		return "arrayClose"
	case ObjectOpenKind:
		return "objectOpen"
	case ObjectCloseKind:
		return "objectClose"
	}
	return "other"
}
