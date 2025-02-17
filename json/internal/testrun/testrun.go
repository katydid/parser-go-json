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

package testrun

import (
	"fmt"
	"testing"

	"github.com/katydid/parser-go/parser"
	"github.com/katydid/parser-go/parser/debug"
)

type InitParser interface {
	parser.Interface
	Init([]byte) error
}

func parse(parser InitParser, s string) (debug.Nodes, error) {
	if err := parser.Init([]byte(s)); err != nil {
		return nil, err
	}
	return debug.Parse(parser)
}

func walk(parser InitParser, s string) error {
	if err := parser.Init([]byte(s)); err != nil {
		return err
	}
	return debug.Walk(parser)
}

func parseValue(parser InitParser, input string) (string, error) {
	jout, err := parse(parser, input)
	if err != nil {
		return "", fmt.Errorf("walk error: %v", err)
	}
	if len(jout) != 1 {
		return "", fmt.Errorf("expected one node, but got %v", jout)
	}
	if len(jout[0].Children) != 0 {
		return "", fmt.Errorf("did not expected any children")
	}
	return jout[0].Label, nil
}

func Parsable(t *testing.T, parser InitParser, s string) {
	t.Helper()
	t.Run("Parsable("+s+")", func(t *testing.T) {
		t.Helper()
		if err := walk(parser, s); err != nil {
			t.Fatal(err)
		}
	})
}

func NotParsable(t *testing.T, parser InitParser, s string) {
	t.Helper()
	t.Run("NotParsable("+s+")", func(t *testing.T) {
		t.Helper()
		parsed, err := parse(parser, s)
		if err != nil {
			return
		}
		t.Fatalf("want error, but got: %v", parsed)
	})
}

func EqualValue(t *testing.T, parser InitParser, input, output string) {
	t.Helper()
	t.Run("EqualValue("+input+","+output+")", func(t *testing.T) {
		t.Helper()
		res, err := parseValue(parser, input)
		if err != nil {
			t.Fatal(err)
		}
		if res != output {
			t.Fatalf("want %q, but got %q", output, res)
		}
	})
}

func SameValue(t *testing.T, parser InitParser, input string) {
	t.Helper()
	t.Run("SameValue("+input+")", func(t *testing.T) {
		t.Helper()
		res, err := parseValue(parser, input)
		if err != nil {
			t.Fatal(err)
		}
		if res != input {
			t.Fatalf("want %q, but got %q", input, res)
		}
	})
}
