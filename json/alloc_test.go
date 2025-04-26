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

package json

import (
	"testing"

	"github.com/katydid/parser-go-json/json/internal/testrun"
	"github.com/katydid/parser-go/parser/debug"
	"github.com/katydid/parser-go/pool"
)

func TestNoAllocsOnAverage(t *testing.T) {
	pool := pool.New()
	p := NewParser()
	testrun.NoAllocsOnAverage(t, func(input []byte) {
		p.Init(input)
		if err := debug.Walk(p); err != nil {
			t.Fatalf("expected EOF, but got %v", err)
		}
		pool.FreeAll()
	})
}

func TestNotASingleAllocAfterWarmUp(t *testing.T) {
	p := NewParser()
	pool := p.(*jsonParser).pool
	testrun.NotASingleAllocAfterWarmUp(t, pool, func(bs []byte) {
		p.Init(bs)
		if err := debug.Walk(p); err != nil {
			t.Fatalf("expected EOF, but got %v", err)
		}
	})
}
