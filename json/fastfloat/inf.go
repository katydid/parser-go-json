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

package fastfloat

import "math"

var inf = math.Inf(1)

func isInf(bs []byte) bool {
	if len(bs) == 3 {
		if bs[0] != 'i' && bs[0] != 'I' {
			return false
		}
		if bs[1] != 'n' && bs[1] != 'N' {
			return false
		}
		if bs[2] != 'f' && bs[2] != 'F' {
			return false
		}
		return true
	} else if len(bs) == 8 {
		if bs[0] != 'i' && bs[0] != 'I' {
			return false
		}
		if bs[1] != 'n' && bs[1] != 'N' {
			return false
		}
		if bs[2] != 'f' && bs[2] != 'F' {
			return false
		}
		if bs[3] != 'i' && bs[3] != 'I' {
			return false
		}
		if bs[4] != 'n' && bs[4] != 'N' {
			return false
		}
		if bs[5] != 'i' && bs[5] != 'I' {
			return false
		}
		if bs[6] != 't' && bs[6] != 'T' {
			return false
		}
		if bs[7] != 'y' && bs[7] != 'Y' {
			return false
		}
		return true
	}
	return false
}
