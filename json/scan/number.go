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

import "github.com/katydid/parser-go-json/json/internal/fork/strconv"

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
func ParseNumber(buf []byte) (offset int, intres int64, intok bool, floatres float64, floatok bool, decimalok bool) {
	state := 's' // start
	sign := int64(1)
	expsign := 1
	var mantissa uint64
	var intdigits int
	var fracdigits int
	var exppart int
	var expdigits int
	maxMantDigits := 19 // 10^19 fits in uint64
	done := false
	for _, c := range buf {
		offset += 1
		switch state {
		case 's': // start
			switch c {
			case '0':
				state = '0'
				intdigits++
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '1'
				if intdigits < maxMantDigits {
					mantissa = (mantissa * 10) + uint64(c-'0')
				}
				intdigits++
			case '-':
				state = '-'
				sign = -1
			default:
				return // ERROR: number starts with a digit or '-'
			}
		case '-': // negative number started
			switch c {
			case '0':
				state = '0'
				intdigits++
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '1'
				if intdigits < maxMantDigits {
					mantissa = (mantissa * 10) + uint64(c-'0')
				}
				intdigits++
			default:
				return // ERROR: negative numbers starts with a digit
			}
		case '0': // integer complete
			switch c {
			case '.':
				state = '.'
			case 'e', 'E':
				state = 'e'
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				return // ERROR: the integer was complete, for example 0, the number does not continue
			default:
				done = true // SUCCESS zero
			}
		case '1':
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if intdigits < maxMantDigits {
					mantissa = (mantissa * 10) + uint64(c-'0')
				}
				intdigits++
			case '.':
				state = '.'
			case 'e', 'E':
				state = 'e'
			default:
				done = true // SUCCESS integer
			}
		case '.': // fraction started
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = '/'
				if fracdigits+intdigits < maxMantDigits {
					mantissa = (mantissa * 10) + uint64(c-'0')
				}
				fracdigits++
			case 'e', 'E':
				state = 'e'
			default:
				return // ERROR: number does not end in a .
			}
		case '/': // fraction ongoing
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if fracdigits+intdigits < maxMantDigits {
					mantissa = (mantissa * 10) + uint64(c-'0')
				}
				fracdigits++
			case 'e', 'E':
				state = 'e'
			case '.':
				return // ERROR: fraction does not contain a .
			default:
				done = true // SUCCESS
			}
		case 'e': // exponent started
			switch c {
			case '-':
				state = 'f'
				expsign = -1
			case '+':
				state = 'f'
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = 'g'
				if exppart < 10000 {
					exppart = (exppart * 10) + int(c-'0')
				}
				expdigits++
			default:
				return // ERROR: number does not end in a e or E
			}
		case 'f': // exponent with sign started
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = 'g'
				if exppart < 10000 {
					exppart = (exppart * 10) + int(c-'0')
				}
				expdigits++
			default:
				return // ERROR: number does not end in '+' or '-'
			}
		case 'g': // exponent ongoing
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if exppart < 10000 {
					exppart = (exppart * 10) + int(c-'0')
				}
				expdigits++
			default:
				done = true // SUCCESS
			}
		default:
			panic("unreachable")
		}
		if done {
			offset = offset - 1
			break
		}
	}
	if state == 's' || state == '-' || state == '.' || state == 'e' || state == 'f' {
		return // ERROR: these are not accepting states
	}

	neg := sign == -1
	// It is an int
	if expdigits == 0 && fracdigits == 0 {
		if intdigits > maxMantDigits {
			// It uses more digits than MaxInt64 and MinInt64, so it is decimal
			decimalok = true
			return
		}
		cutOffInt64 := uint64(1 << uint(64-1))
		if (!neg && mantissa < cutOffInt64) || (neg && mantissa <= cutOffInt64) {
			intres = sign * int64(mantissa)
			intok = true
			return
		} else {
			// It is larger than an int, so it must be decimal
			decimalok = true
			return
		}
	}

	// Check if float needs to be truncated
	trunc := false
	ndMant := intdigits + fracdigits
	if ndMant >= maxMantDigits {
		trunc = true
		ndMant = maxMantDigits
	}

	// calculate exp
	exp := 0
	dp := intdigits
	if expdigits > 0 {
		dp += exppart * expsign
	}
	if mantissa != 0 {
		exp = dp - ndMant
	}

	f, ok := strconv.TryParseFloat(buf[:offset], mantissa, exp, neg, trunc)
	if ok {
		floatres = f
		floatok = true
		return
	}

	// It is larger than an int, so it must be decimal
	decimalok = true
	return
}
