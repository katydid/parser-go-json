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

package json

import (
	"testing"

	"github.com/katydid/parser-go/expect"
	"github.com/katydid/parser-go/parse"
)

func TestJSONSchema(t *testing.T) {
	str := `{"a": ["b", "c"]}`

	// parse as not json schema
	p := NewParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)

	// now parse as json schema
	p = NewJSONSchemaParser()
	p.Init([]byte(str))
	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Tag(t, p, "object")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.String(t, p, "a")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Tag(t, p, "array")

	expect.Hint(t, p, parse.EnterHint)
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 0)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "b")
	expect.Hint(t, p, parse.FieldHint)
	expect.Int(t, p, 1)
	expect.Hint(t, p, parse.ValueHint)
	expect.String(t, p, "c")
	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)

	expect.Hint(t, p, parse.LeaveHint)
	expect.EOF(t, p)
}
