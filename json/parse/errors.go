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

import "errors"

var errExpectedValue = errors.New("expected value")

var errExpectedCommaOrCloseBracket = errors.New("expected ',' or ']'")

var errExpectedStringOrCloseCurly = errors.New("expected '\"' or '}'")

var errExpectedColon = errors.New("expected ':'")

var errCannotSkipUnknown = errors.New("cannot Skip before parsing")

var errNotInt = errors.New("not an int")

var errNotFloat = errors.New("not a float")
