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

	"github.com/katydid/parser-go-json/json/internal/pool"
	"github.com/katydid/parser-go-json/json/internal/testrun"
)

func TestNoAllocsOnAverage(t *testing.T) {
	pool := pool.New()
	s := NewScanner(nil)
	testrun.NoAllocsOnAverage(t, func(input []byte) {
		s.Init(input)
		if err := walk(s); err != nil {
			t.Fatalf("expected EOF, but got %v", err)
		}
		pool.FreeAll()
	})
}

func TestNotASingleAllocAfterWarmUp(t *testing.T) {
	pool := pool.New()
	s := NewScanner(nil)
	testrun.NotASingleAllocAfterWarmUp(t, pool, func(bs []byte) {
		s.Init(bs)
		if err := walk(s); err != nil {
			t.Fatalf("expected EOF, but got %v", err)
		}
	})
}
