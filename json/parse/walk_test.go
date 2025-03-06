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
	"errors"
	"io"

	"github.com/katydid/parser-go-json/json/rand"
	"github.com/katydid/parser-go-json/json/token"
)

var errUnknownToken = errors.New("unknown token")

var errExpectedBool = errors.New("expected bool")

var errExpectedString = errors.New("expected string")

func walkValue(p Parser, hint Hint) error {
	if tokenKind, err := p.Tokenize(); err == nil {
		if tokenKind == token.TrueKind || tokenKind == token.FalseKind || tokenKind == token.NullKind {
			return nil
		}
	}
	if _, err := p.Int(); err == nil {
		return nil
	}
	if _, err := p.Double(); err == nil {
		return nil
	}
	if _, err := p.Bytes(); err == nil {
		return nil
	}
	return errUnknownToken
}

func walk(p Parser) error {
	kind, err := p.Next()
	for err == nil {
		switch kind {
		case ValueHint, KeyHint:
			if err := walkValue(p, kind); err != nil {
				return err
			}
		}
		kind, err = p.Next()
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func randNext(r rand.Rand, p Parser) (Hint, error) {
	skip := r.Intn(2) == 0
	for skip {
		if err := p.Skip(); err != nil {
			return UnknownHint, err
		}
		skip = r.Intn(2) == 0
	}
	return p.Next()
}

func randWalk(r rand.Rand, p Parser) error {
	hint, err := p.Next()
	for err == nil {
		switch hint {
		case ValueHint, KeyHint:
			if err := walkValue(p, hint); err != nil {
				return err
			}
		}
		hint, err = randNext(r, p)
	}
	if err != io.EOF {
		return err
	}
	return nil
}
