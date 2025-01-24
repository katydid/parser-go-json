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

package json_test

import (
	"fmt"
	"io"

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
	v, err = p.Bytes()
	if err == nil {
		return v
	}
	return nil
}

func walk(p parser.Interface) (debug.Nodes, error) {
	a := make(debug.Nodes, 0)
	for {
		if err := p.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		value := getValue(p)
		if p.IsLeaf() {
			a = append(a, debug.Node{Label: fmt.Sprintf("%v", value), Children: nil})
		} else {
			name := fmt.Sprintf("%v", value)
			p.Down()
			v, err := walk(p)
			if err != nil {
				return nil, err
			}
			p.Up()
			a = append(a, debug.Node{Label: name, Children: v})
		}
	}
	return a, nil
}
