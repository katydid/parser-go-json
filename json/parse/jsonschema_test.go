// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parse

import (
	"testing"

	"github.com/katydid/parser-go-json/json/jsonschema"
	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestJSONSchemaExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true,1,2],"obj":{"k":"v","a":[1,2,3],"b":1,"c":2}}`
	p := NewParser(WithBuffer([]byte(s)))
	expect.Hint(t, p, parse.EnterHint)
	if p.JSONSchemaType() != jsonschema.JSONSchemaTypeObject {
		t.Fatalf("expected object type, but got %v", p.JSONSchemaType())
	}

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "num")

	expect.Hint(t, p, parse.ValueHint)
	expect.Float(t, p, 3.14)

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "arr")

	expect.Hint(t, p, parse.EnterHint)
	if p.JSONSchemaType() != jsonschema.JSONSchemaTypeArray {
		t.Fatalf("expected array type, but got %v", p.JSONSchemaType())
	}

	expect.Hint(t, p, parse.ValueHint) // null

	expect.Hint(t, p, parse.ValueHint)
	expect.False(t, p)

	expect.Hint(t, p, parse.ValueHint)
	expect.True(t, p)

	expect.NoErr(t, p.Skip) // skip 1,2]

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "obj")

	expect.Hint(t, p, parse.EnterHint)
	if p.JSONSchemaType() != jsonschema.JSONSchemaTypeObject {
		t.Fatalf("expected object type, but got %v", p.JSONSchemaType())
	}

	expect.Hint(t, p, parse.FieldHint)
	if p.JSONSchemaType() != jsonschema.JSONSchemaTypeUnknown {
		t.Fatalf("expected unknown type, but got %v", p.JSONSchemaType())
	}
	expect.String(t, p, "k")

	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "v")

	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")

	expect.NoErr(t, p.Skip)

	expect.NoErr(t, p.Skip)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
