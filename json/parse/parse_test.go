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

func expectErr[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	got, err := f()
	if err == nil {
		t.Fatalf("expected error, but got %v", got)
	}
}

func TestParseExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true],"obj":{"k":"v", "boring": [1,2,3]}}`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenKind)

	expect(t, p.Next, StringKind)
	expect(t, p.String, "num")

	expect(t, p.Next, NumberKind)
	expectErr(t, p.Int)
	expectErr(t, p.Uint)
	expect(t, p.Double, 3.14)

	expect(t, p.Next, StringKind)
	expect(t, p.String, "arr")

	expect(t, p.Next, ArrayOpenKind)

	expect(t, p.Next, NullKind)

	expect(t, p.Next, FalseKind)
	expect(t, p.Bool, false)

	expect(t, p.Next, TrueKind)
	expect(t, p.Bool, true)

	expect(t, p.Next, ArrayCloseKind)

	expect(t, p.Next, StringKind)
	expect(t, p.String, "obj")

	expect(t, p.Next, ObjectOpenKind)

	expect(t, p.Next, StringKind)
	expect(t, p.String, "k")

	expect(t, p.Next, StringKind)
	expect(t, p.String, "v")

	expect(t, p.Next, StringKind)
	expect(t, p.String, "boring")

	p.Skip()

	expect(t, p.Next, ObjectCloseKind)

	expect(t, p.Next, ObjectCloseKind)
}
