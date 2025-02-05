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
	"testing"

	"github.com/katydid/parser-go-json/json/rand"
)

func TestNoAllocsOnAverage(t *testing.T) {
	num := 100
	r := rand.NewRand()
	values := rand.Values(r, num)
	bs := make([][]byte, len(values))
	for i := range values {
		bs[i] = []byte(values[i])
	}
	scanner := NewScanner(nil)

	const runsPerTest = 100
	for i := 0; i < num; i++ {
		f := func() {
			scanner.Init(bs[i%num])
			_, _, err := scanner.Next()
			for err == nil {
				_, _, err = scanner.Next()
			}
		}
		allocs := testing.AllocsPerRun(runsPerTest, f)
		if allocs != 0 {
			t.Errorf("seed = %v, got %v allocs, want 0 allocs", r.Seed(), allocs)
		}
	}
}

func TestNotASingleAllocAfterWarmUp(t *testing.T) {
	num := 100
	r := rand.NewRand()
	values := rand.Values(r, num)
	bs := make([][]byte, len(values))
	for i := range values {
		bs[i] = []byte(values[i])
	}
	scanner := NewScanner(nil)

	const runsPerTest = 1
	for i := 0; i < num; i++ {
		f := func() {
			scanner.Init(bs[i%num])
			_, _, err := scanner.Next()
			for err == nil {
				_, _, err = scanner.Next()
			}
		}
		allocs := testing.AllocsPerRun(runsPerTest, f)
		if allocs != 0 {
			// there are sometimes allocations made by the testing framework
			// retry to make sure that the allocation is the parser's fault.
			allocs2 := testing.AllocsPerRun(runsPerTest, f)
			if allocs2 != 0 {
				t.Errorf("seed = %v, got %v allocs, want 0 allocs", r.Seed(), allocs)
			}
		}
	}
}
