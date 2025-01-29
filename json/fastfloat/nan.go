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

var nan = math.NaN()

func isNan(bs []byte) bool {
	if len(bs) == 3 {
		if bs[0] != 'n' && bs[0] != 'N' {
			return false
		}
		if bs[1] != 'a' && bs[1] != 'A' {
			return false
		}
		if bs[2] != 'n' && bs[2] != 'N' {
			return false
		}
		return true
	}
	return false
}
