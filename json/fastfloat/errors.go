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

package fastfloat

import "errors"

var errCannotParseNumberFromEmpty = errors.New("cannot parse number from empty string")
var errCannotParseNumber = errors.New("cannot parse number")
var errCannotParseNumberUnparsedTail = errors.New("unparsed tail left after parsing number")
var errCannotParseNumberExponent = errors.New("cannot parse exponent")
var errCannotParseNumberMantissa = errors.New("cannot parse mantissa")
var errCannotParseNumberFindMantissa = errors.New("cannot find mantissa")
var errCannotParseNumberMissingIntInFrac = errors.New("missing integer and fractional part")
