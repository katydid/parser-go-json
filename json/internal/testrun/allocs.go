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
	"testing"

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
	t.Helper()
	r := rand.NewRand()
	numValues := 100
	inputs := rand.Values(r, numValues)

	// warm up buffer pool
	for i := 0; i < numValues; i++ {
		f(inputs[i])
		pool.FreeAll()
	}
	originalPoolSize := pool.Size()

	const runsPerTest = 1
	for i := 0; i < numValues; i++ {
		ff := func() {
			f(inputs[i])
			pool.FreeAll()
		}
		allocs := testing.AllocsPerRun(runsPerTest, ff)
		if allocs != 0 {
			poolallocs := pool.Size() - originalPoolSize
			pool.FreeAll()
			// there are sometimes allocations made by the testing framework
			// retry to make sure that the allocation is the parser's fault.
			allocs2 := testing.AllocsPerRun(runsPerTest, ff)
			if allocs2 != 0 {
				t.Fatalf("input = %s, seed = %v, got %v allocs, want 0 allocs, pool allocs = %v", inputs[i], r.Seed(), allocs, poolallocs)
			}
		}
	}
}
