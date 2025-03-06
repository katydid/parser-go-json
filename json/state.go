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

import "github.com/katydid/parser-go-json/json/parse"

type state struct {
	kind          stateKind
	arrayElemHint parse.Hint
	arrayIndex    int64
}

type stateKind byte

const atStartStateKind = stateKind(0)

const inLeafStateKind = stateKind('l')

const inArrayIndexStateKind = stateKind('i')
const inArrayAfterIndexStateKind = stateKind('a')

const inObjectAtKeyStateKind = stateKind('k')
const inObjectAtValueStateKind = stateKind('v')

const atEOFStateKind = stateKind('$')
