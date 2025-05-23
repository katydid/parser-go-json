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

	"github.com/katydid/parser-go-json/json/internal/testrun"
	"github.com/katydid/parser-go-json/json/rand"
)

func TestParseRandomValues(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := NewParser(WithBuffer([]byte(value)))
			if err := walk(tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v", err)
			}
		})
	}
}

func TestRandomlyParseRandomValues(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := NewParser(WithBuffer([]byte(value)))
			if err := randWalk(r, tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v", err)
			}
		})
	}
}
