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

package token

import (
	"testing"

	"github.com/katydid/parser-go-json/json/scan"
)

func TestTokenizerExample(t *testing.T) {
	s := `{"num":3.14,"arr":[null,false,true],"obj":{"k":"v"}}`
	tzer := NewTokenizer([]byte(s))
	expect(t, tzer.Next, scan.ObjectOpenKind)

	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.String, "num")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.NumberKind)
	expectErr(t, tzer.Int)
	expect(t, tzer.Double, 3.14)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.String, "arr")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.ArrayOpenKind)

	expect(t, tzer.Next, scan.NullKind)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.FalseKind)
	expect(t, tzer.Bool, false)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.TrueKind)
	expect(t, tzer.Bool, true)

	expect(t, tzer.Next, scan.ArrayCloseKind)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.String, "obj")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.ObjectOpenKind)

	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.String, "k")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.String, "v")

	expect(t, tzer.Next, scan.ObjectCloseKind)

	expect(t, tzer.Next, scan.ObjectCloseKind)
}

func TestTokenizerExampleWithSpaces(t *testing.T) {
	str := "  {  \"num\" : 3.14\t\r\n ,\t\"arr\"\n:\r[   null , false    , true],  \"obj\" : { \"k\" : \"v\" }, \"boring\"  : [\n 1 , 2 ,  3  ]  }  "
	s := NewTokenizer([]byte(str))
	expect(t, s.Next, scan.ObjectOpenKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.String, "num")
	expect(t, s.Next, scan.ColonKind)
	expect(t, s.Next, scan.NumberKind)
	expectErr(t, s.Int)
	expect(t, s.Double, 3.14)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.String, "arr")
	expect(t, s.Next, scan.ColonKind)
	expect(t, s.Next, scan.ArrayOpenKind)
	expect(t, s.Next, scan.NullKind)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.FalseKind)
	expect(t, s.Bool, false)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.TrueKind)
	expect(t, s.Bool, true)
	expect(t, s.Next, scan.ArrayCloseKind)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.Next, scan.ColonKind)
	expect(t, s.Next, scan.ObjectOpenKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.Next, scan.ColonKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.Next, scan.ObjectCloseKind)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.StringKind)
	expect(t, s.Next, scan.ColonKind)
	expect(t, s.Next, scan.ArrayOpenKind)
	expect(t, s.Next, scan.NumberKind)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.NumberKind)
	expect(t, s.Next, scan.CommaKind)
	expect(t, s.Next, scan.NumberKind)
	expect(t, s.Next, scan.ArrayCloseKind)
	expect(t, s.Next, scan.ObjectCloseKind)
}
