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
	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go/parser/debug"
)

func TestParseRandom(t *testing.T) {
	r := rand.NewRand()
	numValues := 10000
	values := rand.Values(r, numValues)
	p := NewParser()

	for i := 0; i < numValues; i++ {
		name := testrun.Name(values[i])
		t.Run(name, func(t *testing.T) {
			if err := p.Init(values[i]); err != nil {
				t.Fatalf("seed = %v, err = %v, input = %v", r.Seed(), err, string(values[i]))
			}
			if err := debug.Walk(p); err != nil {
				t.Fatalf("seed = %v, err = %v", r.Seed(), err)
			}
		})
	}
}
