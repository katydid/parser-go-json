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
	"testing"
)

func expect[A comparable](t *testing.T, f func() (A, error), want A) {
	t.Helper()
	got, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectInt(t *testing.T, tzer Tokenizer, want int64) {
	t.Helper()
	tokenKind, err := tzer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != Int64Kind {
		t.Fatalf("expected int64, but got %v", tokenKind)
	}
	gotb, err := tzer.Bytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := castToInt64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFloat(t *testing.T, tzer Tokenizer, want float64) {
	t.Helper()
	tokenKind, err := tzer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != Float64Kind {
		t.Fatalf("expected float64, but got %v", tokenKind)
	}
	gotb, err := tzer.Bytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := castToFloat64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectStr(t *testing.T, tzer Tokenizer, want string) {
	t.Helper()
	tokenKind, err := tzer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != StringKind {
		t.Fatalf("expected string, but got %v", tokenKind)
	}
	gotb, err := tzer.Bytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotf := string(gotb)
	got := fmt.Sprintf("%v", gotf)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectErr[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	got, err := f()
	if err == nil {
		t.Fatalf("expected error, but got %v", got)
	}
}
