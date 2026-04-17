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

package jsonschema

// JSONSchemaAble is an extra method for a Parser that distinguishes between objects and arrays.
// This allows the Parser to tagged to handle JSONSchema `{"type": "object"}` and `{"type":"array"}`.
// Tagging is actually added using the `tag` package, this interface only makes tagging possible.
type JSONSchemaAble interface {
	// JSONSchemaType returns a type that distinguishes between arrays and objects, after Next returned an EnterHint.
	JSONSchemaType() JSONSchemaType
}

type JSONSchemaType byte

const JSONSchemaTypeUnknown = JSONSchemaType(0)

const JSONSchemaTypeObject = JSONSchemaType('{')

const JSONSchemaTypeArray = JSONSchemaType('[')
