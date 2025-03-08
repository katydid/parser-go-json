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
	"io"
	"testing"
)

func TestParseExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true,1,2],"obj":{"k":"v","a":[1,2,3],"b":1,"c":2}}`
	p := NewParser([]byte(s))
	expect(t, p.Next, ObjectOpenHint)

	expect(t, p.Next, KeyHint)
	expectStr(t, p, "num")

	expect(t, p.Next, ValueHint)
	expectFloat(t, p, 3.14)

	expect(t, p.Next, KeyHint)
	expectStr(t, p, "arr")

	expect(t, p.Next, ArrayOpenHint)

	expect(t, p.Next, ValueHint)

	expect(t, p.Next, ValueHint)
	expectFalse(t, p)

	expect(t, p.Next, ValueHint)
	expectTrue(t, p)

	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}

	expect(t, p.Next, KeyHint)
	expectStr(t, p, "obj")

	expect(t, p.Next, ObjectOpenHint)

	expect(t, p.Next, KeyHint)
	expectStr(t, p, "k")

	expect(t, p.Next, ValueHint)
	expectStr(t, p, "v")

	expect(t, p.Next, KeyHint)
	expectStr(t, p, "a")

	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}

	if err := p.Skip(); err != nil {
		t.Fatal(err)
	}

	expect(t, p.Next, ObjectCloseHint)
	if _, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v", err)
	}
}
