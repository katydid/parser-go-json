//  Copyright 2016 Walter Schulze
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

package convert

import (
	"fmt"
	"github.com/katydid/katydid/expr/ast"
	"github.com/katydid/katydid/expr/compose"
	"github.com/katydid/katydid/funcs"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/interp"
	"github.com/katydid/katydid/relapse/nameexpr"
	"sort"
	"strconv"
	"strings"
)

type auto struct {
	start    int
	states   map[int]*state
	patterns []*relapse.Pattern
}

func (this *auto) draw() {
	s := []string{}
	a := func(line string) {
		s = append(s, line)
	}
	a("digraph {")
	a("\trankdir=LR;")
	a("\tsize=\"8,5\"")
	accepts := []string{}
	other := []string{}
	for _, s := range this.states {
		current := this.patterns[s.current].String()
		scurrent := strconv.Quote(current)
		if s.final {
			accepts = append(accepts, scurrent)
		} else {
			other = append(other, scurrent)
		}
	}
	sort.Strings(accepts)
	sort.Strings(other)
	a("\tnode [shape = doublecircle]; " + strings.Join(accepts, " "))
	a("\tnode [shape = circle]; " + strings.Join(other, " "))
	for si, _ := range this.patterns {
		s := this.states[si]
		current := this.patterns[s.current].String()
		scurrent := strconv.Quote(current)
		if this.start == s.current {
			a("\tstart -> " + scurrent)
		}
		for _, t := range s.trans {
			v := funcs.Sprint(t.value)
			a("\t" + scurrent + " -> " + strconv.Quote(this.patterns[t.down].String()) + " [ label = " + strconv.Quote(v+"&uarr;"+current) + " ];")
			for _, u := range t.ups {
				a("\t" + strconv.Quote(this.patterns[u.ret].String()) + " -> " + strconv.Quote(this.patterns[u.dst].String()) + " [ label = " + strconv.Quote(v+"&darr;"+current) + " ];")
			}
		}
	}
	a("}")
	fmt.Printf("%s\n", strings.Join(s, "\n"))
}

func Convert(refs relapse.RefLookup, p *relapse.Pattern) *auto {
	c := newConverter(refs)
	start := c.addPattern(p)
	c.convert(p)
	finals := c.getReachable(c.states[start])
	c.exhaust()
	for _, f := range finals {
		if interp.Nullable(c.refs, c.getPattern(c.states[f].current)) {
			c.states[f].final = true
		}
	}
	a := &auto{start, c.states, c.patterns}
	a.draw()
	return a
}

type converter struct {
	states   map[int]*state
	patterns []*relapse.Pattern
	refs     relapse.RefLookup
}

func newConverter(refs relapse.RefLookup) *converter {
	return &converter{
		states:   make(map[int]*state),
		patterns: []*relapse.Pattern{relapse.NewZAny(), relapse.NewNot(relapse.NewZAny()), relapse.NewEmpty()},
		refs:     refs,
	}
}

func (this *converter) addPattern(p *relapse.Pattern) int {
	p = interp.Simplify(this.refs, p)
	for i, pat := range this.patterns {
		if pat.Equal(p) {
			return i
		}
	}
	this.patterns = append(this.patterns, p)
	return len(this.patterns) - 1
}

func (this *converter) getPattern(i int) *relapse.Pattern {
	return this.patterns[i]
}

func (this *converter) getAny() int {
	return 0
}

func (this *converter) getNone() int {
	return 1
}

func (this *converter) getEmpty() int {
	return 2
}

func (this *converter) exhaust() {
	changed := true
	for changed {
		changed = false
		for i, p := range this.patterns {
			if _, ok := this.states[i]; !ok {
				fmt.Printf("exhausting %v %d\n", p, len(this.states))
				this.convert(p)
				changed = true
			}
		}
	}
}

func setToList(ms map[int]struct{}) []int {
	is := make([]int, 0, len(ms))
	for i, _ := range ms {
		is = append(is, i)
	}
	sort.Ints(is)
	return is
}

func (this *converter) getReachable(state *state) []int {
	allDsts := this.states[state.current].dsts()
	allDsts[state.current] = struct{}{}
	prevDsts := allDsts
	nextDsts := make(map[int]struct{})
	for {
		for _, d := range setToList(prevDsts) {
			if _, ok := this.states[d]; !ok {
				this.convert(this.getPattern(d))
			}
			if this.states[d] == nil {
				fmt.Printf("yeah: %v\n", this.getPattern(d))
				panic("yeah")
			}
			for n, _ := range this.states[d].dsts() {
				if _, ok := allDsts[n]; !ok {
					nextDsts[n] = struct{}{}
					allDsts[n] = struct{}{}
				}
			}
		}
		if len(nextDsts) == 0 {
			break
		}
		prevDsts = make(map[int]struct{})
		for n, _ := range nextDsts {
			prevDsts[n] = struct{}{}
		}
		nextDsts = make(map[int]struct{})
	}
	return setToList(allDsts)
}

type state struct {
	current int
	final   bool
	trans   []tran
}

type tran struct {
	value funcs.Bool
	down  int
	ups   []up
}

func (this tran) Equal(that tran) bool {
	if this.down != that.down {
		return false
	}
	if len(this.ups) != len(that.ups) {
		return false
	}
	for i := range this.ups {
		if this.ups[i].ret != that.ups[i].ret {
			return false
		}
		if this.ups[i].dst != that.ups[i].dst {
			return false
		}
	}
	return true
}

type up struct {
	ret int
	dst int
}

func (this *state) dsts() map[int]struct{} {
	dsts := make(map[int]struct{})
	for _, t := range this.trans {
		for _, u := range t.ups {
			dsts[u.dst] = struct{}{}
		}
	}
	return dsts
}

func (this *converter) toStr(s *state) string {
	current := this.getPattern(s.current)
	ts := []string{}
	for _, t := range s.trans {
		us := []string{}
		for _, u := range t.ups {
			us = append(us, "( "+this.getPattern(u.ret).String()+" ^ "+this.getPattern(u.dst).String()+" )")
		}
		ts = append(ts, funcs.Sprint(t.value)+" -> "+this.getPattern(t.down).String()+" [ "+strings.Join(us, ", ")+" ]")
	}
	return current.String() + " { " + strings.Join(ts, ", ") + " }"
}

func (this *converter) newState(current int, trans []tran) *state {
	fmt.Printf("deduping %v\n", this.toStr(&state{current, false, trans}))
	for i := range trans {
		trans[i].ups = this.dedup2(trans[i].ups)
	}
	mtrans := make(map[string]int, len(trans))
	is := make([]int, 0, len(trans))
	for i := range trans {
		sf := funcs.Sprint(trans[i].value)
		if sf == "false" {
			continue
		}
		if j, ok := mtrans[sf]; ok {
			if !trans[i].Equal(trans[j]) {
				fmt.Printf("deduped! %v\n", this.toStr(&state{current, false, trans}))
				panic("wtf")
			}
		}
		mtrans[funcs.Sprint(trans[i].value)] = i
		is = append(is, i)
	}
	trans2 := make([]tran, 0, len(is))
	for _, i := range is {
		trans2 = append(trans2, trans[i])
	}
	s := &state{
		current: current,
		trans:   trans2,
	}
	fmt.Printf("deduped  %v\n", this.toStr(s))
	return s
}

func (this *converter) dedup2(ups []up) []up {
	rups := make(map[int]int)
	for _, up := range ups {
		if d, ok := rups[up.ret]; !ok {
			rups[up.ret] = up.dst
		} else if d != up.dst {
			fmt.Printf("wtf\n")
			rups[up.ret] = this.addPattern(relapse.NewOr(this.getPattern(up.dst), this.getPattern(d)))
		}
	}
	iups := []int{}
	for i := range rups {
		iups = append(iups, i)
	}
	sort.Ints(iups)
	nups := []up{}
	for _, i := range iups {
		nups = append(nups, up{i, rups[i]})
	}
	return nups
}

func (this *converter) dedup(ups []up) []up {
	rups := make(map[int]map[int]struct{})
	nups := []up{}
	for i, up := range ups {
		if _, ok := rups[up.ret]; !ok {
			nups = append(nups, ups[i])
			rups[up.ret] = make(map[int]struct{})
			rups[up.ret][up.dst] = struct{}{}
		} else if _, ok1 := rups[up.ret][up.dst]; !ok1 {
			nups = append(nups, ups[i])
			rups[up.ret][up.dst] = struct{}{}
			fmt.Printf("wtf %v -> %v \n", this.patterns[up.ret], this.patterns[up.dst])
			//panic("wtf different returns")
		}
	}
	return ups
}

func (this *converter) union(left, right *state) []tran {
	ts := []tran{}
	for _, lt := range left.trans {
		for _, rt := range right.trans {
			pdl := this.getPattern(lt.down)
			pdr := this.getPattern(rt.down)
			ups := []up{}
			for _, lu := range lt.ups {
				for _, ru := range rt.ups {
					ups = append(ups, up{
						ret: this.addPattern(relapse.NewOr(this.getPattern(lu.ret), this.getPattern(ru.ret))),
						dst: this.addPattern(relapse.NewOr(this.getPattern(lu.dst), this.getPattern(ru.dst))),
					})
				}
			}
			tt := tran{
				value: funcs.Simplify(funcs.And(lt.value, rt.value)),
				down:  this.addPattern(relapse.NewOr(pdl, pdr)),
				ups:   ups,
			}
			ts = append(ts, tt)
		}
	}
	return ts
}

func (this *converter) copy(trans []tran) []tran {
	ts := []tran{}
	for _, t := range trans {
		ups := []up{}
		for _, u := range t.ups {
			ups = append(ups, up{
				ret: u.ret,
				dst: u.dst,
			})
		}
		tt := tran{
			value: t.value,
			down:  t.down,
			ups:   ups,
		}
		ts = append(ts, tt)
	}
	return ts
}

func (this *converter) leftConcat(current int, left *state, rightCurrent int) []tran {
	ts := this.copy(left.trans)
	for i, lt := range left.trans {
		for j, lu := range lt.ups {
			ts[i].ups[j].dst = this.addPattern(relapse.NewConcat(this.getPattern(lu.dst), this.getPattern(rightCurrent)))
		}
	}
	return ts
}

func (this *converter) concat(current int, left, right *state) []tran {
	ts := this.leftConcat(current, left, right.current)
	if !interp.Nullable(this.refs, this.getPattern(left.current)) {
		return ts
	}
	newleft := &state{current, false, ts}
	fmt.Printf("concating %v\n", this.getPattern(current))
	fmt.Printf("newleft %v\n", this.toStr(newleft))
	fmt.Printf("right %v\n", this.toStr(right))
	trans := this.union(newleft, right)
	return trans
}

func (this *converter) newDeadEnd(f funcs.Bool) tran {
	return tran{f, this.getNone(), []up{{this.getNone(), this.getNone()}}}
}

func (this *converter) newAnyEnd(f funcs.Bool) tran {
	return tran{f, this.getAny(), []up{{this.getAny(), this.getAny()}}}
}

func nameToFunc(name *relapse.NameExpr) funcs.Bool {
	return funcs.Simplify(nameexpr.NameToFunc(name))
}

func not(f funcs.Bool) funcs.Bool {
	return funcs.Simplify(funcs.Not(f))
}

func exprToFunc(expr *expr.Expr) funcs.Bool {
	f, err := compose.NewBool(expr)
	if err != nil {
		panic(err)
	}
	return funcs.Simplify(f)
}

func (this *converter) convert(p *relapse.Pattern) *state {
	p = interp.Simplify(this.refs, p)
	c := this.addPattern(p)
	if this.states[c] != nil {
		return this.states[c]
	}
	this.states[c] = nil
	typ := p.GetValue()
	switch v := typ.(type) {
	case *relapse.Empty:
		s := this.newState(
			this.getEmpty(),
			[]tran{
				this.newDeadEnd(funcs.BoolConst(true)),
			},
		)
		this.states[s.current] = s
		return s
	case *relapse.TreeNode:
		current := this.addPattern(p)
		f := nameToFunc(v.GetName())
		below := this.convert(v.GetPattern())
		dsts := this.getReachable(below)
		ups := make([]up, 0, len(dsts))
		for _, d := range dsts {
			if interp.Nullable(this.refs, this.getPattern(d)) {
				ups = append(ups, up{d, this.getEmpty()})
			} else {
				ups = append(ups, up{d, this.getNone()})
			}
		}
		s := this.newState(
			current,
			[]tran{
				{f, this.addPattern(v.Pattern), ups},
				this.newDeadEnd(not(f)),
			},
		)
		this.states[s.current] = s
		return s
	case *relapse.LeafNode:
		f := exprToFunc(v.GetExpr())
		s := this.newState(
			this.addPattern(p),
			[]tran{
				{
					f, this.getEmpty(), []up{
						{this.getEmpty(), this.getEmpty()},
						{this.getNone(), this.getNone()},
					},
				},
				this.newDeadEnd(not(f)),
			},
		)
		this.states[s.current] = s
		return s
	case *relapse.Concat:
		left := this.convert(v.GetLeftPattern())
		right := this.convert(v.GetRightPattern())
		current := this.addPattern(p)
		trans := this.concat(current, left, right)
		s := this.newState(current, trans)
		this.states[s.current] = s
		return s
	case *relapse.Or:
		left := this.convert(v.GetLeftPattern())
		right := this.convert(v.GetRightPattern())
		current := this.addPattern(p)
		ts := this.union(left, right)
		s := this.newState(current, ts)
		this.states[s.current] = s
		return s
	case *relapse.And:

	case *relapse.ZeroOrMore:
		elem := this.convert(v.GetPattern())
		current := this.addPattern(p)
		ts := this.leftConcat(current, elem, current)
		// if interp.Nullable(this.refs, v.GetPattern()) {
		// 	panic("not implemented")
		// } else {
		s := this.newState(current, ts)
		this.states[s.current] = s
		return s
		//}
	case *relapse.Reference:

	case *relapse.Not:
		n := this.convert(v.GetPattern())
		trans := make([]tran, len(n.trans))
		for i := range n.trans {
			ups := []up{}
			for j := range n.trans[i].ups {
				d := n.trans[i].ups[j].dst
				dpat := this.getPattern(d)
				ndpat := relapse.NewNot(dpat)
				ups = append(ups, up{
					dst: this.addPattern(ndpat),
					ret: n.trans[i].ups[j].ret,
				})
			}
			trans[i] = tran{
				value: n.trans[i].value,
				down:  n.trans[i].down,
				ups:   ups,
			}
		}
		s := this.newState(this.addPattern(p), trans)
		this.states[s.current] = s
		return s
	case *relapse.ZAny:
		s := this.newState(
			this.getAny(),
			[]tran{
				this.newAnyEnd(funcs.BoolConst(true)),
			},
		)
		this.states[s.current] = s
		return s
	case *relapse.Contains:

	case *relapse.Optional:

	case *relapse.Interleave:

	}
	panic(fmt.Sprintf("unknown pattern typ %T %v", typ, p))
}
