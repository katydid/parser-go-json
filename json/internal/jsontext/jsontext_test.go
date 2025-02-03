// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsontext

import (
	"testing"

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go-json/json/internal/testrun"
)

func TestWalk(t *testing.T) {
	p := json.NewParser()
	testrun.Parsable(t, p, `{"0":0,"1":1} `)
	testrun.Parsable(t, p, `{"0":0,"0":0} `)
	testrun.Parsable(t, p, `{"X":{},"Y":{},"X":{}} `)
}

func TestMoreValues(t *testing.T) {
	p := json.NewParser()
	testrun.EqualValue(t, p, ` null`, "null")
	testrun.EqualValue(t, p, ` null `, "null")
	testrun.SameValue(t, p, `0`)
	testrun.EqualValue(t, p, `0.0`, "0")
	testrun.SameValue(t, p, `123456789`)
	testrun.SameValue(t, p, `0.123456789`)
	testrun.EqualValue(t, p, `0e0`, "0")
	testrun.EqualValue(t, p, `0e+0`, "0")
	testrun.EqualValue(t, p, `0e123456789`, "0") // 0 * 10^123 = 0
	testrun.EqualValue(t, p, `0e+123456789`, "0")
	testrun.EqualValue(t, p, `123.123e+123`, `1.23123e+125`) // 1.23123 x 10^125
	testrun.SameValue(t, p, `123456789.123456789e+123456789`)
	testrun.EqualValue(t, p, `-0`, "0")
	testrun.SameValue(t, p, `-123456789`)
	testrun.EqualValue(t, p, `-0.0`, "0")
	testrun.SameValue(t, p, `-0.123456789`)
	testrun.EqualValue(t, p, `-0e0`, "0")
	testrun.EqualValue(t, p, `-0e-0`, "0")
	testrun.EqualValue(t, p, `-0e123456789`, "0")
	testrun.EqualValue(t, p, `-0e-123456789`, "0")
	testrun.EqualValue(t, p, `-123.123e-123`, `-1.23123e-121`)       // 1.23123 x 10^125
	testrun.EqualValue(t, p, `-123456789.123456789e-123456789`, "0") // (-123456789.123456789) * (10 ^ (-123456789)) = 0
	testrun.EqualValue(t, p, `""`, "")
	testrun.EqualValue(t, p, `"a"`, "a")
	testrun.EqualValue(t, p, `"ab"`, "ab")
	testrun.EqualValue(t, p, `"abc"`, "abc")
	testrun.EqualValue(t, p, `"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"`, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	testrun.EqualValue(t, p, `"\"\\\/\b\f\n\r\t"`, "\"\\/\b\f\n\r\t")
	testrun.EqualValue(t, p, "\"\u0080\u00f6\u20ac\ud799\ue000\ufb33\ufffd\U0001f602\"", "\u0080\u00f6\u20ac\ud799\ue000\ufb33\ufffd\U0001f602")
}

func TestErrors(t *testing.T) {
	p := json.NewParser()
	testrun.NotParsable(t, p, ``)
	testrun.NotParsable(t, p, ` #`)
	testrun.NotParsable(t, p, ` `)
	testrun.NotParsable(t, p, `null null`)
	testrun.NotParsable(t, p, `nul`)
	testrun.NotParsable(t, p, `nulL`)
	testrun.NotParsable(t, p, `fals`)
	testrun.NotParsable(t, p, `falsE`)
	testrun.NotParsable(t, p, `tru`)
	testrun.NotParsable(t, p, `truE`)
	testrun.NotParsable(t, p, `"start`)
	testrun.NotParsable(t, p, `"ok`+"\x00")
	testrun.NotParsable(t, p, `0.`)
	testrun.NotParsable(t, p, `0.e`)
	testrun.NotParsable(t, p, `{`)
	testrun.NotParsable(t, p, `{"0"`)
	testrun.NotParsable(t, p, `{"0":`)
	testrun.NotParsable(t, p, `{"0":0`)
	testrun.NotParsable(t, p, `{"0":0,`)
	testrun.NotParsable(t, p, ` { "fizz" "buzz" } `)
	testrun.NotParsable(t, p, ` { "fizz" , "buzz" } `)
	testrun.NotParsable(t, p, ` { "fizz" # "buzz" } `)
	testrun.NotParsable(t, p, ` { "fizz" : "buzz" "gazz" } `)
	testrun.NotParsable(t, p, ` { "fizz" : "buzz" : "gazz" } `)
	testrun.NotParsable(t, p, ` { "fizz" : "buzz" # "gazz" } `)
	testrun.NotParsable(t, p, ` { , } `)
	testrun.NotParsable(t, p, ` { "fizz" : "buzz" , } `)
	testrun.NotParsable(t, p, ` { null : null } `)
	testrun.NotParsable(t, p, ` { false : false } `)
	testrun.NotParsable(t, p, ` { true : true } `)
	testrun.NotParsable(t, p, ` { 0 : 0 } `)
	testrun.NotParsable(t, p, ` { {} : {} } `)
	testrun.NotParsable(t, p, ` { [] : [] } `)
	testrun.NotParsable(t, p, ` { ] `)
	testrun.NotParsable(t, p, `[`)
	testrun.NotParsable(t, p, `[0`)
	testrun.NotParsable(t, p, `[0,`)
	testrun.NotParsable(t, p, ` [ "fizz" "buzz" ] `)
	testrun.NotParsable(t, p, ` [ } `)
	testrun.NotParsable(t, p, `"",`)
	testrun.NotParsable(t, p, `{:`)
	testrun.NotParsable(t, p, `{"",`)
	testrun.NotParsable(t, p, `{"":`)
	testrun.NotParsable(t, p, `{"":"":`)
	testrun.NotParsable(t, p, `{"":"",`)
	testrun.NotParsable(t, p, `[,`)
	testrun.NotParsable(t, p, `["":`)
	testrun.NotParsable(t, p, `["",`)
}

// Test that JSON parser doesn't break with invalid UTF8
// This might not be the behaviour we want and future we might expect an error.
func TestInvalidUTF8(t *testing.T) {
	p := json.NewParser()
	testrun.Parsable(t, p, "\"living\xde\xad\xbe\xef\"")
	testrun.Parsable(t, p, ` "a`+"\xff"+`0" `)
	testrun.Parsable(t, p, ` [ "a`+"\xff"+`1" ] `)
	testrun.Parsable(t, p, ` [ "a1" , "b`+"\xff"+`1" ] `)
	testrun.Parsable(t, p, ` [ [ "a`+"\xff"+`2" ] ] `)
	testrun.Parsable(t, p, ` [ "a1" , [ "a`+"\xff"+`2" ] ] `)
	testrun.Parsable(t, p, ` [ [ "a2" , "b`+"\xff"+`2" ] ] `)
	testrun.Parsable(t, p, ` [ "a1" , [ "a2" , "b`+"\xff"+`2" ] ] `)
	testrun.Parsable(t, p, ` { "a`+"\xff"+`1" : "b1" } `)
	testrun.Parsable(t, p, ` { "a1" : "b`+"\xff"+`1" } `)
	testrun.Parsable(t, p, ` { "a1" : "b1" , "c`+"\xff"+`1" : "d1" } `)
	testrun.Parsable(t, p, ` { "a1" : "b1" , "c1" : "d`+"\xff"+`1" } `)
	testrun.Parsable(t, p, ` { "a1" : { "a`+"\xff"+`2" : "b2" } } `)
	testrun.Parsable(t, p, ` { "a1" : { "a2" : "b`+"\xff"+`2" } } `)
	testrun.Parsable(t, p, ` { "a1" : { "a2" : "b2" , "c`+"\xff"+`2" : "d2" } } `)
	testrun.Parsable(t, p, ` { "a1" : { "a2" : "b2" , "c2" : "d`+"\xff"+`2" } } `)
	testrun.Parsable(t, p, ` [ "a1" , { "a2" : "b`+"\xff"+`2" } ] `)
	testrun.Parsable(t, p, ` { "a1" : "b1" , "c1" : [ "a2" , "b`+"\xff"+`2" ] } `)
	testrun.Parsable(t, p, ` [ { "a1" : [ "a2" , { "a3" : "b3" , "c3" : [ "a4" , "b`+"\xff"+`4" ] } ] } ] `)
}
