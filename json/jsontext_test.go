// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Original these tests were copied from https://github.com/go-json-experiment/json/blob/master/jsontext/decode_test.go
package json_test

import "testing"

func TestWalk(t *testing.T) {
	testWalk(t, `{"0":0,"1":1} `)
	testWalk(t, `{"0":0,"0":0} `)
	testWalk(t, `{"X":{},"Y":{},"X":{}} `)
}

func TestMoreValues(t *testing.T) {
	testValue(t, ` null`, "null")
	testValue(t, ` null `, "null")
	testSame(t, `0`)
	testValue(t, `0.0`, "0")
	testSame(t, `123456789`)
	testSame(t, `0.123456789`)
	testValue(t, `0e0`, "0")
	testValue(t, `0e+0`, "0")
	testValue(t, `0e123456789`, "0") // 0 * 10^123 = 0
	testValue(t, `0e+123456789`, "0")
	testValue(t, `123.123e+123`, `1.23123e+125`) // 1.23123 x 10^125
	testSame(t, `123456789.123456789e+123456789`)
	testValue(t, `-0`, "0")
	testSame(t, `-123456789`)
	testValue(t, `-0.0`, "0")
	testSame(t, `-0.123456789`)
	testValue(t, `-0e0`, "0")
	testValue(t, `-0e-0`, "0")
	testValue(t, `-0e123456789`, "0")
	testValue(t, `-0e-123456789`, "0")
	testValue(t, `-123.123e-123`, `-1.23123e-121`)       // 1.23123 x 10^125
	testValue(t, `-123456789.123456789e-123456789`, "0") // (-123456789.123456789) * (10 ^ (-123456789)) = 0
	testValue(t, `""`, "")
	testValue(t, `"a"`, "a")
	testValue(t, `"ab"`, "ab")
	testValue(t, `"abc"`, "abc")
	testValue(t, `"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"`, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	// testSame(t, `"\"\\\/\b\f\n\r\t"`)
	// testSame(t, `"\u0022\u005c\u002f\u0008\u000c\u000a\u000d\u0009"`)
	// testSame(t, `"\ud800\udead"`)
	testValue(t, "\"\u0080\u00f6\u20ac\ud799\ue000\ufb33\ufffd\U0001f602\"", "\u0080\u00f6\u20ac\ud799\ue000\ufb33\ufffd\U0001f602")
	// testSame(t, `"\u0080\u00f6\u20ac\ud799\ue000\ufb33\ufffd\ud83d\ude02"`)
}

func TestErrors(t *testing.T) {
	testError(t, ``)
	testError(t, ` #`)
	testError(t, ` `)
	// testError(t, `null null`)
	testError(t, `nul`)
	testError(t, `nulL`)
	testError(t, `fals`)
	testError(t, `falsE`)
	testError(t, `tru`)
	testError(t, `truE`)
	testError(t, `"start`)
	testError(t, `"ok`+"\x00")
	// testError(t, `0.`)
	// testError(t, `0.e`)
	testError(t, `{`)
	testError(t, `{"0"`)
	// testError(t, `{"0":`)
	testError(t, `{"0":0`)
	testError(t, `{"0":0,`)
	testError(t, ` { "fizz" "buzz" } `)
	testError(t, ` { "fizz" , "buzz" } `)
	testError(t, ` { "fizz" # "buzz" } `)
	testError(t, ` { "fizz" : "buzz" "gazz" } `)
	testError(t, ` { "fizz" : "buzz" : "gazz" } `)
	testError(t, ` { "fizz" : "buzz" # "gazz" } `)
	testError(t, ` { , } `)
	// testError(t, ` { "fizz" : "buzz" , } `)
	testError(t, ` { null : null } `)
	testError(t, ` { false : false } `)
	testError(t, ` { true : true } `)
	testError(t, ` { 0 : 0 } `)
	testError(t, ` { {} : {} } `)
	testError(t, ` { [] : [] } `)
	testError(t, ` { ] `)
	testError(t, `[`)
	testError(t, `[0`)
	testError(t, `[0,`)
	testError(t, ` [ "fizz" "buzz" ] `)
	testError(t, ` [ } `)
	// testError(t, `"",`)
	testError(t, `{:`)
	testError(t, `{"",`)
	// testError(t, `{"":`)
	testError(t, `{"":"":`)
	testError(t, `{"":"",`)
	testError(t, `[,`)
	testError(t, `["":`)
	testError(t, `["",`)
}

func DisabledTestInvalidUTF8(t *testing.T) {
	testError(t, "\"living\xde\xad\xbe\xef\"")
	testError(t, ` "a`+"\xff"+`0" `)
	testError(t, ` [ "a`+"\xff"+`1" ] `)
	testError(t, ` [ "a1" , "b`+"\xff"+`1" ] `)
	testError(t, ` [ [ "a`+"\xff"+`2" ] ] `)
	testError(t, ` [ "a1" , [ "a`+"\xff"+`2" ] ] `)
	testError(t, ` [ [ "a2" , "b`+"\xff"+`2" ] ] `)
	testError(t, ` [ "a1" , [ "a2" , "b`+"\xff"+`2" ] ] `)
	testError(t, ` { "a`+"\xff"+`1" : "b1" } `)
	testError(t, ` { "a1" : "b`+"\xff"+`1" } `)
	testError(t, ` { "a1" : "b1" , "c`+"\xff"+`1" : "d1" } `)
	testError(t, ` { "a1" : "b1" , "c1" : "d`+"\xff"+`1" } `)
	testError(t, ` { "a1" : { "a`+"\xff"+`2" : "b2" } } `)
	testError(t, ` { "a1" : { "a2" : "b`+"\xff"+`2" } } `)
	testError(t, ` { "a1" : { "a2" : "b2" , "c`+"\xff"+`2" : "d2" } } `)
	testError(t, ` { "a1" : { "a2" : "b2" , "c2" : "d`+"\xff"+`2" } } `)
	testError(t, ` [ "a1" , { "a2" : "b`+"\xff"+`2" } ] `)
	testError(t, ` { "a1" : "b1" , "c1" : [ "a2" , "b`+"\xff"+`2" ] } `)
	testError(t, ` [ { "a1" : [ "a2" , { "a3" : "b3" , "c3" : [ "a4" , "b`+"\xff"+`4" ] } ] } ] `)
}
