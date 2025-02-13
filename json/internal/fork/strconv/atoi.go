// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv

// lower(c) is a lower-case letter if and only if
// c is either that lower-case letter or the equivalent upper-case letter.
// Instead of writing c == 'x' || c == 'X' one can write lower(c) == 'x'.
// Note that lower of non-letters can produce other non-letters.
func lower(c byte) byte {
	return c | ('x' - 'X')
}

const maxUint64 = 1<<64 - 1

// ParseUint is like [ParseInt] but for unsigned numbers.
//
// A sign prefix is not permitted.
func ParseUint(s []byte) (uint64, error) {
	return parseUint(s, 64)
}

// Cutoff is the smallest number such that cutoff*base > maxUint64.
// Use compile-time constants for common cases.
const cutOffUint64 uint64 = maxUint64/10 + 1

func parseUint(s []byte, bitSize int) (uint64, error) {
	base := 10

	if len(s) == 0 {
		return 0, errSyntaxParseUint
	}

	if bitSize == 0 {
		bitSize = 64
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		default:
			return 0, errSyntaxParseUint
		}

		if n >= cutOffUint64 {
			// n*base overflows
			return maxVal, errRangeParseUint
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal, errRangeParseUint
		}
		n = n1
	}

	return n, nil
}

// ParseInt interprets a string s in the given base (0, 2 to 36) and
// bit size (0 to 64) and returns the corresponding value i.
//
// The string may begin with a leading sign: "+" or "-".
//
// If the base argument is 0, the true base is implied by the string's
// prefix following the sign (if present): 2 for "0b", 8 for "0" or "0o",
// 16 for "0x", and 10 otherwise. Also, for argument base 0 only,
// underscore characters are permitted as defined by the Go syntax for
// [integer literals].
//
// The bitSize argument specifies the integer type
// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
// correspond to int, int8, int16, int32, and int64.
// If bitSize is below 0 or above 64, an error is returned.
//
// The errors that ParseInt returns have concrete type [*NumError]
// and include err.Num = s. If s is empty or contains invalid
// digits, err.Err = [ErrSyntax] and the returned value is 0;
// if the value corresponding to s cannot be represented by a
// signed integer of the given size, err.Err = [ErrRange] and the
// returned value is the maximum magnitude integer of the
// appropriate bitSize and sign.
//
// [integer literals]: https://go.dev/ref/spec#Integer_literals
func ParseInt(s []byte) (i int64, err error) {
	return parseInt(s)
}

const cutOffInt64 uint64 = uint64(1 << uint(64-1))

func parseInt(s []byte) (i int64, err error) {

	if len(s) == 0 {
		return 0, errSyntaxParseInt
	}

	// Pick off leading sign.
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	var un uint64
	un, err = parseUint(s, 64)
	if err != nil && err.(*NumError).Err != ErrRange {
		return 0, errSyntaxParseInt
	}

	cutoff := cutOffInt64
	if !neg && un >= cutoff {
		return int64(cutoff - 1), errRangeParseInt
	}
	if neg && un > cutoff {
		return -int64(cutoff), errRangeParseInt
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Atoi(s []byte) (int, error) {

	sLen := len(s)
	if 0 < sLen && sLen < 19 {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0, errSyntaxAtoi
			}
		}

		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0, errSyntaxAtoi
			}
			n = n*10 + int(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid, big, or underscored integers.
	i64, err := parseInt(s)
	if nerr, ok := err.(*NumError); ok {
		if nerr.Err == ErrRange {
			return int(i64), errRangeAtoi
		} else {
			return int(i64), errSyntaxAtoi
		}
	}
	return int(i64), err
}
