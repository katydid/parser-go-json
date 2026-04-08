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

type options struct {
	tags  bool
	index bool
	alloc func(int) []byte
	buf   []byte
}

func newOptions(opts ...Option) *options {
	o := &options{
		// set default allocator.
		alloc: func(size int) []byte { return make([]byte, size) },
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Option is used set options when creating a new JSON Parser.
type Option func(*options)

// WithAllocator replaces the default `func(size int) []byte { return make([]byte, size) }` allocator
// with a different allocator function.
// Usually an allocator that uses a pool.
func WithAllocator(alloc func(int) []byte) func(*options) {
	return func(o *options) {
		o.alloc = alloc
	}
}

// WithBuffer passes in a buffer to parse.
func WithBuffer(buf []byte) func(*options) {
	return func(o *options) {
		o.buf = buf
	}
}
