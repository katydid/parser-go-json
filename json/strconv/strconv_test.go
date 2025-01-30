// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv_test

import (
	"runtime"
	"strings"
	"testing"

	. "github.com/katydid/parser-go-json/json/strconv"
)

var (
	globalBuf [64]byte
	nextToOne = "1.00000000000000011102230246251565404236316680908203125" + strings.Repeat("0", 10000) + "1"

	mallocTest = []struct {
		count int
		desc  string
		fn    func()
	}{
		{0, `AppendInt(localBuf[:0], 123, 10)`, func() {
			var localBuf [64]byte
			AppendInt(localBuf[:0], 123, 10)
		}},
		{0, `AppendInt(globalBuf[:0], 123, 10)`, func() { AppendInt(globalBuf[:0], 123, 10) }},
		{0, `AppendFloat(localBuf[:0], 1.23, 'g', 5, 64)`, func() {
			var localBuf [64]byte
			AppendFloat(localBuf[:0], 1.23, 'g', 5, 64)
		}},
		{0, `AppendFloat(globalBuf[:0], 1.23, 'g', 5, 64)`, func() { AppendFloat(globalBuf[:0], 1.23, 'g', 5, 64) }},
		// In practice we see 7 for the next one, but allow some slop.
		// Before pre-allocation in appendQuotedWith, we saw 39.
		{0, `ParseFloat("123.45", 64)`, func() { ParseFloat([]byte("123.45")) }},
		{0, `ParseFloat("123.456789123456789", 64)`, func() { ParseFloat([]byte("123.456789123456789")) }},
		{0, `ParseFloat("1.000000000000000111022302462515654042363166809082031251", 64)`, func() {
			ParseFloat([]byte("1.000000000000000111022302462515654042363166809082031251"))
		}},
		{0, `ParseFloat("1.0000000000000001110223024625156540423631668090820312500...001", 64)`, func() {
			ParseFloat([]byte(nextToOne))
		}},
	}
)

var oneMB []byte // Will be allocated to 1MB of random data by TestCountMallocs.

func TestCountMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}
	if runtime.GOMAXPROCS(0) > 1 {
		t.Skip("skipping; GOMAXPROCS>1")
	}
	// Allocate a big messy buffer for AppendQuoteToASCII's test.
	oneMB = make([]byte, 1e6)
	for i := range oneMB {
		oneMB[i] = byte(i)
	}
	for _, mt := range mallocTest {
		allocs := testing.AllocsPerRun(0, mt.fn)
		if max := float64(mt.count); allocs > max {
			t.Errorf("%s: %v allocs, want <=%v", mt.desc, allocs, max)
		}
	}
}

// Sink makes sure the compiler cannot optimize away the benchmarks.
var Sink struct {
	Bool       bool
	Int        int
	Int64      int64
	Uint64     uint64
	Float64    float64
	Complex128 complex128
	Error      error
	Bytes      []byte
}

func TestAllocationsFromBytes(t *testing.T) {
	const runsPerTest = 100
	bytes := struct{ Bool, Number, String, Buffer []byte }{
		Bool:   []byte("false"),
		Number: []byte("123456789"),
		String: []byte("hello, world!"),
		Buffer: make([]byte, 1024),
	}

	checkNoAllocs := func(f func()) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			if allocs := testing.AllocsPerRun(runsPerTest, f); allocs != 0 {
				t.Errorf("got %v allocs, want 0 allocs", allocs)
			}
		}
	}

	t.Run("Atoi", checkNoAllocs(func() {
		Sink.Int, Sink.Error = Atoi(bytes.Number)
	}))
	t.Run("ParseInt", checkNoAllocs(func() {
		Sink.Int64, Sink.Error = ParseInt(bytes.Number)
	}))
	t.Run("ParseUint", checkNoAllocs(func() {
		Sink.Uint64, Sink.Error = ParseUint(bytes.Number)
	}))
	t.Run("ParseFloat", checkNoAllocs(func() {
		Sink.Float64, Sink.Error = ParseFloat(bytes.Number)
	}))
}

func TestErrorPrefixes(t *testing.T) {
	_, errInt := Atoi([]byte("INVALID"))
	_, errFloat := ParseFloat([]byte("INVALID"))
	_, errInt64 := ParseInt([]byte("INVALID"))
	_, errUint64 := ParseUint([]byte("INVALID"))

	vectors := []struct {
		err  error  // Input error
		want string // Function name wanted
	}{
		{errInt, "Atoi"},
		{errFloat, "ParseFloat"},
		{errInt64, "ParseInt"},
		{errUint64, "ParseUint"},
	}

	for _, v := range vectors {
		nerr, ok := v.err.(*NumError)
		if !ok {
			t.Errorf("test %s, error was not a *NumError", v.want)
			continue
		}
		if got := nerr.Func; got != v.want {
			t.Errorf("mismatching Func: got %s, want %s", got, v.want)
		}
	}

}
