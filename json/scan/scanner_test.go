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

func expect[A comparable, B any](t *testing.T, f func() (A, B, error), want A) {
	t.Helper()
	got, _, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func TestScannerExample(t *testing.T) {
	str := `{"num":3.14,"arr":[null,false,true],"obj":{"k":"v"}}`
	s := NewScanner([]byte(str))
	expect(t, s.Next, ObjectOpenKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, NumberKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, ArrayOpenKind)
	expect(t, s.Next, NullKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, FalseKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, TrueKind)
	expect(t, s.Next, ArrayCloseKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, ObjectOpenKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ObjectCloseKind)
	expect(t, s.Next, ObjectCloseKind)
}

func TestScannerExampleWithSpaces(t *testing.T) {
	str := "  {  \"num\" : 3.14\t\r\n ,\t\"arr\"\n:\r[   null , false    , true],  \"obj\" : { \"k\" : \"v\" }, \"boring\"  : [\n 1 , 2 ,  3  ]  }  "
	s := NewScanner([]byte(str))
	expect(t, s.Next, ObjectOpenKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, NumberKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, ArrayOpenKind)
	expect(t, s.Next, NullKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, FalseKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, TrueKind)
	expect(t, s.Next, ArrayCloseKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, ObjectOpenKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ObjectCloseKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, StringKind)
	expect(t, s.Next, ColonKind)
	expect(t, s.Next, ArrayOpenKind)
	expect(t, s.Next, NumberKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, NumberKind)
	expect(t, s.Next, CommaKind)
	expect(t, s.Next, NumberKind)
	expect(t, s.Next, ArrayCloseKind)
	expect(t, s.Next, ObjectCloseKind)
}
