// Code generated by goderive DO NOT EDIT.

package ast

import (
	types "github.com/katydid/katydid/relapse/types"
)

func deriveCopyToGrammar(this, that *Grammar) {
	if this.TopPattern == nil {
		that.TopPattern = nil
	} else {
		that.TopPattern = new(Pattern)
		deriveCopyToPattern(this.TopPattern, that.TopPattern)
	}
	if this.PatternDecls == nil {
		that.PatternDecls = nil
	} else {
		if that.PatternDecls != nil {
			if len(this.PatternDecls) > len(that.PatternDecls) {
				if cap(that.PatternDecls) >= len(this.PatternDecls) {
					that.PatternDecls = (that.PatternDecls)[:len(this.PatternDecls)]
				} else {
					that.PatternDecls = make([]*PatternDecl, len(this.PatternDecls))
				}
			} else if len(this.PatternDecls) < len(that.PatternDecls) {
				that.PatternDecls = (that.PatternDecls)[:len(this.PatternDecls)]
			}
		} else {
			that.PatternDecls = make([]*PatternDecl, len(this.PatternDecls))
		}
		deriveCopyToSliceOfPtrToPatternDecl(this.PatternDecls, that.PatternDecls)
	}
	if this.After == nil {
		that.After = nil
	} else {
		that.After = new(Space)
		deriveCopyToPtrToSpace(this.After, that.After)
	}
}

func deriveCopyToPattern(this, that *Pattern) {
	if this.Empty == nil {
		that.Empty = nil
	} else {
		that.Empty = new(Empty)
		deriveCopyToPtrToEmpty(this.Empty, that.Empty)
	}
	if this.TreeNode == nil {
		that.TreeNode = nil
	} else {
		that.TreeNode = new(TreeNode)
		deriveCopyToPtrToTreeNode(this.TreeNode, that.TreeNode)
	}
	if this.LeafNode == nil {
		that.LeafNode = nil
	} else {
		that.LeafNode = new(LeafNode)
		deriveCopyToPtrToLeafNode(this.LeafNode, that.LeafNode)
	}
	if this.Concat == nil {
		that.Concat = nil
	} else {
		that.Concat = new(Concat)
		deriveCopyToPtrToConcat(this.Concat, that.Concat)
	}
	if this.Or == nil {
		that.Or = nil
	} else {
		that.Or = new(Or)
		deriveCopyToPtrToOr(this.Or, that.Or)
	}
	if this.And == nil {
		that.And = nil
	} else {
		that.And = new(And)
		deriveCopyToPtrToAnd(this.And, that.And)
	}
	if this.ZeroOrMore == nil {
		that.ZeroOrMore = nil
	} else {
		that.ZeroOrMore = new(ZeroOrMore)
		deriveCopyToPtrToZeroOrMore(this.ZeroOrMore, that.ZeroOrMore)
	}
	if this.Reference == nil {
		that.Reference = nil
	} else {
		that.Reference = new(Reference)
		deriveCopyToPtrToReference(this.Reference, that.Reference)
	}
	if this.Not == nil {
		that.Not = nil
	} else {
		that.Not = new(Not)
		deriveCopyToPtrToNot(this.Not, that.Not)
	}
	if this.ZAny == nil {
		that.ZAny = nil
	} else {
		that.ZAny = new(ZAny)
		deriveCopyToPtrToZAny(this.ZAny, that.ZAny)
	}
	if this.Contains == nil {
		that.Contains = nil
	} else {
		that.Contains = new(Contains)
		deriveCopyToPtrToContains(this.Contains, that.Contains)
	}
	if this.Optional == nil {
		that.Optional = nil
	} else {
		that.Optional = new(Optional)
		deriveCopyToPtrToOptional(this.Optional, that.Optional)
	}
	if this.Interleave == nil {
		that.Interleave = nil
	} else {
		that.Interleave = new(Interleave)
		deriveCopyToPtrToInterleave(this.Interleave, that.Interleave)
	}
}

func deriveCopyToExpr(this, that *Expr) {
	if this.RightArrow == nil {
		that.RightArrow = nil
	} else {
		that.RightArrow = new(Keyword)
		deriveCopyToPtrToKeyword(this.RightArrow, that.RightArrow)
	}
	if this.Comma == nil {
		that.Comma = nil
	} else {
		that.Comma = new(Keyword)
		deriveCopyToPtrToKeyword(this.Comma, that.Comma)
	}
	if this.Terminal == nil {
		that.Terminal = nil
	} else {
		that.Terminal = new(Terminal)
		deriveCopyToPtrToTerminal(this.Terminal, that.Terminal)
	}
	if this.List == nil {
		that.List = nil
	} else {
		that.List = new(List)
		deriveCopyToPtrToList(this.List, that.List)
	}
	if this.Function == nil {
		that.Function = nil
	} else {
		that.Function = new(Function)
		deriveCopyToPtrToFunction(this.Function, that.Function)
	}
	if this.BuiltIn == nil {
		that.BuiltIn = nil
	} else {
		that.BuiltIn = new(BuiltIn)
		deriveCopyToPtrToBuiltIn(this.BuiltIn, that.BuiltIn)
	}
}

func deriveCopyToSliceOfPtrToPatternDecl(this, that []*PatternDecl) {
	for this_i, this_value := range this {
		if this_value == nil {
			that[this_i] = nil
		} else {
			that[this_i] = new(PatternDecl)
			deriveCopyToPtrToPatternDecl(this_value, that[this_i])
		}
	}
}

func deriveCopyToPtrToSpace(this, that *Space) {
	if this.Space == nil {
		that.Space = nil
	} else {
		if that.Space != nil {
			if len(this.Space) > len(that.Space) {
				if cap(that.Space) >= len(this.Space) {
					that.Space = (that.Space)[:len(this.Space)]
				} else {
					that.Space = make([]string, len(this.Space))
				}
			} else if len(this.Space) < len(that.Space) {
				that.Space = (that.Space)[:len(this.Space)]
			}
		} else {
			that.Space = make([]string, len(this.Space))
		}
		copy(that.Space, this.Space)
	}
}

func deriveCopyToPtrToEmpty(this, that *Empty) {
	if this.Empty == nil {
		that.Empty = nil
	} else {
		that.Empty = new(Keyword)
		deriveCopyToPtrToKeyword(this.Empty, that.Empty)
	}
}

func deriveCopyToPtrToTreeNode(this, that *TreeNode) {
	if this.Name == nil {
		that.Name = nil
	} else {
		that.Name = new(NameExpr)
		deriveCopyToPtrToNameExpr(this.Name, that.Name)
	}
	if this.Colon == nil {
		that.Colon = nil
	} else {
		that.Colon = new(Keyword)
		deriveCopyToPtrToKeyword(this.Colon, that.Colon)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
}

func deriveCopyToPtrToLeafNode(this, that *LeafNode) {
	if this.Expr == nil {
		that.Expr = nil
	} else {
		that.Expr = new(Expr)
		deriveCopyToExpr(this.Expr, that.Expr)
	}
}

func deriveCopyToPtrToConcat(this, that *Concat) {
	if this.OpenBracket == nil {
		that.OpenBracket = nil
	} else {
		that.OpenBracket = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenBracket, that.OpenBracket)
	}
	if this.LeftPattern == nil {
		that.LeftPattern = nil
	} else {
		that.LeftPattern = new(Pattern)
		deriveCopyToPattern(this.LeftPattern, that.LeftPattern)
	}
	if this.Comma == nil {
		that.Comma = nil
	} else {
		that.Comma = new(Keyword)
		deriveCopyToPtrToKeyword(this.Comma, that.Comma)
	}
	if this.RightPattern == nil {
		that.RightPattern = nil
	} else {
		that.RightPattern = new(Pattern)
		deriveCopyToPattern(this.RightPattern, that.RightPattern)
	}
	if this.ExtraComma == nil {
		that.ExtraComma = nil
	} else {
		that.ExtraComma = new(Keyword)
		deriveCopyToPtrToKeyword(this.ExtraComma, that.ExtraComma)
	}
	if this.CloseBracket == nil {
		that.CloseBracket = nil
	} else {
		that.CloseBracket = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseBracket, that.CloseBracket)
	}
}

func deriveCopyToPtrToOr(this, that *Or) {
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.LeftPattern == nil {
		that.LeftPattern = nil
	} else {
		that.LeftPattern = new(Pattern)
		deriveCopyToPattern(this.LeftPattern, that.LeftPattern)
	}
	if this.Pipe == nil {
		that.Pipe = nil
	} else {
		that.Pipe = new(Keyword)
		deriveCopyToPtrToKeyword(this.Pipe, that.Pipe)
	}
	if this.RightPattern == nil {
		that.RightPattern = nil
	} else {
		that.RightPattern = new(Pattern)
		deriveCopyToPattern(this.RightPattern, that.RightPattern)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}

func deriveCopyToPtrToAnd(this, that *And) {
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.LeftPattern == nil {
		that.LeftPattern = nil
	} else {
		that.LeftPattern = new(Pattern)
		deriveCopyToPattern(this.LeftPattern, that.LeftPattern)
	}
	if this.Ampersand == nil {
		that.Ampersand = nil
	} else {
		that.Ampersand = new(Keyword)
		deriveCopyToPtrToKeyword(this.Ampersand, that.Ampersand)
	}
	if this.RightPattern == nil {
		that.RightPattern = nil
	} else {
		that.RightPattern = new(Pattern)
		deriveCopyToPattern(this.RightPattern, that.RightPattern)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}

func deriveCopyToPtrToZeroOrMore(this, that *ZeroOrMore) {
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
	if this.Star == nil {
		that.Star = nil
	} else {
		that.Star = new(Keyword)
		deriveCopyToPtrToKeyword(this.Star, that.Star)
	}
}

func deriveCopyToPtrToReference(this, that *Reference) {
	if this.At == nil {
		that.At = nil
	} else {
		that.At = new(Keyword)
		deriveCopyToPtrToKeyword(this.At, that.At)
	}
	that.Name = this.Name
}

func deriveCopyToPtrToNot(this, that *Not) {
	if this.Exclamation == nil {
		that.Exclamation = nil
	} else {
		that.Exclamation = new(Keyword)
		deriveCopyToPtrToKeyword(this.Exclamation, that.Exclamation)
	}
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}

func deriveCopyToPtrToZAny(this, that *ZAny) {
	if this.Star == nil {
		that.Star = nil
	} else {
		that.Star = new(Keyword)
		deriveCopyToPtrToKeyword(this.Star, that.Star)
	}
}

func deriveCopyToPtrToContains(this, that *Contains) {
	if this.Dot == nil {
		that.Dot = nil
	} else {
		that.Dot = new(Keyword)
		deriveCopyToPtrToKeyword(this.Dot, that.Dot)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
}

func deriveCopyToPtrToOptional(this, that *Optional) {
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
	if this.QuestionMark == nil {
		that.QuestionMark = nil
	} else {
		that.QuestionMark = new(Keyword)
		deriveCopyToPtrToKeyword(this.QuestionMark, that.QuestionMark)
	}
}

func deriveCopyToPtrToInterleave(this, that *Interleave) {
	if this.OpenCurly == nil {
		that.OpenCurly = nil
	} else {
		that.OpenCurly = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenCurly, that.OpenCurly)
	}
	if this.LeftPattern == nil {
		that.LeftPattern = nil
	} else {
		that.LeftPattern = new(Pattern)
		deriveCopyToPattern(this.LeftPattern, that.LeftPattern)
	}
	if this.SemiColon == nil {
		that.SemiColon = nil
	} else {
		that.SemiColon = new(Keyword)
		deriveCopyToPtrToKeyword(this.SemiColon, that.SemiColon)
	}
	if this.RightPattern == nil {
		that.RightPattern = nil
	} else {
		that.RightPattern = new(Pattern)
		deriveCopyToPattern(this.RightPattern, that.RightPattern)
	}
	if this.ExtraSemiColon == nil {
		that.ExtraSemiColon = nil
	} else {
		that.ExtraSemiColon = new(Keyword)
		deriveCopyToPtrToKeyword(this.ExtraSemiColon, that.ExtraSemiColon)
	}
	if this.CloseCurly == nil {
		that.CloseCurly = nil
	} else {
		that.CloseCurly = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseCurly, that.CloseCurly)
	}
}

func deriveCopyToPtrToKeyword(this, that *Keyword) {
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	that.Value = this.Value
}

func deriveCopyToPtrToTerminal(this, that *Terminal) {
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	that.Literal = this.Literal
	if this.DoubleValue == nil {
		that.DoubleValue = nil
	} else {
		that.DoubleValue = new(float64)
		*that.DoubleValue = *this.DoubleValue
	}
	if this.IntValue == nil {
		that.IntValue = nil
	} else {
		that.IntValue = new(int64)
		*that.IntValue = *this.IntValue
	}
	if this.UintValue == nil {
		that.UintValue = nil
	} else {
		that.UintValue = new(uint64)
		*that.UintValue = *this.UintValue
	}
	if this.BoolValue == nil {
		that.BoolValue = nil
	} else {
		that.BoolValue = new(bool)
		*that.BoolValue = *this.BoolValue
	}
	if this.StringValue == nil {
		that.StringValue = nil
	} else {
		that.StringValue = new(string)
		*that.StringValue = *this.StringValue
	}
	if this.BytesValue == nil {
		that.BytesValue = nil
	} else {
		if that.BytesValue != nil {
			if len(this.BytesValue) > len(that.BytesValue) {
				if cap(that.BytesValue) >= len(this.BytesValue) {
					that.BytesValue = (that.BytesValue)[:len(this.BytesValue)]
				} else {
					that.BytesValue = make([]byte, len(this.BytesValue))
				}
			} else if len(this.BytesValue) < len(that.BytesValue) {
				that.BytesValue = (that.BytesValue)[:len(this.BytesValue)]
			}
		} else {
			that.BytesValue = make([]byte, len(this.BytesValue))
		}
		copy(that.BytesValue, this.BytesValue)
	}
	if this.Variable == nil {
		that.Variable = nil
	} else {
		that.Variable = new(Variable)
		*that.Variable = *this.Variable
	}
}

func deriveCopyToPtrToList(this, that *List) {
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	that.Type = this.Type
	if this.OpenCurly == nil {
		that.OpenCurly = nil
	} else {
		that.OpenCurly = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenCurly, that.OpenCurly)
	}
	if this.Elems == nil {
		that.Elems = nil
	} else {
		if that.Elems != nil {
			if len(this.Elems) > len(that.Elems) {
				if cap(that.Elems) >= len(this.Elems) {
					that.Elems = (that.Elems)[:len(this.Elems)]
				} else {
					that.Elems = make([]*Expr, len(this.Elems))
				}
			} else if len(this.Elems) < len(that.Elems) {
				that.Elems = (that.Elems)[:len(this.Elems)]
			}
		} else {
			that.Elems = make([]*Expr, len(this.Elems))
		}
		deriveCopyToSliceOfPtrToExpr(this.Elems, that.Elems)
	}
	if this.CloseCurly == nil {
		that.CloseCurly = nil
	} else {
		that.CloseCurly = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseCurly, that.CloseCurly)
	}
}

func deriveCopyToPtrToFunction(this, that *Function) {
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	that.Name = this.Name
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Params == nil {
		that.Params = nil
	} else {
		if that.Params != nil {
			if len(this.Params) > len(that.Params) {
				if cap(that.Params) >= len(this.Params) {
					that.Params = (that.Params)[:len(this.Params)]
				} else {
					that.Params = make([]*Expr, len(this.Params))
				}
			} else if len(this.Params) < len(that.Params) {
				that.Params = (that.Params)[:len(this.Params)]
			}
		} else {
			that.Params = make([]*Expr, len(this.Params))
		}
		deriveCopyToSliceOfPtrToExpr(this.Params, that.Params)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}

func deriveCopyToPtrToBuiltIn(this, that *BuiltIn) {
	if this.Symbol == nil {
		that.Symbol = nil
	} else {
		that.Symbol = new(Keyword)
		deriveCopyToPtrToKeyword(this.Symbol, that.Symbol)
	}
	if this.Expr == nil {
		that.Expr = nil
	} else {
		that.Expr = new(Expr)
		deriveCopyToExpr(this.Expr, that.Expr)
	}
}

func deriveCopyToPtrToPatternDecl(this, that *PatternDecl) {
	if this.Hash == nil {
		that.Hash = nil
	} else {
		that.Hash = new(Keyword)
		deriveCopyToPtrToKeyword(this.Hash, that.Hash)
	}
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	that.Name = this.Name
	if this.Eq == nil {
		that.Eq = nil
	} else {
		that.Eq = new(Keyword)
		deriveCopyToPtrToKeyword(this.Eq, that.Eq)
	}
	if this.Pattern == nil {
		that.Pattern = nil
	} else {
		that.Pattern = new(Pattern)
		deriveCopyToPattern(this.Pattern, that.Pattern)
	}
}

func deriveCopyToPtrToNameExpr(this, that *NameExpr) {
	if this.Name == nil {
		that.Name = nil
	} else {
		that.Name = new(Name)
		deriveCopyToPtrToName(this.Name, that.Name)
	}
	if this.AnyName == nil {
		that.AnyName = nil
	} else {
		that.AnyName = new(AnyName)
		deriveCopyToPtrToAnyName(this.AnyName, that.AnyName)
	}
	if this.AnyNameExcept == nil {
		that.AnyNameExcept = nil
	} else {
		that.AnyNameExcept = new(AnyNameExcept)
		deriveCopyToPtrToAnyNameExcept(this.AnyNameExcept, that.AnyNameExcept)
	}
	if this.NameChoice == nil {
		that.NameChoice = nil
	} else {
		that.NameChoice = new(NameChoice)
		deriveCopyToPtrToNameChoice(this.NameChoice, that.NameChoice)
	}
}

func deriveCopyToSliceOfPtrToExpr(this, that []*Expr) {
	for this_i, this_value := range this {
		if this_value == nil {
			that[this_i] = nil
		} else {
			that[this_i] = new(Expr)
			deriveCopyToExpr(this_value, that[this_i])
		}
	}
}

func deriveCopyToPtrToName(this, that *Name) {
	if this.Before == nil {
		that.Before = nil
	} else {
		that.Before = new(Space)
		deriveCopyToPtrToSpace(this.Before, that.Before)
	}
	if this.DoubleValue == nil {
		that.DoubleValue = nil
	} else {
		that.DoubleValue = new(float64)
		*that.DoubleValue = *this.DoubleValue
	}
	if this.IntValue == nil {
		that.IntValue = nil
	} else {
		that.IntValue = new(int64)
		*that.IntValue = *this.IntValue
	}
	if this.UintValue == nil {
		that.UintValue = nil
	} else {
		that.UintValue = new(uint64)
		*that.UintValue = *this.UintValue
	}
	if this.BoolValue == nil {
		that.BoolValue = nil
	} else {
		that.BoolValue = new(bool)
		*that.BoolValue = *this.BoolValue
	}
	if this.StringValue == nil {
		that.StringValue = nil
	} else {
		that.StringValue = new(string)
		*that.StringValue = *this.StringValue
	}
	if this.BytesValue == nil {
		that.BytesValue = nil
	} else {
		if that.BytesValue != nil {
			if len(this.BytesValue) > len(that.BytesValue) {
				if cap(that.BytesValue) >= len(this.BytesValue) {
					that.BytesValue = (that.BytesValue)[:len(this.BytesValue)]
				} else {
					that.BytesValue = make([]byte, len(this.BytesValue))
				}
			} else if len(this.BytesValue) < len(that.BytesValue) {
				that.BytesValue = (that.BytesValue)[:len(this.BytesValue)]
			}
		} else {
			that.BytesValue = make([]byte, len(this.BytesValue))
		}
		copy(that.BytesValue, this.BytesValue)
	}
}

func deriveCopyToPtrToAnyName(this, that *AnyName) {
	if this.Underscore == nil {
		that.Underscore = nil
	} else {
		that.Underscore = new(Keyword)
		deriveCopyToPtrToKeyword(this.Underscore, that.Underscore)
	}
}

func deriveCopyToPtrToAnyNameExcept(this, that *AnyNameExcept) {
	if this.Exclamation == nil {
		that.Exclamation = nil
	} else {
		that.Exclamation = new(Keyword)
		deriveCopyToPtrToKeyword(this.Exclamation, that.Exclamation)
	}
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Except == nil {
		that.Except = nil
	} else {
		that.Except = new(NameExpr)
		deriveCopyToPtrToNameExpr(this.Except, that.Except)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}

func deriveCopyToPtrToNameChoice(this, that *NameChoice) {
	if this.OpenParen == nil {
		that.OpenParen = nil
	} else {
		that.OpenParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.OpenParen, that.OpenParen)
	}
	if this.Left == nil {
		that.Left = nil
	} else {
		that.Left = new(NameExpr)
		deriveCopyToPtrToNameExpr(this.Left, that.Left)
	}
	if this.Pipe == nil {
		that.Pipe = nil
	} else {
		that.Pipe = new(Keyword)
		deriveCopyToPtrToKeyword(this.Pipe, that.Pipe)
	}
	if this.Right == nil {
		that.Right = nil
	} else {
		that.Right = new(NameExpr)
		deriveCopyToPtrToNameExpr(this.Right, that.Right)
	}
	if this.CloseParen == nil {
		that.CloseParen = nil
	} else {
		that.CloseParen = new(Keyword)
		deriveCopyToPtrToKeyword(this.CloseParen, that.CloseParen)
	}
}
