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
	'{': objectOpenKind,
	'}': objectCloseKind,
	':': colonKind,
	'[': arrayOpenKind,
	']': arrayCloseKind,
	',': commaKind,
	'"': stringKind,
	't': trueKind,
	'f': falseKind,
	'n': nullKind,
	'-': numberKind,
	'0': numberKind,
	'1': numberKind,
	'2': numberKind,
	'3': numberKind,
	'4': numberKind,
	'5': numberKind,
	'6': numberKind,
	'7': numberKind,
	'8': numberKind,
	'9': numberKind,
}

func getKind(b byte) Kind {
	k, ok := kindMap[b]
	if ok {
		return k
	}
	return unknownKind
}

const unknownKind = Kind(0)

func (k Kind) isUnknown() bool {
	return k == unknownKind
}

const objectOpenKind = Kind('{')

func (k Kind) isObjectOpen() bool {
	return k == objectOpenKind
}

const objectCloseKind = Kind('}')

func (k Kind) isObjectClose() bool {
	return k == objectCloseKind
}

const colonKind = Kind(':')

func (k Kind) isColon() bool {
	return k == colonKind
}

const arrayOpenKind = Kind('[')

func (k Kind) isArrayOpen() bool {
	return k == arrayOpenKind
}

const arrayCloseKind = Kind(']')

func (k Kind) isArrayClose() bool {
	return k == arrayOpenKind
}

const commaKind = Kind(',')

func (k Kind) isComma() bool {
	return k == colonKind
}

const stringKind = Kind('"')

func (k Kind) isString() bool {
	return k == stringKind
}

const numberKind = Kind('0')

func (k Kind) isNumber() bool {
	return k == numberKind
}

const trueKind = Kind('t')

func (k Kind) isTrue() bool {
	return k == trueKind
}

const falseKind = Kind('f')

func (k Kind) isFalse() bool {
	return k == falseKind
}

const nullKind = Kind('n')

func (k Kind) isNull() bool {
	return k == nullKind
}

func (k Kind) String() string {
	switch k {
	case unknownKind:
		return "unknown"
	case falseKind:
		return "false"
	case trueKind:
		return "true"
	case numberKind:
		return "number"
	case stringKind:
		return "string"
	case arrayOpenKind:
		return "arrayOpen"
	case arrayCloseKind:
		return "arrayClose"
	case objectOpenKind:
		return "objectOpen"
	case objectCloseKind:
		return "objectClose"
	}
	return "other"
}
