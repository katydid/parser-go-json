// Copyright (c) 2018 Aliaksandr Valialkin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package fastfloat

import (
	"math"
	"strconv"
)

// ParseUint64 parses uint64 from s.
//
// It is equivalent to strconv.ParseUint(s, 10, 64), but is faster.
//
// See also ParseUint64BestEffort.
func ParseUint64(s []byte) (uint64, error) {
	if len(s) == 0 {
		return 0, errCannotParseNumberFromEmpty
	}
	i := uint(0)
	d := uint64(0)
	j := i
	for i < uint(len(s)) {
		if s[i] >= '0' && s[i] <= '9' {
			d = d*10 + uint64(s[i]-'0')
			i++
			if i > 18 {
				// The integer part may be out of range for uint64.
				// Fall back to slow parsing.
				dd, err := strconv.ParseUint(string(s), 10, 64)
				if err != nil {
					return 0, err
				}
				return dd, nil
			}
			continue
		}
		break
	}
	if i <= j {
		return 0, errCannotParseNumber
	}
	if i < uint(len(s)) {
		// Unparsed tail left.
		return 0, errCannotParseNumberUnparsedTail
	}
	return d, nil
}

// ParseInt64 parses int64 number s.
//
// It is equivalent to strconv.ParseInt(s, 10, 64), but is faster.
//
// See also ParseInt64BestEffort.
func ParseInt64(s []byte) (int64, error) {
	if len(s) == 0 {
		return 0, errCannotParseNumberFromEmpty
	}
	i := uint(0)
	minus := s[0] == '-'
	if minus {
		i++
		if i >= uint(len(s)) {
			return 0, errCannotParseNumber
		}
	}

	d := int64(0)
	j := i
	for i < uint(len(s)) {
		if s[i] >= '0' && s[i] <= '9' {
			d = d*10 + int64(s[i]-'0')
			i++
			if i > 18 {
				// The integer part may be out of range for int64.
				// Fall back to slow parsing.
				dd, err := strconv.ParseInt(string(s), 10, 64)
				if err != nil {
					return 0, err
				}
				return dd, nil
			}
			continue
		}
		break
	}
	if i <= j {
		return 0, errCannotParseNumber
	}
	if i < uint(len(s)) {
		// Unparsed tail left.
		return 0, errCannotParseNumberUnparsedTail
	}
	if minus {
		d = -d
	}
	return d, nil
}

// Exact powers of 10.
//
// This works faster than math.Pow10, since it avoids additional multiplication.
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16,
}

// Parse parses floating-point number s.
//
// It is equivalent to strconv.ParseFloat(s, 64), but is faster.
//
// See also ParseBestEffort.
func ParseFloat(s []byte) (float64, error) {
	if len(s) == 0 {
		return 0, errCannotParseNumberFromEmpty
	}
	i := uint(0)
	minus := s[0] == '-'
	if minus {
		i++
		if i >= uint(len(s)) {
			return 0, errCannotParseNumber
		}
	}

	// the integer part might be elided to remain compliant
	// with https://go.dev/ref/spec#Floating-point_literals
	if s[i] == '.' && (i+1 >= uint(len(s)) || s[i+1] < '0' || s[i+1] > '9') {
		return 0, errCannotParseNumberMissingIntInFrac
	}

	d := uint64(0)
	j := i
	for i < uint(len(s)) {
		if s[i] >= '0' && s[i] <= '9' {
			d = d*10 + uint64(s[i]-'0')
			i++
			if i > 18 {
				// The integer part may be out of range for uint64.
				// Fall back to slow parsing.
				f, err := strconv.ParseFloat(string(s), 64)
				if err != nil && !math.IsInf(f, 0) {
					return 0, err
				}
				return f, nil
			}
			continue
		}
		break
	}
	if i <= j && s[i] != '.' {
		ss := s[i:]
		if len(ss) > 0 && ss[0] == '+' {
			ss = ss[1:]
		}
		if isInf(ss) {
			if minus {
				return -inf, nil
			}
			return inf, nil
		}
		if isNan(ss) {
			return nan, nil
		}
		return 0, errCannotParseNumberUnparsedTail
	}
	f := float64(d)
	if i >= uint(len(s)) {
		// Fast path - just integer.
		if minus {
			f = -f
		}
		return f, nil
	}

	if s[i] == '.' {
		// Parse fractional part.
		i++
		if i >= uint(len(s)) {
			// the fractional part might be elided to remain compliant
			// with https://go.dev/ref/spec#Floating-point_literals
			return f, nil
		}
		k := i
		for i < uint(len(s)) {
			if s[i] >= '0' && s[i] <= '9' {
				d = d*10 + uint64(s[i]-'0')
				i++
				if i-j >= uint(len(float64pow10)) {
					// The mantissa is out of range. Fall back to standard parsing.
					f, err := strconv.ParseFloat(string(s), 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0, errCannotParseNumberMantissa
					}
					return f, nil
				}
				continue
			}
			break
		}
		if i < k {
			return 0, errCannotParseNumberFindMantissa
		}
		// Convert the entire mantissa to a float at once to avoid rounding errors.
		f = float64(d) / float64pow10[i-k]
		if i >= uint(len(s)) {
			// Fast path - parsed fractional number.
			if minus {
				f = -f
			}
			return f, nil
		}
	}
	if s[i] == 'e' || s[i] == 'E' {
		// Parse exponent part.
		i++
		if i >= uint(len(s)) {
			return 0, errCannotParseNumberExponent
		}
		expMinus := false
		if s[i] == '+' || s[i] == '-' {
			expMinus = s[i] == '-'
			i++
			if i >= uint(len(s)) {
				return 0, errCannotParseNumberExponent
			}
		}
		exp := int16(0)
		j := i
		for i < uint(len(s)) {
			if s[i] >= '0' && s[i] <= '9' {
				exp = exp*10 + int16(s[i]-'0')
				i++
				if exp > 300 {
					// The exponent may be too big for float64.
					// Fall back to standard parsing.
					f, err := strconv.ParseFloat(string(s), 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0, errCannotParseNumberExponent
					}
					return f, nil
				}
				continue
			}
			break
		}
		if i <= j {
			return 0, errCannotParseNumberExponent
		}
		if expMinus {
			exp = -exp
		}
		f *= math.Pow10(int(exp))
		if i >= uint(len(s)) {
			if minus {
				f = -f
			}
			return f, nil
		}
	}
	return 0, errCannotParseNumber
}
