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

// String returns the offset after the quoted string.
// A string starts and ends with a double quote '"'.
// The string BNF:
// string := '"' characters '"'
// characters := "" | character characters
// character := '0020' . '10FFFF' - '"' - '\' | '\' escape
// escape := '"' | '\' | '/' | 'b' | 'f' | 'n' | 'r' | 't' | 'u' hex hex hex hex
// hex := digit | 'A' . 'F' | 'a' . 'f'
func String(buf []byte) (int, error) {
	if len(buf) == 0 || buf[0] != '"' {
		return 0, errScanString
	}
	i := 1
	for i < len(buf) {
		c := buf[i]
		i++
		isplain := plaintable[c]
		if isplain == 0 {
			continue
		}
		if c == '\\' {
			if i >= len(buf) {
				return 0, errScanString
			}
			switch buf[i] {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				i++
			case 'u':
				i += 4
				if i >= len(buf) {
					return 0, errScanString
				}
				c4 := buf[i]
				c3 := buf[i-1]
				c2 := buf[i-2]
				c1 := buf[i-3]
				if !hextable[c1] || !hextable[c2] || !hextable[c3] || !hextable[c4] {
					return 0, errScanString
				}
			default:
				return 0, errScanString
			}
			continue
		}
		if c == '"' {
			return i, nil
		}
		return 0, errScanString
	}
	return 0, errScanString
}

var plaintable = [256]byte{}

func init() {
	for i := range 0x20 {
		plaintable[i] = 1
	}
	plaintable['"'] = 1
	plaintable['\\'] = 1
}

var hextable = [256]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
}
