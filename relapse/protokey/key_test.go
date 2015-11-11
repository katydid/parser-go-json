//  Copyright 2015 Walter Schulze
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

package protokey

import (
	"github.com/katydid/katydid/expr/ast"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/serialize/debug"
	"testing"
)

func TestKeyField(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny())
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().GetName().GetUintValue() != 1 {
		t.Fatalf("expected field 1, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyOr(t *testing.T) {
	p := relapse.NewOr(
		relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()),
		relapse.NewTreeNode(relapse.NewStringName("B"), relapse.NewZAny()),
	)
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetLeftPattern().GetTreeNode().GetName().GetName().GetUintValue() != 1 {
		t.Fatalf("expected field 1, but got %v", gkey)
	}
	if gkey.GetTopPattern().GetOr().GetRightPattern().GetTreeNode().GetName().GetName().GetUintValue() != 2 {
		t.Fatalf("expected field 2, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyTree(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewStringName("C"), relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()))
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().GetName().GetUintValue() != 3 {
		t.Fatalf("expected field 3, but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetName().GetUintValue() != 1 {
		t.Fatalf("expected field 1, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyAnyName(t *testing.T) {
	p := relapse.NewOr(
		relapse.NewTreeNode(relapse.NewNameChoice(relapse.NewAnyName(), relapse.NewStringName("C")), relapse.NewZAny()),
		relapse.NewTreeNode(relapse.NewStringName("B"), relapse.NewZAny()),
	)
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetLeftPattern().GetTreeNode().GetName().GetNameChoice().GetLeft().GetAnyName() == nil {
		t.Fatalf("expected any field, but got %v", gkey)
	}
	if gkey.GetTopPattern().GetOr().GetLeftPattern().GetTreeNode().GetName().GetNameChoice().GetRight().GetName().GetUintValue() != 3 {
		t.Fatalf("expected field 3, but got %v", gkey)
	}
	if gkey.GetTopPattern().GetOr().GetRightPattern().GetTreeNode().GetName().GetName().GetUintValue() != 2 {
		t.Fatalf("expected field 2, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyRecursive(t *testing.T) {
	p := relapse.NewOr(
		relapse.NewTreeNode(relapse.NewNameChoice(relapse.NewAnyName(), relapse.NewStringName("C")), relapse.NewReference("main")),
		relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()),
	)
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetLeftPattern().GetTreeNode().GetName().GetNameChoice().GetRight().GetName().GetUintValue() != 3 {
		t.Fatalf("expected field 3, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyLeftRecursive(t *testing.T) {
	p := relapse.NewOr(
		relapse.NewReference("a"),
		relapse.NewTreeNode(relapse.NewStringName("C"), relapse.NewReference("main")),
		relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()),
	)
	g := p.Grammar().AddRef("a", relapse.NewReference("main"))
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetRightPattern().GetOr().GetLeftPattern().GetTreeNode().GetName().GetName().GetUintValue() != 3 {
		t.Fatalf("expected field 3, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestKeyLeaf(t *testing.T) {
	p := relapse.NewLeafNode(expr.NewStringVar())
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", gkey)
}

func TestKeyAny(t *testing.T) {
	p := relapse.NewConcat(
		relapse.NewZAny(),
		relapse.NewTreeNode(relapse.NewStringName("E"),
			relapse.NewTreeNode(relapse.NewAnyName(),
				relapse.NewConcat(
					relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()),
					relapse.NewTreeNode(relapse.NewStringName("B"), relapse.NewZAny()),
				),
			),
		),
		relapse.NewZAny(),
	)
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", gkey)
}

func TestRepeatedMessageWithNoFieldsOfTypeMessage(t *testing.T) {
	p := relapse.NewConcat(
		relapse.NewZAny(),
		relapse.NewTreeNode(relapse.NewStringName("KeyValue"),
			relapse.NewTreeNode(relapse.NewAnyName(),
				relapse.NewConcat(
					relapse.NewTreeNode(relapse.NewStringName("Key"), relapse.NewZAny()),
					relapse.NewTreeNode(relapse.NewStringName("Value"), relapse.NewZAny()),
				),
			),
		),
		relapse.NewZAny(),
	)
	g := p.Grammar()
	gkey, err := KeyTheGrammar("protokey", "ProtoKey", ProtokeyDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", gkey)
}

func TestUnreachable(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewStringName("NotC"), relapse.NewTreeNode(relapse.NewStringName("C"), relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny())))
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().GetAnyNameExcept().GetExcept().AnyName == nil {
		t.Fatalf("expected field !(_), but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetAnyNameExcept().GetExcept().AnyName == nil {
		t.Fatalf("expected field !(_), but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetAnyNameExcept().GetExcept().AnyName == nil {
		t.Fatalf("expected field !(_), but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestNotUnreachable(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewAnyNameExcept(relapse.NewStringName("NotC")), relapse.NewTreeNode(relapse.NewStringName("C"), relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny())))
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().GetAnyNameExcept().GetExcept().GetAnyNameExcept().GetExcept().AnyName == nil {
		t.Fatalf("expected field !(!(_)), but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetName().GetUintValue() != 3 {
		t.Fatalf("expected field 3, but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetName().GetUintValue() != 1 {
		t.Fatalf("expected field 1, but got %v", gkey)
	}
	t.Logf("%v", gkey)
}

func TestNotUnreachableArray(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewAnyNameExcept(relapse.NewStringName("NotC")), relapse.NewTreeNode(relapse.NewStringName("F"),
		relapse.NewConcat(relapse.NewZAny(), relapse.NewTreeNode(relapse.NewAnyName(),
			relapse.NewZAny(),
		))))
	g := p.Grammar()
	gkey, err := KeyTheGrammar("debug", "Debug", debug.DebugDescription(), g)
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().GetAnyNameExcept().GetExcept().GetAnyNameExcept().GetExcept().AnyName == nil {
		t.Fatalf("expected field !(!(_)), but got %v", gkey)
	}
	if gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetName().GetName().GetUintValue() != 6 {
		t.Fatalf("expected field 6, but got %v", gkey)
	}
	if name := gkey.GetTopPattern().GetTreeNode().GetPattern().GetTreeNode().GetPattern().GetConcat().GetRightPattern().GetTreeNode().GetName(); name.AnyName == nil {
		t.Fatalf("expected field _, but got %v", name)
	}
	t.Logf("%v", gkey)
}

func TestTopsyTurvy(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewAnyName(), relapse.NewTreeNode(relapse.NewStringName("A"), relapse.NewZAny()))
	gkey, err := KeyTheGrammar("protokey", "TopsyTurvy", ProtokeyDescription(), p.Grammar())
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().AnyName != nil {
		t.Fatalf("did not expected any name, since this causes a name conflict")
	}
	//TODO more checks
	t.Logf("%v", gkey)
}

func TestKnot(t *testing.T) {
	p := relapse.NewTreeNode(relapse.NewAnyName(), relapse.NewTreeNode(relapse.NewAnyName(), relapse.NewTreeNode(relapse.NewStringName("Elbow"), relapse.NewZAny())))
	gkey, err := KeyTheGrammar("protokey", "Knot", ProtokeyDescription(), p.Grammar())
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetTreeNode().GetName().AnyName != nil {
		t.Fatalf("did not expected any name, since this causes a name conflict")
	}
	//TODO more checks
	t.Logf("%v", gkey)
}

func TestRecursiveKnotTurn(t *testing.T) {
	p := relapse.NewOr(relapse.NewTreeNode(relapse.NewAnyName(), relapse.NewReference("main")), relapse.NewTreeNode(relapse.NewStringName("Turn"), relapse.NewZAny()))
	gkey, err := KeyTheGrammar("protokey", "Knot", ProtokeyDescription(), p.Grammar())
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetRightPattern().GetTreeNode().GetName().AnyNameExcept != nil {
		t.Fatalf("did not expected not")
	}
	//TODO more checks
	t.Logf("%v", gkey)
}

func TestRecursiveKnotElbow(t *testing.T) {
	p := relapse.NewOr(relapse.NewTreeNode(relapse.NewAnyName(), relapse.NewReference("main")), relapse.NewTreeNode(relapse.NewStringName("Elbow"), relapse.NewZAny()))
	gkey, err := KeyTheGrammar("protokey", "Knot", ProtokeyDescription(), p.Grammar())
	if err != nil {
		t.Fatal(err)
	}
	if gkey.GetTopPattern().GetOr().GetRightPattern().GetTreeNode().GetName().AnyNameExcept != nil {
		t.Fatalf("did not expected not")
	}
	//TODO more checks
	t.Logf("%v", gkey)
}
