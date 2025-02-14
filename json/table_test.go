//  Copyright 2015 Walter Schulze
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
	"fmt"
	"testing"

	"github.com/katydid/parser-go-json/json/internal/testrun"
)

func TestValues(t *testing.T) {
	p := NewParser()
	testrun.EqualValue(t, p, "0", "0")
	testrun.EqualValue(t, p, "1", "1")
	testrun.EqualValue(t, p, "-1", "-1")
	testrun.EqualValue(t, p, "123", "123")
	testrun.EqualValue(t, p, "1.1", "1.1")
	testrun.EqualValue(t, p, "1.123", "1.123")
	testrun.EqualValue(t, p, "1.1e1", "11")
	testrun.EqualValue(t, p, "1.1e-1", "0.11")
	testrun.EqualValue(t, p, "1.1e10", "11000000000")
	testrun.EqualValue(t, p, "1.1e+10", "11000000000")
	testrun.EqualValue(t, p, `"a"`, "a")
	testrun.EqualValue(t, p, `"abc"`, "abc")
	testrun.EqualValue(t, p, `""`, "")
	testrun.EqualValue(t, p, `"\b"`, "\b")
	testrun.EqualValue(t, p, `true`, "true")
	testrun.EqualValue(t, p, `false`, "false")
	testrun.EqualValue(t, p, `null`, fmt.Sprintf("%v", []byte("null")))
}

func TestArray(t *testing.T) {
	p := NewParser()
	testrun.Parsable(t, p, `[]`)
	testrun.NotParsable(t, p, `[`)
	testrun.Parsable(t, p, `[1]`)
	testrun.NotParsable(t, p, `[1 2]`)
	testrun.Parsable(t, p, `[1,2]`)
	testrun.Parsable(t, p, `[1,2.3e5]`)
	testrun.Parsable(t, p, `[1,"a"]`)
	testrun.Parsable(t, p, `[1,2,3]`)
	testrun.Parsable(t, p, `[true,false,null]`)
	testrun.Parsable(t, p, `[ true  , false , null   ]`)
	testrun.Parsable(t, p, `[{"a": true, "b": [1,2]}]`)
	testrun.Parsable(t, p, `["["]`)
	testrun.Parsable(t, p, `["]"]`)
}

func TestObject(t *testing.T) {
	p := NewParser()
	testrun.Parsable(t, p, `{}`)
	testrun.Parsable(t, p, `{"a":1}`)
	testrun.Parsable(t, p, `{"a":1,"b":2}`)
	testrun.Parsable(t, p, `{"a":1,"b":2,"c":3}`)
	testrun.NotParsable(t, p, `{"a":1,"b"}`)
	testrun.NotParsable(t, p, `{"a"}`)
	testrun.NotParsable(t, p, `{"a" "b"}`)
	testrun.Parsable(t, p, `{"{":null}`)
	testrun.Parsable(t, p, `{"}":null}`)
	testrun.Parsable(t, p, `{"a":"{"}`)
	testrun.Parsable(t, p, `{"a":"}"}`)
	testrun.Parsable(t, p, `{"a":true,"b":false}`)
	testrun.Parsable(t, p, `{"a": true , "b": false}`)
	testrun.Parsable(t, p, `{"a":[1]}`)
	testrun.Parsable(t, p, `{"a":true,"b":[1,2]}`)
	testrun.Parsable(t, p, `{"a": true, "b": [1,2]}`)
}
