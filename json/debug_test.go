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

	"katydid.org.za/go/parser-go/hedge"
	"katydid.org.za/go/parser-go/parse"
	"katydid.org.za/go/parser-go/parse/debug"
	"katydid.org.za/go/parser-go/parse/example"
	"katydid.org.za/go/parser-go/parse/log"
	"katydid.org.za/go/parser-go/pool"
	"katydid.org.za/go/parser-go/rand"
)

func TestDebugParse(t *testing.T) {
	var p parse.ParserWithInit = NewParser()
	data, err := json.Marshal(example.Input)
	if err != nil {
		t.Fatal(err)
	}
	p.Init(data)
	got, err := hedge.ParseInto(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := output.VerboseEqual(got); err != nil {
		t.Fatalf("want %v got %v: %v", output, got, err)
	}
}

var output = hedge.Hedge{
	{Label: hedge.NewStringToken(`A`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(1), Children: nil}}},
	{Label: hedge.NewStringToken(`B`), Children: hedge.Hedge{
		{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b2"), Children: nil}}},
		{Label: hedge.NewInt64Token(1), Children: hedge.Hedge{{Label: hedge.NewStringToken("b3"), Children: nil}}},
	}},
	{Label: hedge.NewStringToken(`C`), Children: hedge.Hedge{
		{Label: hedge.NewStringToken(`A`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(2), Children: nil}}},
		{Label: hedge.NewStringToken(`D`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(3), Children: nil}}},
		{Label: hedge.NewStringToken(`E`), Children: hedge.Hedge{
			{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{
				{Label: hedge.NewStringToken(`A`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(0), Children: nil}}},
				{Label: hedge.NewStringToken(`B`), Children: hedge.Hedge{
					{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b4"), Children: nil}}},
				}},
			}},
			{Label: hedge.NewInt64Token(1), Children: hedge.Hedge{
				{Label: hedge.NewStringToken(`A`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(0), Children: nil}}},
				{Label: hedge.NewStringToken(`B`), Children: hedge.Hedge{
					{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewStringToken("b5"), Children: nil}}},
				}},
			}},
		}},
	}},
	{Label: hedge.NewStringToken(`D`), Children: hedge.Hedge{{Label: hedge.NewInt64Token(4), Children: nil}}},
	{Label: hedge.NewStringToken(`F`), Children: hedge.Hedge{
		{Label: hedge.NewInt64Token(0), Children: hedge.Hedge{{Label: hedge.NewInt64Token(5), Children: nil}}},
	}},
}

func TestDebugRandomWalk(t *testing.T) {
	p := NewParser()
	p.(*jsonParser).pool = pool.None()
	data, err := json.Marshal(example.Input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			p.Init(data)
			l := log.WrapParser(p)
			if err := debug.RandomWalk(l, rand.NewRand(), 10, 3); err != nil {
				t.Fatal(err)
			}
		})
	}
}
