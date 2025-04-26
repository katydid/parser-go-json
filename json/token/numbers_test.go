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

package token

import (
	"fmt"
	"math"
	"testing"

	"github.com/katydid/parser-go-json/json/internal/cast"
	"github.com/katydid/parser-go/parse"
)

func TestNumbersMaxInt64(t *testing.T) {
	input := "9223372036854775807" // math.MaxInt64
	var want int64 = math.MaxInt64
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	_, got, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	goti := cast.ToInt64(got)
	if goti != want {
		t.Fatalf("got %v, but want %v", goti, want)
	}
}

func TestNumbersMaxInt64Plus1(t *testing.T) {
	input := "9223372036854775808" // math.MaxInt64 + 1
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, got, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.DecimalKind {
		t.Fatalf("expected decimal")
	}
	var want string = input
	if string(got) != want {
		t.Fatalf("got %v, but want %v", string(got), want)
	}
}

func TestNumbersMinInt64(t *testing.T) {
	input := "-9223372036854775808" // math.MinInt64
	var want int64 = math.MinInt64
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	_, got, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	goti := cast.ToInt64(got)
	if goti != want {
		t.Fatalf("got %v, but want %v", goti, want)
	}
}

func TestNumbersMinInt64Min1(t *testing.T) {
	input := "-9223372036854775809" // math.MinInt64 - 1
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.DecimalKind {
		t.Fatalf("expected decimal")
	}
	want := input
	got := string(gotb)
	if got != want {
		t.Fatalf("got %v, but want %v", got, want)
	}
}

func TestNumbersMaxUint64(t *testing.T) {
	input := "18446744073709551615" // math.MaxUint64
	var want string = input
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	_, got, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("got %v, but want %v", string(got), want)
	}
}

func TestNumbersMaxUint64Plus1(t *testing.T) {
	input := "18446744073709551616" // math.MaxUint64 + 1
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.DecimalKind {
		t.Fatalf("expected decimal")
	}
	want := input
	got := string(gotb)
	if got != want {
		t.Fatalf("got %v, but want %v", got, want)
	}
}

func TestNumbersMaxFloat64(t *testing.T) {
	input := "1.79769313486231570814527423731704356798070e+308" // math.MaxFloat64
	var want float64 = math.MaxFloat64
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, got, err := tzer.Token()
	if tokenKind != parse.Float64Kind {
		t.Fatalf("expected float64")
	}
	gotf := cast.ToFloat64(got)
	if gotf != want {
		t.Fatalf("got %v, but want %v", gotf, want)
	}
}

func TestNumbersLargerThanMaxFloat64(t *testing.T) {
	input := "2.79769313486231570814527423731704356798070e+308" // > math.MaxFloat64
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.DecimalKind {
		t.Fatal("expected decimal")
	}
	want := input
	got := string(gotb)
	if got != want {
		t.Fatalf("got %v, but want %v", got, want)
	}
}

func TestNumbersSmallestNonZeroFloat64(t *testing.T) {
	input := "4.9406564584124654417656879286822137236505980e-324" // math.SmallestNonzeroFloat64
	var want float64 = math.SmallestNonzeroFloat64
	tzer := NewTokenizer([]byte(input))
	kind, err := tzer.Next()
	if err != nil {
		t.Fatal(err)
	}
	if !kind.IsNumber() {
		t.Fatal("expected number")
	}
	tokenKind, got, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.Float64Kind {
		t.Fatal("expected float64")
	}
	gotf := cast.ToFloat64(got)
	if gotf != want {
		t.Fatalf("got %v, but want %v", gotf, want)
	}
}

func TestNumbersIntOutsideOfFloatingPointPrecision(t *testing.T) {
	inputs := []string{
		"9007199254740993",  // 2^53 + 1
		"9007199254740994",  // 2^53 + 2
		"18014398509481984", // 2^54
		"72057594037927936", // 2^56
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			tzer := NewTokenizer([]byte(input))
			kind, err := tzer.Next()
			if err != nil {
				t.Fatal(err)
			}
			if !kind.IsNumber() {
				t.Fatal("expected number")
			}
			want := input
			tokenKind, gotb, err := tzer.Token()
			if err != nil {
				t.Fatal(err)
			}
			if tokenKind != parse.Int64Kind {
				t.Fatal("expected int64")
			}
			goti := cast.ToInt64(gotb)
			got := fmt.Sprintf("%v", goti)
			if got != want {
				t.Fatalf("got %v, but want %v", got, want)
			}
		})
	}
}

func TestNumbersUintOutsideOfFloatingPointPrecision(t *testing.T) {
	inputs := []string{
		"9007199254740993",  // 2^53 + 1
		"9007199254740994",  // 2^53 + 2
		"18014398509481984", // 2^54
		"72057594037927936", // 2^56
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			tzer := NewTokenizer([]byte(input))
			kind, err := tzer.Next()
			if err != nil {
				t.Fatal(err)
			}
			if !kind.IsNumber() {
				t.Fatal("expected number")
			}
			want := input
			tokenKind, gotb, err := tzer.Token()
			if err != nil {
				t.Fatal(err)
			}
			if tokenKind != parse.Int64Kind {
				t.Fatal("expected int64")
			}
			goti := cast.ToInt64(gotb)
			got := fmt.Sprintf("%v", goti)
			if got != want {
				t.Fatalf("got %v, but want %v", got, want)
			}
		})
	}
}
