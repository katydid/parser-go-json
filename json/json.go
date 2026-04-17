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
	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go-json/json/tag"
	"github.com/katydid/parser-go/compat/downgrade"
	"github.com/katydid/parser-go/parser"
	"github.com/katydid/parser-go/pool"
)

// Interface is a parser for JSON
type Interface interface {
	parser.Interface
	//Init initialises the parser with a byte buffer containing JSON.
	Init(buf []byte) error
}

type jsonParser struct {
	Interface
	underlying jsonparse.Parser
	pool       pool.Pool
}

// NewParser returns a new JSON parser.
func NewParser() Interface {
	p := pool.New()
	underlyingParser := jsonparse.NewParser(jsonparse.WithAllocator(p.Alloc))
	j := tag.NewTagger(underlyingParser, tag.WithAllocator(p.Alloc), tag.WithIndexes())
	return &jsonParser{
		Interface:  downgrade.ParserWithInit(&resetInit{j}),
		underlying: underlyingParser,
		pool:       p,
	}
}

type resetInit struct {
	tag.Parser
}

func (r *resetInit) Init(buf []byte) {
	r.Reset()
}

func (p *jsonParser) Init(buf []byte) error {
	// This Init only resets the tagger and the downgrade parser.
	p.Interface.Init(buf)
	// This Init really inits the underlying parser with the new buffer.
	p.underlying.Init(buf)
	p.pool.FreeAll()
	return nil
}
