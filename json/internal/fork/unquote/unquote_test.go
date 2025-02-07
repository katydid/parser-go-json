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

package unquote

import (
	"bytes"
	"testing"
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
	alloc := func(size int) []byte { return make([]byte, size) }
	for _, test := range unquotetests {
		got, gotOk := unquoteBytes(alloc, []byte(test.in))
		if !bytes.Equal(got, []byte(test.out)) || !gotOk {
			t.Errorf("Unquote(%q) = (%q, %v), want (%q, %v)", test.in, got, gotOk, test.out, true)
		}
	}
}
