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
	expectStr(t, tzer, "num")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.NumberKind)
	expectFloat(t, tzer, 3.14)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "arr")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.ArrayOpenKind)

	expect(t, tzer.Next, scan.NullKind)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.FalseKind)
	expect(t, tzer.Tokenize, FalseKind)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.TrueKind)
	expect(t, tzer.Tokenize, TrueKind)

	expect(t, tzer.Next, scan.ArrayCloseKind)

	expect(t, tzer.Next, scan.CommaKind)

	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "obj")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.ObjectOpenKind)

	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "k")

	expect(t, tzer.Next, scan.ColonKind)

	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "v")

	expect(t, tzer.Next, scan.ObjectCloseKind)

	expect(t, tzer.Next, scan.ObjectCloseKind)
}

func TestTokenizerExampleWithSpaces(t *testing.T) {
	str := "  {  \"num\" : 3.14\t\r\n ,\t\"arr\"\n:\r[   null , false    , true],  \"obj\" : { \"k\" : \"v\" }, \"boring\"  : [\n 1 , 2 ,  3  ]  }  "
	tzer := NewTokenizer([]byte(str))
	expect(t, tzer.Next, scan.ObjectOpenKind)
	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "num")
	expect(t, tzer.Next, scan.ColonKind)
	expect(t, tzer.Next, scan.NumberKind)
	expectFloat(t, tzer, 3.14)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.StringKind)
	expectStr(t, tzer, "arr")
	expect(t, tzer.Next, scan.ColonKind)
	expect(t, tzer.Next, scan.ArrayOpenKind)
	expect(t, tzer.Next, scan.NullKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.FalseKind)
	expect(t, tzer.Tokenize, FalseKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.TrueKind)
	expect(t, tzer.Tokenize, TrueKind)
	expect(t, tzer.Next, scan.ArrayCloseKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.Next, scan.ColonKind)
	expect(t, tzer.Next, scan.ObjectOpenKind)
	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.Next, scan.ColonKind)
	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.Next, scan.ObjectCloseKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.StringKind)
	expect(t, tzer.Next, scan.ColonKind)
	expect(t, tzer.Next, scan.ArrayOpenKind)
	expect(t, tzer.Next, scan.NumberKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.NumberKind)
	expect(t, tzer.Next, scan.CommaKind)
	expect(t, tzer.Next, scan.NumberKind)
	expect(t, tzer.Next, scan.ArrayCloseKind)
	expect(t, tzer.Next, scan.ObjectCloseKind)
}
