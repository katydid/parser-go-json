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
	"testing"

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestParseExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true,1,2],"obj":{"k":"v","a":[1,2,3],"b":1,"c":2}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.ObjectOpenHint)

	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "num")

	expect.Hint(t, p, parse.ValueHint)
	expect.Float(t, p, 3.14)

	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "arr")

	expect.Hint(t, p, parse.ArrayOpenHint)

	expect.Hint(t, p, parse.ValueHint)

	expect.Hint(t, p, parse.ValueHint)
	expect.False(t, p)

	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "obj")

	expect.Hint(t, p, parse.ObjectOpenHint)

	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "k")

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "v")

	expect.Hint(t, p, parse.KeyHint)
	expect.String(t, p, "a")

	expect.NoErr(t, p.Skip)

	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.ObjectCloseHint)
	expect.EOF(t, p)
}
