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

import (
	"testing"

	"github.com/katydid/parser-go-json/json/internal/testrun"
	"github.com/katydid/parser-go-json/json/rand"
)

func TestString(t *testing.T) {
	valid := map[string]int{
		`""`:    2,
		`"a"`:   3,
		`"a" `:  3,
		`"abc"`: 5,
		`"\b"`:  4,
		`"\b" `: 4,
		`"\""`:  4,
	}
	invalid := []string{
		`"`,
	}
	for input, want := range valid {
		t.Run("Valid("+input+")", func(t *testing.T) {
			got, err := String([]byte(input))
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
			_, err := String([]byte(input))
			if err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestRandString(t *testing.T) {
	r := rand.NewRand()
	for i := 0; i < 100; i++ {
		s := rand.String(r)
		name := testrun.Name([]byte(s))
		t.Run(name, func(t *testing.T) {
			buf := []byte(s)
			got, err := String(buf)
			if err != nil {
				t.Fatal(err)
			}
			if got != len(buf) {
				t.Fatalf("expected offset = %d, but got %d", len(buf), got)
			}
		})
	}
}
