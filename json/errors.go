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

package json

import "errors"

var errLongBuffer = errors.New("json is complete, but buffer is still has more to read")

// errUnquote returns an error that resulted from trying to unquote a string.
var errUnquote = errors.New("unable to unquote string")

var errExpectedOpenCurly = errors.New("expected '{'")

var errExpectedOpenBracket = errors.New("expected '['")

var errExpectedCloseBracket = errors.New("expected ']'")

var errExpectedComma = errors.New("expected ','")

var errExpectedCloseCurly = errors.New("expected '}'")

var errExpectedColon = errors.New("expected ':'")

var errExpectedValue = errors.New("expected a json value")

var errExpectedCommaOrCloseCurly = errors.New("expected ',' or '}'")

var errExpectedCommaOrCloseBracket = errors.New("expected ',' or ']'")

var errNotLeaf = errors.New("not leaf")
