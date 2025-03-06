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
	"math"
	"testing"
)

func TestCastInt64(t *testing.T) {
	want := int64(123)
	bs := castFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestCastMaxInt64(t *testing.T) {
	want := int64(math.MaxInt64)
	bs := castFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestCastMinInt64(t *testing.T) {
	want := int64(math.MinInt64)
	bs := castFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestUnsafeCastInt64(t *testing.T) {
	want := int64(123)
	bs := unsafeCastFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestUnsafeCastMaxInt64(t *testing.T) {
	want := int64(math.MaxInt64)
	bs := unsafeCastFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestUnsafeCastMinInt64(t *testing.T) {
	want := int64(math.MinInt64)
	bs := unsafeCastFromInt64(want)
	got := castToInt64(bs)
	if got != want {
		t.Fatalf("want %d got %d", want, got)
	}
}

func TestAllocCastUnsafe(t *testing.T) {
	f := func() {
		want := int64(1233456578)
		bs := unsafeCastFromInt64(want)
		got := castToInt64(bs)
		if got != want {
			t.Fatalf("want %d got %d", want, got)
		}
	}
	for i := 0; i < 10000; i++ {
		allocs := testing.AllocsPerRun(1, f)
		if allocs > 0 {
			t.Fatalf("UnsafeCast Allocs = %f", allocs)
		}
	}
}

func TestAllocCastDeprecated(t *testing.T) {
	f := func() {
		want := int64(1233456578)
		bs := castFromInt64(want)
		got := castToInt64(bs)
		if got != want {
			t.Fatalf("want %d got %d", want, got)
		}
	}
	for i := 0; i < 10000; i++ {
		allocs := testing.AllocsPerRun(1, f)
		if allocs > 0 {
			t.Fatalf("Cast Allocs = %f", allocs)
		}
	}
}
