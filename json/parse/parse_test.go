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
)

func TestParseExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true,1,2],"obj":{"k":"v","a":[1,2,3],"b":1,"c":2}}`
	p := NewParser(WithBuffer([]byte(s)))
	expectHint(t, p, ObjectOpenHint)

	expectHint(t, p, KeyHint)
	expectString(t, p, "num")

	expectHint(t, p, ValueHint)
	expectFloat(t, p, 3.14)

	expectHint(t, p, KeyHint)
	expectString(t, p, "arr")

	expectHint(t, p, ArrayOpenHint)

	expectHint(t, p, ValueHint) // null

	expectHint(t, p, ValueHint)
	expectFalse(t, p)

	expectHint(t, p, ValueHint)
	expectTrue(t, p)

	expectNoErr(t, p.Skip) // skip 1,2]

	expectHint(t, p, KeyHint)
	expectString(t, p, "obj")

	expectHint(t, p, ObjectOpenHint)

	expectHint(t, p, KeyHint)
	expectString(t, p, "k")

	expectHint(t, p, ValueHint)
	expectString(t, p, "v")

	expectHint(t, p, KeyHint)
	expectString(t, p, "a")

	expectNoErr(t, p.Skip)

	expectNoErr(t, p.Skip)

	expectHint(t, p, ObjectCloseHint)
	expectEOF(t, p)
}
