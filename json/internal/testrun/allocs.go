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

package testrun

import (
	"runtime"
	"testing"
	"time"

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/rand"
)

func NoAllocsOnAverage(t *testing.T, f func(bs []byte)) {
	t.Helper()
	const numValues = 1000
	r := rand.NewRand()
	inputs := rand.Values(r, numValues)
	const runsPerTest = 100

	for i := 0; i < numValues; i++ {
		ff := func() { f(inputs[i]) }
		allocs := testing.AllocsPerRun(runsPerTest, ff)
		if allocs != 0 {
			t.Fatalf("seed = %v, got %v allocs, want 0 allocs", r.Seed(), allocs)
		}
	}
}

func NotASingleAllocAfterWarmUp(t *testing.T, pool pool.Pool, f func(bs []byte)) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	t.Helper()
	r := rand.NewRand()
	seed := r.Seed()
	numValues := 100

	inputs := rand.Values(r, numValues)

	testRunWarmup(pool, f, inputs)

	testRunAfterWarmup(t, seed, pool, f, inputs)
}

func testRunWarmup(pool pool.Pool, f func(bs []byte), inputs [][]byte) {
	// warm up buffer pool
	for i := 0; i < len(inputs); i++ {
		f(inputs[i])
		pool.FreeAll()
	}
}

func testRunAfterWarmup(t *testing.T, seed int64, pool pool.Pool, f func(bs []byte), inputs [][]byte) {
	originalPoolSize := pool.Size()
	for i := 0; i < len(inputs); i++ {
		ff := func() {
			f(inputs[i])
			pool.FreeAll()
		}
		allocs := allocsForSingleRun(ff)
		retries := 10
		for allocs != 0 && retries > 0 {
			// there are sometimes allocations made by the testing framework
			// retry to make sure that the allocation is the parser's fault.
			time.Sleep(100e6)
			allocs = allocsForSingleRun(ff)
			if allocs == 0 {
				break
			}
			retries -= 1
		}
		if allocs != 0 {
			poolallocs := pool.Size() - originalPoolSize
			t.Fatalf("input = %s, seed = %v, got %v allocs, pool allocs = %d", inputs[i], seed, allocs, poolallocs)
		}
	}
}

// allocsForSingleRun should be called with GOMAXPROCS=1
// Hint add the following line to the caller
//
//	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
func allocsForSingleRun(f func()) (avg uint64) {
	var memstats runtime.MemStats

	// Measure the starting statistics
	runtime.ReadMemStats(&memstats)
	mallocs := 0 - memstats.Mallocs

	// Run the function
	f()

	// Read the final statistics
	runtime.ReadMemStats(&memstats)
	mallocs += memstats.Mallocs

	return mallocs
}
