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

// Option is used set options when creating a new JSON Parser.
type Option func(*tagger)

// WithObjectTag tags each object with an object key, for example `{"a": null}` is parsed as `{"object": {"a": null}}`.
func WithObjectTag() func(*tagger) {
	return func(t *tagger) {
		t.tagObjects = true
	}
}

// WithArrayTag tags each array with an array key, for example `{"a": []}` is parsed as `{"a": {"array": []}}`.
func WithArrayTag() func(*tagger) {
	return func(t *tagger) {
		t.tagArrays = true
	}
}
