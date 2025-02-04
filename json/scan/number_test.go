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

import "testing"

func TestNumber(t *testing.T) {
	valid := map[string]int{
		"0":       1,
		"1":       1,
		"123":     3,
		"1 ":      1,
		"1 .":     1,
		"-1":      2,
		"1.1":     3,
		"1.123":   5,
		"1.1E1":   5,
		"1.1e1":   5,
		"1.1e-1":  6,
		"1.1E+1":  6,
		"1.1e10":  6,
		"1.1e+10": 7,
	}
	invalid := []string{
		"01",
		"1E+",
		"1E",
		"01 ",
		"1E+ ",
		"1E ",
	}
	for input, want := range valid {
		t.Run("Valid("+input+")", func(t *testing.T) {
			got, err := Number([]byte(input))
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("offset want %d, but got %d", want, got)
			}
		})
	}
	for _, input := range invalid {
		t.Run("Invalid("+input+")", func(t *testing.T) {
			_, err := Number([]byte(input))
			if err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}
