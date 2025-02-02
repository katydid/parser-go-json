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
	"io"
	"testing"

	"github.com/katydid/parser-go/parser"
	"github.com/katydid/parser-go/parser/debug"
)

func getValue(p parser.Interface) interface{} {
	var v interface{}
	var err error
	v, err = p.Int()
	if err == nil {
		return v
	}
	v, err = p.Uint()
	if err == nil {
		return v
	}
	v, err = p.Double()
	if err == nil {
		return v
	}
	v, err = p.Bool()
	if err == nil {
		return v
	}
	v, err = p.String()
	if err == nil {
		return v
	}
	bs, err := p.Bytes()
	if err == nil {
		return string(bs)
	}
	return nil
}

func parse(p parser.Interface) (debug.Nodes, error) {
	a := make(debug.Nodes, 0)
	for {
		fmt.Printf("Next\n")
		if err := p.Next(); err != nil {
			if err == io.EOF {
				fmt.Printf("EOF\n")
				break
			} else {
				fmt.Printf("err = %v\n", err)
				return nil, err
			}
		}
		value := getValue(p)
		if p.IsLeaf() {
			fmt.Printf("IsLeaf\n")
			a = append(a, debug.Node{Label: fmt.Sprintf("%v", value), Children: nil})
		} else {
			name := fmt.Sprintf("%v", value)
			fmt.Printf("Down\n")
			p.Down()
			v, err := parse(p)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Up\n")
			p.Up()
			a = append(a, debug.Node{Label: name, Children: v})
		}
	}
	return a, nil
}

func walkValue(p parser.Interface) {
	if _, err := p.Int(); err == nil {
		return
	}
	if _, err := p.Uint(); err == nil {
		return
	}
	if _, err := p.Double(); err == nil {
		return
	}
	if _, err := p.Bool(); err == nil {
		return
	}
	if _, err := p.String(); err == nil {
		return
	}
	if _, err := p.Bytes(); err == nil {
		return
	}
	return
}

func walk(p parser.Interface) error {
	for {
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		walkValue(p)
		if !p.IsLeaf() {
			p.Down()
			if err := walk(p); err != nil {
				return err
			}
			p.Up()
		}
	}
	return nil
}

func TestDebug(t *testing.T) {
	p := NewJsonParser()
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Init(data); err != nil {
		t.Fatal(err)
	}
	m := debug.Walk(p)
	if !m.Equal(debug.Output) {
		t.Fatalf("expected %s but got %s", debug.Output, m)
	}
}

func TestRandomDebug(t *testing.T) {
	p := NewJsonParser()
	data, err := json.Marshal(debug.Input)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if err := p.Init(data); err != nil {
			t.Fatal(err)
		}
		//l := debug.NewLogger(p, debug.NewLineLogger())
		debug.RandomWalk(p, debug.NewRand(), 10, 3)
		//t.Logf("original %v vs random %v", debug.Output, m)
	}
}
