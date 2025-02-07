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

import "bytes"

func scanConst(buf []byte, valBytes []byte, err error) (int, error) {
	if len(buf) < len(valBytes) {
		return 0, err
	}
	if !bytes.Equal(buf[0:len(valBytes)], valBytes) {
		return 0, err
	}
	return len(valBytes), nil
}

var trueBytes = []byte{'t', 'r', 'u', 'e'}

// True returns 4, nil if the prefix of the bytes matches true, otherwise error.
func True(buf []byte) (int, error) {
	return scanConst(buf, trueBytes, errExpectedTrue)
}

var falseBytes = []byte{'f', 'a', 'l', 's', 'e'}

// False returns 5, nil if the prefix of the bytes matches false, otherwise error.
func False(buf []byte) (int, error) {
	return scanConst(buf, falseBytes, errExpectedFalse)
}

var nullBytes = []byte{'n', 'u', 'l', 'l'}

// Null returns 4, nil if the prefix of the bytes matches null, otherwise error.
func Null(buf []byte) (int, error) {
	return scanConst(buf, nullBytes, errExpectedNull)
}
