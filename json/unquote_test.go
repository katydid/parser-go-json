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
	"bytes"
	"testing"

	"github.com/katydid/parser-go-json/json/pool"
)

var unquotetests = []struct {
	in  string
	out string
}{
	{`""`, ""},
	{`"a"`, "a"},
	{`"abc"`, "abc"},
	{`"☺"`, "☺"},
	{`"hello world"`, "hello world"},
	{`"\u1234"`, "\u1234"},
	{`"'"`, "'"},
}

func TestUnquote(t *testing.T) {
	for _, test := range unquotetests {
		got, gotOk := unquoteBytes(pool.New(), []byte(test.in))
		if !bytes.Equal(got, []byte(test.out)) || !gotOk {
			t.Errorf("Unquote(%q) = (%q, %v), want (%q, %v)", test.in, got, gotOk, test.out, true)
		}
	}
}
