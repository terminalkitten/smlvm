package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseCases(p *parser) []*ast.Case {
	var ret []*ast.Case
	for !(p.SeeOp("}") || p.See(lexing.EOF)) {
		if c := parseCase(p); c != nil {
			ret = append(ret, c)
		}
		p.skipErrStmt()
	}
	return ret
}

func parseCase(p *parser) *ast.Case {
	ret := new(ast.Case)
	if p.SeeKeyword("case") {
		ret.Kw = p.Shift()
		ret.Expr = parseExpr(p)
		if ret.Expr == nil {
			return nil
		}
	} else if p.SeeKeyword("default") {
		ret.Kw = p.Shift()
		ret.Expr = nil
	} else {
		p.CodeErrorfHere("pl.missingCaseInSwitch",
			"must start with keyword case/default in switch")
		return nil
	}
	ret.Colon = p.ExpectOp(":")
	if ret.Colon == nil {
		return nil
	}
	for !(p.SeeKeyword("case") || p.SeeKeyword("default") ||
		p.SeeOp("}") || p.See(lexing.EOF)) {
		if p.SeeKeyword("fallthrough") {
			return parseFallthrough(p, ret)
		}
		if stmt := p.parseStmt(); stmt != nil {
			ret.Stmts = append(ret.Stmts, stmt)
		}
		p.skipErrStmt()
	}
	return ret
}

func parseFallthrough(p *parser, ret *ast.Case) *ast.Case {
	fall := new(ast.FallthroughStmt)
	fall.Kw = p.Shift()
	fall.Semi = p.ExpectSemi()
	ret.Fall = fall
	if p.InError() {
		return nil
	}
	if p.SeeKeyword("case") || p.SeeKeyword("default") {
		return ret
	}
	p.CodeErrorfHere("pl.wrongFallthroughPos",
		"fallthrough out of place, must be the last statement"+
			"of a CASE and cannot be in the final CASE in a switch")
	return nil
}
