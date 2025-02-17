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

var errDown = errors.New("cannot go Down")

var errPop = errors.New("stack is length zero, cannot go Up")

var errNextShouldBeCalled = errors.New("Next should also be called at the start of parsing, after Down and after Up")

var errDownLeaf = errors.New("cannot call Down in Leaf")

var errDownEOF = errors.New("cannot call Down at EOF")

var errExpectedEOF = errors.New("expected EOF")
