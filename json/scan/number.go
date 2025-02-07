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

// Number returns the offset after the prefix of a valid number.
// The number BNF:
// number := integer fraction exponent
// integer := digit | onenine digits | '-' digit | '-' onenine digits
// digits := digit | digit digits
// digit := '0' | onenine
// onenine := '1' . '9'
// fraction := "" | '.' digits
// exponent := "" | 'E' sign digits | 'e' sign digits
// sign := "" | '+' | '-'
func Number(buf []byte) (int, error) {
	state := 's' // start
	offset := 0
	for _, c := range buf[offset:] {
		offset += 1
		switch state {
		case 's': // start
			switch c {
			case '0':
				state = '0'
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '1'
			case '-':
				state = '-'
			default:
				return 0, errScanNumber
			}
		case '-': // negative number started
			switch c {
			case '0':
				state = '0'
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '1'
			default:
				return 0, errScanNumber
			}
		case '0': // integer complete
			switch c {
			case '.':
				state = '.'
			case 'e', 'E':
				state = 'e'
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				return 0, errScanNumber
			default:
				return offset - 1, nil // zero
			}
		case '1':
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '.':
				state = '.'
			case 'e', 'E':
				state = 'e'
			default:
				return offset - 1, nil // integer
			}
		case '.': // fraction started
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '/'
			case 'e', 'E':
				state = 'e'
			default:
				return 0, errScanNumber // number does not end in a .
			}
		case '/': // fraction ongoing
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case 'e', 'E':
				state = 'e'
			case '.':
				return 0, errScanNumber // fraction does not contain a .
			default:
				return offset - 1, nil // integer with fraction
			}
		case 'e': // exponent started
			switch c {
			case '-', '+':
				state = 'f'
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = 'g'
			default:
				return 0, errScanNumber // number does not end in a e or E
			}
		case 'f': // exponent with sign started
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = 'g'
			default:
				return 0, errScanNumber // number does not end in '+' or '-'
			}
		case 'g': // exponent ongoing
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				return offset - 1, nil
			}
		default:
			panic("unreachable")
		}
	}
	if state == 's' || state == '-' || state == '.' || state == 'e' || state == 'f' {
		return 0, errScanNumber
	}
	return offset, nil
}
