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

package tag_test

import (
	"testing"

	"github.com/katydid/parser-go-json/json/internal/testrun"
	"github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go-json/json/tag"
)

func TestParseRandomValuesWithTagsAndIndexes(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithTags(), tag.WithIndexes())
			if err := walk(tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestRandomlyParseRandomValuesWithTagsAndIndexes(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithTags(), tag.WithIndexes())
			if err := randWalk(r, tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestParseRandomValuesWithTagsOnly(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithTags())
			if err := walk(tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestRandomlyParseRandomValuesWithIndexesOnly(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithIndexes())
			if err := randWalk(r, tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestParseRandomValuesWithIndexesOnly(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithIndexes())
			if err := walk(tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestRandomlyParseRandomValuesWithTagsOnly(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))), tag.WithTags())
			if err := randWalk(r, tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestParseRandomValuesWithoutTagsAndIndexes(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))))
			if err := walk(tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}

func TestRandomlyParseRandomValuesWithoutTagsAndIndexes(t *testing.T) {
	r := rand.NewRand()
	values := rand.Values(r, 100)
	for _, value := range values {
		name := testrun.Name(value)
		t.Run(name, func(t *testing.T) {
			tokenizer := tag.NewTagger(parse.NewParser(parse.WithBuffer([]byte(value))))
			if err := randWalk(r, tokenizer); err != nil {
				t.Fatalf("expected EOF, but got %v using seed %v", err, r.Seed())
			}
		})
	}
}
