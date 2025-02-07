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

func isSpace(c byte) bool {
	return (c == ' ') || (c == '\n') || (c == '\r') || (c == '\t')
}

// Space returns the next character that is not a prefix.
// Spaces are limited to the following characters ' ', '\n', '\r', '\t'.
// If there are no spaces in the prefix, then Space returns 0
func Space(buf []byte) int {
	for i, c := range buf {
		if !isSpace(c) {
			return i
		}
	}
	return len(buf)
}
