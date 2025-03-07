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

package parse

import (
	"fmt"
	"testing"

	"github.com/katydid/parser-go-json/json/token"
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

func expectFalse(t *testing.T, tzer Parser) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != token.FalseKind {
		t.Fatalf("expected false, but got %v", tokenKind)
	}
}

func expectTrue(t *testing.T, tzer Parser) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != token.TrueKind {
		t.Fatalf("expected true, but got %v", tokenKind)
	}
}

func expectInt(t *testing.T, tzer Parser, want int64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != token.Int64Kind {
		t.Fatalf("expected int64, but got %v", tokenKind)
	}
	got := castToInt64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFloat(t *testing.T, tzer Parser, want float64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != token.Float64Kind {
		t.Fatalf("expected float64, but got %v", tokenKind)
	}
	got := castToFloat64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectStr(t *testing.T, tzer Parser, want string) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != token.StringKind {
		t.Fatalf("expected string, but got %v", tokenKind)
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
