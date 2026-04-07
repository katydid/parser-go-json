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

import (
	"fmt"

	jsonparse "github.com/katydid/parser-go-json/json/parse"
)

type state struct {
	kind       stateKind
	hint       jsonparse.Hint
	arrayIndex int64
}

func (s state) String() string {
	return fmt.Sprintf("%c %v %d", s.kind, s.hint, s.arrayIndex)
}

type stateKind byte

const startState = stateKind(0)

const objectTagOpenState = stateKind('{')

const objectTagKeyOpenState = stateKind('O')

const objectTagKeyCloseState = stateKind('C')

const objectTagCloseState = stateKind('}')

const arrayTagOpenState = stateKind('[')

const arrayTagKeyOpenState = stateKind('A')

const arrayTagKeyCloseState = stateKind('Z')

const arrayTagIndexState = stateKind('I')

const arrayTagElemState = stateKind('E')

const endState = stateKind('$')
