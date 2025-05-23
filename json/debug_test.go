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
	"encoding/json"
	"fmt"
	"testing"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
	"github.com/katydid/parser-go/compat/downgrade"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/parser/debug"
	"github.com/katydid/parser-go/pool"
)

func TestDebugParse(t *testing.T) {
	p := NewParser()
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	m, err := debug.Parse(p)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Equal(debug.Output) {
		t.Fatalf("expected %s but got %s", debug.Output, m)
	}
}

func TestDebugRandomWalk(t *testing.T) {
	p := &jsonParser{
		Interface: downgrade.ParserWithInit(parse.NewLogger(jsonparse.NewParser())),
		pool:      pool.None(),
	}
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if err := p.Init(data); err != nil {
				t.Fatal(err)
			}
			l := debug.NewLogger(p, debug.NewLineLogger())
			if err := debug.RandomWalk(l, debug.NewRand(), 10, 3); err != nil {
				t.Fatal(err)
			}
		})
	}
}
