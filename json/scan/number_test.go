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
	"fmt"
	"testing"

	"github.com/katydid/parser-go-json/json/internal/fork/strconv"
	"github.com/katydid/parser-go-json/json/rand"
)

func TestNumber(t *testing.T) {
	valid := map[string]int{
		"0":       1,
		"1":       1,
		"123":     3,
		"1 ":      1,
		"1 .":     1,
		"-1":      2,
		"1.1":     3,
		"1.0":     3,
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

func TestRandNumber(t *testing.T) {
	r := rand.NewRand()
	for i := 0; i < 100; i++ {
		s := rand.Number(r)
		t.Run(s, func(t *testing.T) {
			buf := []byte(s)
			got, err := Number(buf)
			if err != nil {
				t.Fatal(err)
			}
			if got != len(buf) {
				t.Fatalf("expected offset = %d, but got %d", len(buf), got)
			}
		})
	}
}

func TestParseRandomNumber(t *testing.T) {
	r := rand.NewRand()
	for i := 0; i < 1000; i++ {
		s := rand.Number(r)
		t.Run(s, func(t *testing.T) {
			buf := []byte(s)
			if err := checkNumber(buf); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func notParseableInteger(token []byte) bool {
	for _, b := range token {
		if b == '.' || b == 'e' || b == 'E' {
			return true
		}
	}
	return false
}

func checkNumber(token []byte) error {
	offset, intval, intok, floatval, floatok, decok := ParseNumber(token)
	if notParseableInteger(token) {
		slowFloatVal, err := strconv.ParseFloat(token)
		if err != nil {
			if !decok {
				return fmt.Errorf("expected decimal, since we could not parse the float and it is not an integer, but %v is not a decimal", string(token))
			}
			if offset != len(token) {
				return fmt.Errorf("expected decimal %s to parse to end %d, but only parsed up to %d", string(token), len(token), offset)
			}
			return nil
		}
		if !floatok {
			return fmt.Errorf("expected float, since we could parse the float and it is not an integer, but %v is not a float", string(token))
		}
		if slowFloatVal != floatval {
			return fmt.Errorf("want %v got %v", slowFloatVal, floatval)
		}
		return nil
	}
	slowIntVal, err := strconv.ParseInt(token)
	if err != nil {
		if !decok {
			return fmt.Errorf("expected decimal, since we could not parse the integer, but %v is not a decimal", string(token))
		}
		if offset != len(token) {
			return fmt.Errorf("expected decimal %s to parse to end %d, but only parsed up to %d", string(token), len(token), offset)
		}
		return nil
	}
	if !intok {
		return fmt.Errorf("expected float, since we could parse the float and it is not an integer, but %v is not a float", string(token))
	}
	if slowIntVal != intval {
		return fmt.Errorf("want %v got %v", slowIntVal, intval)
	}
	return nil
}
