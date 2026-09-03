package parse

import (
	"testing"

	"katydid.org.za/go/parser-go/expect"
	"katydid.org.za/go/parser-go/parse"
)

// Test adapted from parser-go-xml

func TestSkipValue(t *testing.T) {
	// elemStr := `<A><B>C</B></A>`
	str := `{"A": {"B": "C"}}`
	x := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipSingleValue(t *testing.T) {
	// elemStr := `<A>B</A>`
	elemStr := `{"A": "B"}`
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipElement(t *testing.T) {
	// elemStr := `<A><B><C>D</C></B><E/></A>`
	elemStr := `{"A": {"B": {"C": "D"}, "E": {}}}`
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "E")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipElementAfterEnter(t *testing.T) {
	// elemStr := `<A><B><C>D</C></B><E/></A>`
	elemStr := `{"A": {"B": {"C": "D"}, "E": {}}}`
	x := NewParser(WithBuffer([]byte(elemStr)))
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.Hint(t, x, parse.EnterHint)
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "E")
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.LeaveHint)

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipAttributeValue(t *testing.T) {
	// str := `<A B="C" D="E"/>`
	str := `{"A": {"B": "C", "D": "E"}}`
	x := NewParser(WithBuffer([]byte(str)))

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "D")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "E")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemStart(t *testing.T) {
	// str := `<A><B>C</B></A>`
	str := `{"A": {"B": "C"}}`
	x := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemMid(t *testing.T) {
	// str := `<A><B>C</B><D>E</D><F>G</F></A>`
	str := `{"A": {"B": "C", "D": "E", "F": "G"}}`
	x := NewParser(WithBuffer([]byte(str)))
	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestSkipInElemAttr(t *testing.T) {
	// str := `<A b="c" d="e"><B>C</B><D>E</D></A>`
	str := `{"A": {"b": "c", "d": "e", "B": "C", "D": "E"}}`
	x := NewParser(WithBuffer([]byte(str)))

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "b")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "c")

	expect.NoErr(t, x.Skip)

	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}

func TestParseInElemAttr(t *testing.T) {
	// str := `<A b="c" d="e"><B>C</B><D>E</D></A>`
	str := `{"A": {"b": "c", "d": "e", "B": "C", "D": "E"}}`
	x := NewParser(WithBuffer([]byte(str)))

	expect.Hint(t, x, parse.EnterHint)
	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "A")
	expect.Hint(t, x, parse.EnterHint)

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "b")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "c")

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "d")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "e")

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "B")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "C")

	expect.Hint(t, x, parse.FieldHint)
	expect.String(t, x, "D")
	expect.Hint(t, x, parse.ValueHint)
	expect.String(t, x, "E")

	expect.Hint(t, x, parse.LeaveHint)
	expect.Hint(t, x, parse.LeaveHint)
	expect.EOF(t, x)
}
