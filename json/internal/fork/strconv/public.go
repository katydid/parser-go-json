// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strconv

import "math"

func TryParseFloat(s []byte, mantissa uint64, exp int, neg bool, trunc bool) (float64, bool) {
	if !trunc {
		if f, ok := tryExactConvertToFloat(mantissa, exp, neg); ok {
			return f, ok
		}
	}
	if f, ok := tryEiselLemireConvertToFloat(mantissa, exp, neg, trunc); ok {
		return f, ok
	}
	return trySlowApproxConvertToFloat(s)
}

func tryExactConvertToFloat(mantissa uint64, exp int, neg bool) (float64, bool) {
	return atof64exact(mantissa, exp, neg)
}

func tryEiselLemireConvertToFloat(mantissa uint64, exp int, neg bool, trunc bool) (float64, bool) {
	f, ok := eiselLemire64(mantissa, exp, neg)
	if ok {
		if !trunc {
			return f, true
		}
		// Even if the mantissa was truncated, we may
		// have found the correct result. Confirm by
		// converting the upper mantissa bound.
		fUp, ok := eiselLemire64(mantissa+1, exp, neg)
		if ok && f == fUp {
			return f, true
		}
	}
	return 0, false
}

func trySlowApproxConvertToFloat(s []byte) (float64, bool) {
	var d decimal
	if !d.set(s) {
		return 0, false
	}
	b, ovf := d.floatBits(&float64info)
	f := math.Float64frombits(b)
	return f, !ovf
}
