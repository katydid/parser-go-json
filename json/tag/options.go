// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tag

// Option is used set options when creating a new JSON Parser.
type Option func(*tagger)

// WithTags tags
// 1. each object with an object key, for example `{"a": null}` is parsed as `{"object": {"a": null}}`.
// 2. each array with an array key, for example `{"a": []}` is parsed as `{"a": {"array": []}}`.
func WithTags() func(*tagger) {
	return func(t *tagger) {
		t.tag = true
	}
}

// WithIndexes tags each array item with an index:
// for example `["a", "b"]` is parsed as `[0: "a", 1: "b"]`.
func WithIndexes() func(*tagger) {
	return func(t *tagger) {
		t.index = true
	}
}

// WithAllocator replaces the default `func(size int) []byte { return make([]byte, size) }` allocator
// with a different allocator function.
// Usually an allocator that uses a pool.
func WithAllocator(alloc func(int) []byte) func(*tagger) {
	return func(o *tagger) {
		o.alloc = alloc
	}
}
