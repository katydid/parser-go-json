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

package tag

import "github.com/katydid/parser-go/parse"

type state struct {
	kind          stateKind
	arrayElemHint parse.Hint
	arrayIndex    int64
}

type stateKind byte

const startState = stateKind(0)

const objectTagOpenState = stateKind('{')

const objectTagKeyState = stateKind('O')

const objectTagCloseState = stateKind('}')

const arrayTagOpenState = stateKind('[')

const arrayTagKeyState = stateKind('A')

const arrayTagIndexState = stateKind('I')

const arrayTagElemState = stateKind('E')

const arrayTagCloseState = stateKind(']')

const endState = stateKind('$')
