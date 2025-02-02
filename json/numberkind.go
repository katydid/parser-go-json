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

// a number can be:
// * unknown, since it is not parsed yet.
// * int, uint and double (for example 123)
// * int and double, but not uint (negative number)
// * uint and double, but not int (a large positive number)
// * double, but not int or uint (a fraction)
// * none, since it is a number too large to fit even in double.

type kindOfNumber byte

const unknownNumberKind = kindOfNumber(0)

const anyOfNumber = kindOfNumber('0')

const noneOfNumber = kindOfNumber('>')

const intOfNumber = kindOfNumber('-')

func (k kindOfNumber) isInt() bool {
	return k == anyOfNumber || k == intOfNumber
}

const uintOfNumber = kindOfNumber('+')

func (k kindOfNumber) isUint() bool {
	return k == anyOfNumber || k == uintOfNumber
}

const doubleOfNumber = kindOfNumber('.')

func (k kindOfNumber) isDouble() bool {
	return k == anyOfNumber || k == uintOfNumber || k == intOfNumber || k == doubleOfNumber
}
