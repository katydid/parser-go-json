//  Copyright 2013 Walter Schulze
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

// downgrade package downgrades a new parse.Parser implementation to an old parser.Interface implementation.
package downgrade

import (
	"fmt"
	"io"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
	"github.com/katydid/parser-go/parser"
)

type parserWithInit interface {
	parse.Parser
	Init([]byte)
}

type interfaceWithInit interface {
	parser.Interface
	Init([]byte) error
}

type downgradeParser struct {
	action  action
	actions []action
	state   state
	stack   []state
	parser  parserWithInit
}

// Parser downgrades a new parse.Parser implementation to an old parser.Interface implementation with an Init method.
func Parser(parser parse.Parser) parser.Interface {
	parserWithInit, ok := parser.(parserWithInit)
	if !ok {
		parserWithInit = &noopInit{parser}
	}
	return &downgradeParser{
		stack:   make([]state, 0, 10),
		actions: make([]action, 0, 10),
		parser:  parserWithInit,
	}
}

type noopInit struct {
	parse.Parser
}

func (n *noopInit) Init([]byte) {}

// ParserWithInit downgrades a new parse.Parser implementation to an old parser.Interface implementation with an Init method.
func ParserWithInit(parser parserWithInit) interfaceWithInit {
	return &downgradeParser{
		stack:   make([]state, 0, 10),
		actions: make([]action, 0, 10),
		parser:  parser,
	}
}

func (p *downgradeParser) Init(buf []byte) error {
	p.parser.Init(buf)
	p.action = nextAction
	p.actions = p.actions[:0]
	p.state = atStartState
	p.stack = p.stack[:0]
	return nil
}

func (p *downgradeParser) nextAtStartState(action action) error {
	switch action {
	case nextAction:
		parseHint, err := p.parser.Next()
		if err != nil {
			return err
		}
		switch parseHint {
		case parse.EnterHint:
			p.state = atFieldState
			parseHintNext, err := p.parser.Next()
			if err != nil {
				return err
			}
			if parseHintNext == parse.LeaveHint {
				return p.eof()
			}
			return nil
		case parse.ValueHint, parse.FieldHint:
			p.state = inLeafState
			return nil
		case parse.LeaveHint:
			return errExpectedLeave
		}
		panic(fmt.Sprintf("unreachable: %v", parseHint))
	case downAction:
		return errNextShouldBeCalled
	case upAction:
		return errNextShouldBeCalled
	}
	panic("unreachable")
}

func (p *downgradeParser) nextInLeafState(action action) error {
	switch action {
	case nextAction:
		// We already parsed the leaf, so there is no next element.
		return p.eof()
	case downAction:
		// Cannot call Down when in leaf, since we are the bottom.
		return errDownLeaf
	case upAction:
		// We can go up, if we are an array element or value for a key in an object.
		if err := p.pop(); err != nil {
			return err
		}
		// If we were in an object, then move onto next key.
		// If we were in an array, them ove onto the next element.
		return p.next()
	}
	panic("unreachable")
}

func (p *downgradeParser) nextAtEOF(action action) error {
	switch action {
	case nextAction:
		// If Next is called too many times, just keep on return EOF
		return p.eof()
	case downAction:
		// We cannot go down if we are at the EOF
		return errDownEOF
	case upAction:
		// We can go up, if we are an array element or value for a key in an object.
		if err := p.pop(); err != nil {
			return err
		}
		// If we were in an object, then move onto next key.
		// If we were in an array, them ove onto the next element.
		return p.next()
	}
	panic("unreachable")
}

func (p *downgradeParser) nextInAtFieldState(action action) error {
	// inAtFieldStateKind represents that we have scanned a key
	switch action {
	case nextAction:
		// We want to skip over value and move onto next key.
		// We start by skipping over the value.
		if err := p.parser.Skip(); err != nil {
			return err
		}
		// Next we move onto the Next key.
		parseKind, err := p.parser.Next()
		if parseKind == parse.LeaveHint {
			// If the Object has ended, we return eof
			return p.eof()
		}
		return err
	case downAction:
		// Set the state to be ready to parse to next key, when Up is called.
		p.state = atValueState
		// We do not want to skip over the value, we want to continue into value.
		// We start by pushing the current state to the stack.
		if err := p.push(); err != nil {
			return err
		}
		// The state is reset to be the start state.
		// We can call this' Next method, instead of the parser's Next method.
		return p.next()
	case upAction:
		// We want to skip over value and the rest of the object and move onto the next key.
		// We start by skipping over the the value.
		if err := p.parser.Skip(); err != nil {
			return err
		}
		// Next we skip over the rest of the object.
		if err := p.parser.Skip(); err != nil {
			return err
		}
		// Now we pop the stack
		if err := p.pop(); err != nil {
			return err
		}
		// Finally we move onto the next key or element, if we were in an array.
		return p.next()
	}
	panic("unreachable")
}

func (p *downgradeParser) nextInAtValueState(action action) error {
	// inAtValueStateKind represents that we have scanned a value and Up was called.
	switch action {
	case nextAction:
		// Up was just called and we need to scan to the Next key.
		parseKind, err := p.parser.Next()
		if parseKind == parse.LeaveHint {
			// If the Object has ended, we return eof
			return p.eof()
		}
		// Set the state to the next key
		p.state = atFieldState
		return err
	case downAction:
		// We can't call Down right, while at the end of value.
		return errDown
	case upAction:
		// We want to skip over the rest of the object and move onto the next key.
		// We skip over the rest of the object.
		if err := p.parser.Skip(); err != nil {
			return err
		}
		// Now we pop the stack
		if err := p.pop(); err != nil {
			return err
		}
		// Finally we move onto the next key or element, if we were in an array.
		return p.next()
	}
	panic("unreachable")
}

func (p *downgradeParser) eof() error {
	if len(p.stack) == 0 {
		// if we are at the top of stack, then check that there is no more input left.
		_, err := p.parser.Next()
		if err == nil {
			return errExpectedEOF
		}
		if err != io.EOF {
			return err
		}
	}
	// When EOF is returned also set the state to an EOF state.
	// This state allows us to call Up.
	p.state = atEOFState
	return io.EOF
}

func (p *downgradeParser) next() error {
	action := p.action
	// do not forget to reset action
	p.action = nextAction
	switch p.state {
	case atStartState:
		return p.nextAtStartState(action)
	case inLeafState:
		return p.nextInLeafState(action)
	case atFieldState:
		return p.nextInAtFieldState(action)
	case atValueState:
		return p.nextInAtValueState(action)
	case atEOFState:
		return p.nextAtEOF(action)
	}
	panic("unreachable")
}

func (p *downgradeParser) nexts() error {
	lastAction := p.action
	for i := 0; i < len(p.actions); i++ {
		p.action = p.actions[i]
		if err := p.next(); err != nil {
			// ignore EOF if we still have more actions to perform.
			if err != io.EOF || i == len(p.actions) {
				return err
			}
		}
	}
	p.actions = p.actions[:0]
	p.action = lastAction
	return p.next()
}

func (p *downgradeParser) Next() error {
	return p.nexts()
}

func (p *downgradeParser) Down() {
	p.pushAction(downAction)
}

func (p *downgradeParser) Up() {
	p.pushAction(upAction)
}

func (p *downgradeParser) pushAction(newAction action) {
	if p.action == nextAction {
		p.action = newAction
		return
	}
	// when Up is called straight after Down, we simply call next.
	if p.action == downAction && newAction == upAction {
		p.action = nextAction
		return
	}
	p.actions = append(p.actions, p.action)
	p.action = newAction
}

func (p *downgradeParser) push() error {
	// Append the current state to the stack.
	p.stack = append(p.stack, p.state)
	p.state = atStartState
	return nil
}

func (p *downgradeParser) pop() error {
	if len(p.stack) == 0 {
		return errPop
	}
	top := len(p.stack) - 1
	// Set the current state to the state on top of the stack.
	p.state = p.stack[top]
	// Remove the state on the top the stack from the stack,
	// but do it in a way that keeps the capacity,
	// so we can reuse it the next time Down is called.
	p.stack = p.stack[:top]
	return nil
}

func (p *downgradeParser) IsLeaf() bool {
	return p.state == inLeafState
}

func (p *downgradeParser) Bool() (bool, error) {
	tokenKind, _, err := p.parser.Token()
	if err != nil {
		return false, err
	}
	if tokenKind == parse.FalseKind {
		return false, nil
	}
	if tokenKind == parse.TrueKind {
		return true, nil
	}
	return false, errNotBool
}

func (p *downgradeParser) Int() (int64, error) {
	tokenKind, bs, err := p.parser.Token()
	if err != nil {
		return 0, err
	}
	if tokenKind != parse.Int64Kind {
		return 0, errNotInt
	}
	return cast.ToInt64(bs), nil
}

func (p *downgradeParser) Uint() (uint64, error) {
	i, err := p.Int()
	if err != nil {
		return 0, err
	}
	if i >= 0 {
		return uint64(i), nil
	}
	return 0, errNotUint
}

func (p *downgradeParser) Double() (float64, error) {
	tokenKind, bs, err := p.parser.Token()
	if err != nil {
		return 0, err
	}
	if tokenKind != parse.Float64Kind {
		return 0, errNotFloat
	}
	return cast.ToFloat64(bs), nil
}

func (p *downgradeParser) String() (string, error) {
	tokenKind, bs, err := p.parser.Token()
	if err != nil {
		return "", err
	}
	if tokenKind != parse.StringKind && tokenKind != parse.DecimalKind {
		return "", errNotString
	}
	return cast.ToString(bs), nil
}

var nullBytes = []byte{'n', 'u', 'l', 'l'}

func (p *downgradeParser) Bytes() ([]byte, error) {
	tokenKind, bs, err := p.parser.Token()
	if err != nil {
		return nil, err
	}
	if tokenKind == parse.NullKind {
		return nullBytes, nil
	}
	return bs, nil
}
