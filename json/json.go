//  Copyright 2013 Walter Schulze
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

// Package json contains the implementation of a JSON parser.
package json

import (
	"github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/tag"
	goparse "github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/pool"
)

type Parser interface {
	goparse.Parser
	Reset()
	// Init restarts the parser with a new byte buffer, without allocating a new parser.
	Init([]byte)
}

type parserWithReset interface {
	goparse.Parser
	Reset()
}

type jsonParser struct {
	parserWithReset
	underlying Parser
	pool       pool.Pool
}

// NewParser returns a new JSON parser with indexes.
// Use this parser with other Katydid tools, such as the validator.
func NewParser() Parser {
	p := pool.New()
	underlyingParser := parse.NewParser(parse.WithAllocator(p.Alloc))
	tagged := tag.NewTagger(underlyingParser, tag.WithAllocator(p.Alloc), tag.WithIndexes())
	return &jsonParser{parserWithReset: tagged, underlying: underlyingParser, pool: p}
}

// NewJSONSchemaParser returns a new JSON parser that tags objects and arrays, so that the types can be checked by JSONSchema.
// The following json: `{"a": ["b", "c"]}`
// is parsed as: `{"object": {"a": {"array": {0: "b", 1: "c"}}}}`.
// The kind returned from the Token method for "object" and "array" will be parse.TagKind.
func NewJSONSchemaParser() Parser {
	p := pool.New()
	underlyingParser := parse.NewParser(parse.WithAllocator(p.Alloc))
	j := tag.NewTagger(underlyingParser, tag.WithAllocator(p.Alloc), tag.WithIndexes(), tag.WithTags())
	return &jsonParser{parserWithReset: j, pool: p}
}

func (p *jsonParser) Init(buf []byte) {
	// This Init really inits the underlying parser with the new buffer.
	p.parserWithReset.Reset()
	p.underlying.Init(buf)
	p.pool.FreeAll()
	return
}
