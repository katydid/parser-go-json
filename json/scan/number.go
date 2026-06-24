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

package scan

import "github.com/katydid/parser-go-json/json/internal/fork/strconv"

// Number returns the offset after the prefix of a valid number.
// The number BNF:
// number := integer fraction exponent
// integer := digit | onenine digits | '-' digit | '-' onenine digits
// digits := digit | digit digits
// digit := '0' | onenine
// onenine := '1' . '9'
// fraction := "" | '.' digits
// exponent := "" | 'E' sign digits | 'e' sign digits
// sign := "" | '+' | '-'
func Number(buf []byte) (int, error) {
	state := StateStart // start
	offset := 0
	for offset < len(buf) {
		c := buf[offset]
		state = machine[state][c].state
		if state == StateError || state == StateSuccess {
			break
		}
		offset += 1
	}
	if isFailState[state] {
		return 0, errScanNumber
	}
	return offset, nil
}

var isFailState = [256]bool{'s': true, '-': true, '.': true, 'e': true, 'f': true, StateError: true}

type dst struct {
	state  byte
	action byte
}

const StateSuccess = byte(0)
const StateError = byte(255)

const StateStart = byte('s')
const StateIntegerComplete = byte('0')
const StateIntegerContinue = byte('1')
const StateIntegerNegative = byte('-')
const StateFractionStarted = byte('.')
const StateFractionOngoing = byte('/')
const StateExponentStarted = byte('e')
const StateExponentWithSign = byte('f')
const StateExponentOngoing = byte('g')

const ActionNothing = 0
const ActionIntMatinsa = 1      // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
const ActionFracMatinsa = 2     // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
const ActionSetSignNegative = 3 // sign = -1
const ActionSetExponentSign = 4 // expsign = -1
const ActionExpPart = 5         // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++

var machine = [256][256]dst{
	// start
	StateStart: { // default is error, see init below
		'0': {StateIntegerComplete, ActionIntMatinsa},      // state = '0'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'1': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'2': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'3': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'4': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'5': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'6': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'7': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'8': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'9': {StateIntegerContinue, ActionIntMatinsa},      // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'-': {StateIntegerNegative, ActionSetSignNegative}, // state = '-'; sign = -1
	},
	StateIntegerNegative: { // default is error, see init below
		'0': {StateIntegerComplete, ActionIntMatinsa}, // state = '0'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'1': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'2': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'3': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'4': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'5': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'6': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'7': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'8': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'9': {StateIntegerContinue, ActionIntMatinsa}, // state = '1'; if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
	},
	StateIntegerComplete: { // default is success
		'.': {StateFractionStarted, ActionNothing}, // state = '.'
		'e': {StateExponentStarted, ActionNothing}, // state = 'e'
		'E': {StateExponentStarted, ActionNothing}, // state = 'e'
		'0': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'1': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'2': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'3': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'4': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'5': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'6': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'7': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'8': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
		'9': {StateError, ActionNothing},           // return // ERROR: the integer was complete, for example 0, the number does not continue
	},
	StateIntegerContinue: { // default is success
		'0': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'1': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'2': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'3': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'4': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'5': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'6': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'7': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'8': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'9': {StateIntegerContinue, ActionIntMatinsa}, // if intdigits < maxMantDigits  { mantissa = (mantissa * 10) + uint64(c-'0') }; intdigits++
		'.': {StateFractionStarted, ActionNothing},    // state = '.'
		'e': {StateExponentStarted, ActionNothing},    // state = 'e'
		'E': {StateExponentStarted, ActionNothing},    // state = 'e'
	},
	StateFractionStarted: { // default is error, see init below
		'0': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'1': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'2': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'3': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'4': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'5': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'6': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'7': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'8': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'9': {StateFractionOngoing, ActionFracMatinsa}, // state = '/'; if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'e': {StateExponentStarted, ActionNothing},     // state = 'e'
		'E': {StateExponentStarted, ActionNothing},     // state = 'e'
	},
	StateFractionOngoing: { // default is success
		'0': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'1': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'2': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'3': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'4': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'5': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'6': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'7': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'8': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'9': {StateFractionOngoing, ActionFracMatinsa}, // if fracdigits+intdigits < maxMantDigits { mantissa = (mantissa * 10) + uint64(c-'0') }; fracdigits++
		'e': {StateExponentStarted, ActionNothing},     // state = 'e'
		'E': {StateExponentStarted, ActionNothing},     // state = 'e'
		'.': {StateError, ActionNothing},               // return // ERROR: fraction does not contain a .
	},
	StateExponentStarted: { // default is error, see init below
		'-': {StateExponentWithSign, ActionSetExponentSign}, // state = 'f'; expsign = -1
		'+': {StateExponentWithSign, ActionNothing},         // state = 'f'
		'0': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'1': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'2': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'3': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'4': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'5': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'6': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'7': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'8': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'9': {StateExponentOngoing, ActionExpPart},          // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
	},
	StateExponentWithSign: { // default is error, see init below
		'0': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'1': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'2': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'3': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'4': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'5': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'6': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'7': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'8': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'9': {StateExponentOngoing, ActionExpPart}, // state = 'g'; if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
	},
	StateExponentOngoing: { // default is success
		'0': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'1': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'2': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'3': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'4': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'5': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'6': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'7': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'8': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
		'9': {StateExponentOngoing, ActionExpPart}, // if exppart < 10000 { exppart = (exppart * 10) + int(c-'0') }; expdigits++
	},
}

func init() {
	// set the default states to error for some states.
	zerodst := dst{state: 0, action: 0}
	defaultIsError := []byte{StateStart, StateIntegerNegative, StateFractionStarted, StateExponentStarted, StateExponentWithSign}
	for _, s := range defaultIsError {
		for c, v := range machine[s] {
			if v == zerodst {
				machine[s][c] = dst{state: StateError}
			}
		}
	}
}

// Number returns the offset after the prefix of a valid number.
// The number BNF:
// number := integer fraction exponent
// integer := digit | onenine digits | '-' digit | '-' onenine digits
// digits := digit | digit digits
// digit := '0' | onenine
// onenine := '1' . '9'
// fraction := "" | '.' digits
// exponent := "" | 'E' sign digits | 'e' sign digits
// sign := "" | '+' | '-'
func ParseNumber(buf []byte) (offset int, intres int64, intok bool, floatres float64, floatok bool, decimalok bool) {
	sign := int64(1)
	expsign := 1
	var mantissa uint64
	var intdigits int
	var fracdigits int
	var exppart int
	var expdigits int
	maxMantDigits := 19 // 10^19 fits in uint64
	state := StateStart // start
	for _, c := range buf {
		dst := machine[state][c]
		state = dst.state
		switch dst.action {
		case ActionNothing:
		case ActionIntMatinsa:
			if intdigits < maxMantDigits {
				mantissa = (mantissa * 10) + uint64(c-'0')
			}
			intdigits++
		case ActionFracMatinsa:
			if intdigits+fracdigits < maxMantDigits {
				mantissa = (mantissa * 10) + uint64(c-'0')
			}
			fracdigits++
		case ActionSetSignNegative:
			sign = -1
		case ActionSetExponentSign:
			expsign = -1
		case ActionExpPart:
			if exppart < 10000 {
				exppart = (exppart * 10) + int(c-'0')
			}
			expdigits++
		}
		if state == StateError || state == StateSuccess {
			break
		}
		offset++
	}
	if isFailState[state] {
		return // ERROR: these are not accepting states
	}

	neg := sign == -1
	// It is an int
	if expdigits == 0 && fracdigits == 0 {
		if intdigits > maxMantDigits {
			// It uses more digits than MaxInt64 and MinInt64, so it is decimal
			decimalok = true
			return
		}
		cutOffInt64 := uint64(1 << uint(64-1))
		if (!neg && mantissa < cutOffInt64) || (neg && mantissa <= cutOffInt64) {
			intres = sign * int64(mantissa)
			intok = true
			return
		} else {
			// It is larger than an int, so it must be decimal
			decimalok = true
			return
		}
	}

	// Check if float needs to be truncated
	trunc := false
	ndMant := intdigits + fracdigits
	if ndMant >= maxMantDigits {
		trunc = true
		ndMant = maxMantDigits
	}

	// calculate exp
	exp := 0
	dp := intdigits
	if expdigits > 0 {
		dp += exppart * expsign
	}
	if mantissa != 0 {
		exp = dp - ndMant
	}

	f, ok := strconv.TryParseFloat(buf[:offset], mantissa, exp, neg, trunc)
	if ok {
		floatres = f
		floatok = true
		return
	}

	// It is larger than an int, so it must be decimal
	decimalok = true
	return
}
