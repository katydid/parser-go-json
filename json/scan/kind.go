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

type Kind byte

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

const UnknownKind = Kind(0)

func (k Kind) isUnknown() bool {
	return k == UnknownKind
}

const ObjectOpenKind = Kind('{')

func (k Kind) isObjectOpen() bool {
	return k == ObjectOpenKind
}

const ObjectCloseKind = Kind('}')

func (k Kind) isObjectClose() bool {
	return k == ObjectCloseKind
}

const ColonKind = Kind(':')

func (k Kind) isColon() bool {
	return k == ColonKind
}

const ArrayOpenKind = Kind('[')

func (k Kind) isArrayOpen() bool {
	return k == ArrayOpenKind
}

const ArrayCloseKind = Kind(']')

func (k Kind) isArrayClose() bool {
	return k == ArrayCloseKind
}

const CommaKind = Kind(',')

func (k Kind) isComma() bool {
	return k == CommaKind
}

const StringKind = Kind('"')

func (k Kind) isString() bool {
	return k == StringKind
}

const NumberKind = Kind('0')

func (k Kind) isNumber() bool {
	return k == NumberKind
}

const TrueKind = Kind('t')

func (k Kind) isTrue() bool {
	return k == TrueKind
}

const FalseKind = Kind('f')

func (k Kind) isFalse() bool {
	return k == FalseKind
}

const NullKind = Kind('n')

func (k Kind) isNull() bool {
	return k == NullKind
}

func (k Kind) String() string {
	switch k {
	case UnknownKind:
		return "unknown"
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
	case ArrayCloseKind:
		return "arrayClose"
	case ObjectOpenKind:
		return "objectOpen"
	case ObjectCloseKind:
		return "objectClose"
	}
	return "other"
}
