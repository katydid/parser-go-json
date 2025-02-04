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

package token

import (
	"errors"
)

// ErrNotDouble is an error that represents a type error.
var ErrNotDouble = errors.New("value is not a double")

// ErrNotInt is an error that represents a type error.
var ErrNotInt = errors.New("value is not a int")

// ErrNotUint is an error that represents a type error.
var ErrNotUint = errors.New("value is not a uint")

// ErrNotBool is an error that represents a type error.
var ErrNotBool = errors.New("value is not a bool")

// ErrNotString is an error that represents a type error.
var ErrNotString = errors.New("value is not a string")

// ErrNotBytes is an error that represents a type error.
var ErrNotBytes = errors.New("value is not a bytes")

// ErrNotValue is an error that says that all attempts at getting a value failed.
var ErrNotValue = errors.New("all values return an error")

// errUnquote returns an error that resulted from trying to unquote a string.
var errUnquote = errors.New("unable to unquote string")
