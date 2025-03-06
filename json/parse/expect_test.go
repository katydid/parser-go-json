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
	"bytes"
	"io"
	"testing"
)

func expect[A comparable](t *testing.T, f func() (A, error), want A) {
	t.Helper()
	got, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectStr(t *testing.T, f func() ([]byte, error), want string) {
	t.Helper()
	got, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("want %v, but got %v", want, string(got))
	}
}

func expectErr[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	got, err := f()
	if err == nil || err == io.EOF {
		t.Fatalf("expected error, but got %v with err = %v", got, err)
	}
}
