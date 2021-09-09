// Code generated by goyacc -l -o pipeline/parser/parser_y.go pipeline/parser/parser.y. DO NOT EDIT.
package parser

import __yyfmt__ "fmt"

import (
	"fmt"
)

type yySymType struct {
	yys     int
	node    Node
	nodes   []Node
	item    Item
	strings []string
	float   float64
}

const SEMICOLON = 57346
const COMMA = 57347
const COMMENT = 57348
const DOT = 57349
const EOF = 57350
const ERROR = 57351
const ID = 57352
const LEFT_PAREN = 57353
const LEFT_BRACKET = 57354
const NUMBER = 57355
const RIGHT_PAREN = 57356
const RIGHT_BRACKET = 57357
const SPACE = 57358
const STRING = 57359
const QUOTED_STRING = 57360
const operatorsStart = 57361
const ADD = 57362
const DIV = 57363
const GTE = 57364
const GT = 57365
const LT = 57366
const LTE = 57367
const MOD = 57368
const MUL = 57369
const NEQ = 57370
const EQ = 57371
const POW = 57372
const SUB = 57373
const operatorsEnd = 57374
const keywordsStart = 57375
const TRUE = 57376
const FALSE = 57377
const IDENTIFIER = 57378
const AND = 57379
const OR = 57380
const NIL = 57381
const NULL = 57382
const RE = 57383
const JP = 57384
const keywordsEnd = 57385
const startSymbolsStart = 57386
const START_PIPELINE = 57387
const startSymbolsEnd = 57388

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SEMICOLON",
	"COMMA",
	"COMMENT",
	"DOT",
	"EOF",
	"ERROR",
	"ID",
	"LEFT_PAREN",
	"LEFT_BRACKET",
	"NUMBER",
	"RIGHT_PAREN",
	"RIGHT_BRACKET",
	"SPACE",
	"STRING",
	"QUOTED_STRING",
	"operatorsStart",
	"ADD",
	"DIV",
	"GTE",
	"GT",
	"LT",
	"LTE",
	"MOD",
	"MUL",
	"NEQ",
	"EQ",
	"POW",
	"SUB",
	"operatorsEnd",
	"keywordsStart",
	"TRUE",
	"FALSE",
	"IDENTIFIER",
	"AND",
	"OR",
	"NIL",
	"NULL",
	"RE",
	"JP",
	"keywordsEnd",
	"startSymbolsStart",
	"START_PIPELINE",
	"startSymbolsEnd",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 31,
	11, 48,
	-2, 59,
}

const yyPrivate = 57344

const yyLast = 206

var yyAct = [...]int{

	19, 17, 25, 26, 9, 56, 65, 9, 31, 32,
	35, 8, 10, 33, 8, 10, 16, 35, 42, 61,
	40, 60, 3, 35, 91, 100, 99, 86, 87, 44,
	11, 41, 64, 11, 38, 39, 46, 85, 43, 36,
	37, 53, 54, 98, 97, 56, 69, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 81, 82, 83,
	84, 70, 68, 67, 66, 2, 63, 90, 93, 94,
	95, 88, 89, 9, 30, 65, 33, 14, 45, 46,
	35, 10, 13, 40, 53, 54, 4, 96, 56, 57,
	62, 1, 5, 33, 41, 24, 27, 38, 39, 11,
	40, 28, 36, 37, 29, 9, 30, 18, 33, 20,
	21, 41, 35, 10, 22, 40, 23, 6, 59, 15,
	12, 7, 34, 0, 0, 0, 41, 0, 0, 38,
	39, 11, 92, 0, 36, 37, 29, 0, 45, 46,
	47, 48, 51, 52, 53, 54, 55, 58, 56, 57,
	0, 0, 0, 0, 0, 49, 50, 45, 46, 47,
	48, 51, 52, 53, 54, 55, 58, 56, 57, 0,
	0, 0, 0, 0, 49, 50, 45, 46, 47, 48,
	51, 52, 53, 54, 55, 58, 56, 57, 0, 0,
	0, 0, 0, 49, 45, 46, 47, 48, 51, 52,
	53, 54, 55, 58, 56, 57,
}
var yyPact = [...]int{

	20, 78, -3, -1000, -1000, -3, -1000, 71, -1000, -1000,
	-1000, 66, -1000, 95, -7, 24, -1000, 137, 0, -1000,
	-1000, -1000, -1000, -1000, 83, -1000, -1000, -1000, -1000, 55,
	63, 52, 51, -1000, 49, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 32, -1000, 95, 63, 63, 63, 63, 63,
	63, 63, 63, 63, 63, 63, 63, 63, 63, 22,
	13, -1000, -6, 6, 118, 80, 80, 80, -1000, -1000,
	-1000, 15, -25, 58, 58, 174, 156, 58, 58, -25,
	-25, 58, -25, 15, 58, -1000, 0, -1000, 52, 51,
	30, 29, -1000, 13, 11, 10, -1000, -1000, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 122, 121, 8, 119, 0, 118, 116, 1, 16,
	114, 110, 109, 101, 3, 96, 2, 9, 95, 92,
	91,
}
var yyR1 = [...]int{

	0, 20, 20, 20, 19, 19, 8, 8, 8, 8,
	8, 8, 1, 1, 14, 15, 15, 13, 13, 11,
	10, 4, 4, 4, 4, 9, 9, 6, 6, 6,
	5, 5, 5, 5, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 2, 16,
	16, 12, 12, 3, 3, 3, 17, 17, 17, 18,
	18, 18, 18,
}
var yyR2 = [...]int{

	0, 2, 2, 1, 1, 2, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	4, 3, 2, 1, 0, 1, 3, 3, 1, 0,
	1, 1, 1, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 1, 1,
	2, 4, 4, 1, 1, 4, 4, 4, 3, 1,
	1, 3, 3,
}
var yyChk = [...]int{

	-1000, -20, 45, 2, 8, -19, -10, -2, -3, 10,
	18, 36, -10, 11, 11, -4, -9, -8, 12, -5,
	-12, -11, -10, -7, -18, -16, -14, -15, -13, 41,
	11, -3, -17, 13, -1, 17, 39, 40, 34, 35,
	20, 31, -14, 14, 5, 20, 21, 22, 23, 37,
	38, 24, 25, 26, 27, 28, 30, 31, 29, -6,
	-16, -5, 7, 11, -8, 12, 12, 12, 13, 14,
	-9, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, 15, 5, 15, -3, -17,
	-14, 18, 14, -16, -16, -16, -5, 14, 14, 15,
	15,
}
var yyDef = [...]int{

	0, -2, 0, 3, 2, 1, 4, 0, 48, 53,
	54, 0, 5, 24, 0, 0, 23, 25, 29, 6,
	7, 8, 9, 10, 11, 30, 31, 32, 33, 0,
	0, -2, 60, 49, 0, 14, 15, 16, 17, 18,
	12, 13, 0, 20, 22, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	30, 28, 0, 0, 0, 0, 0, 0, 50, 55,
	21, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	43, 44, 45, 46, 47, 26, 0, 58, 61, 62,
	0, 0, 19, 0, 0, 0, 27, 51, 52, 56,
	57,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yylex.(*parser).parseResult = yyDollar[2].node
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yylex.(*parser).unexpected("", "")
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &Ast{Functions: []*FuncExpr{yyDollar[1].node.(*FuncExpr)}}
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			ast := yyDollar[1].node.(*Ast)
			ast.Functions = append(ast.Functions, yyDollar[2].node.(*FuncExpr))
			yyVAL.node = yyDollar[1].node
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &StringLiteral{Val: yylex.(*parser).unquoteString(yyDollar[1].item.Val)}
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &NilLiteral{}
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &NilLiteral{}
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &BoolLiteral{Val: true}
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &BoolLiteral{Val: false}
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = &ParenExpr{Param: yyDollar[2].node}
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.node = yylex.(*parser).newFunc(yyDollar[1].item.Val, yyDollar[3].nodes)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.nodes = append(yyVAL.nodes, yyDollar[3].node)
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.nodes = []Node{yyDollar[1].node}
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
		{
			yyVAL.nodes = nil
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = yyDollar[1].node
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = getFuncArgList(yyDollar[2].node.(NodeList))
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			nl := yyVAL.node.(NodeList)
			nl = append(nl, yyDollar[3].node)
			yyVAL.node = nl
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = NodeList{yyDollar[1].node}
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
		{
			yyVAL.node = NodeList{}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			yyVAL.node = bexpr
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			yyVAL.node = bexpr
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			yyVAL.node = bexpr
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			yyVAL.node = bexpr
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			bexpr := yylex.(*parser).newBinExpr(yyDollar[1].node, yyDollar[3].node, yyDollar[2].item)
			bexpr.ReturnBool = true
			yyVAL.node = bexpr
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.item = yyDollar[1].item
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = yylex.(*parser).number(yyDollar[1].item.Val)
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			num := yylex.(*parser).number(yyDollar[2].item.Val)
			switch yyDollar[1].item.Typ {
			case ADD: // pass
			case SUB:
				if num.IsInt {
					num.Int = -num.Int
				} else {
					num.Float = -num.Float
				}
			}
			yyVAL.node = num
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.node = &Regex{Regex: yyDollar[3].node.(*StringLiteral).Val}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.node = &Regex{Regex: yylex.(*parser).unquoteString(yyDollar[3].item.Val)}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.item.Val = yylex.(*parser).unquoteString(yyDollar[1].item.Val)
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.item.Val = yyDollar[3].node.(*StringLiteral).Val
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			nl := yyDollar[3].node.(*NumberLiteral)
			if !nl.IsInt {
				yylex.(*parser).addParseErr(nil,
					fmt.Errorf("array index should be int, got `%f'", nl.Float))
				yyVAL.node = nil
			} else {
				yyVAL.node = &IndexExpr{Obj: &Identifier{Name: yyDollar[1].item.Val}, Index: []int64{nl.Int}}
			}
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
		{

			nl := yyDollar[3].node.(*NumberLiteral)
			if !nl.IsInt {
				yylex.(*parser).addParseErr(nil,
					fmt.Errorf("array index should be int, got `%f'", nl.Float))
				yyVAL.node = nil
			} else {
				in := yyDollar[1].node.(*IndexExpr)
				in.Index = append(in.Index, nl.Int)
				yyVAL.node = in
			}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			nl := yyDollar[2].node.(*NumberLiteral)
			if !nl.IsInt {
				yylex.(*parser).addParseErr(nil,
					fmt.Errorf("array index should be int, got `%f'", nl.Float))
				yyVAL.node = nil
			} else {
				yyVAL.node = &IndexExpr{Index: []int64{nl.Int}}
			}
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &AttrExpr{Obj: &Identifier{Name: yyDollar[1].item.Val}}
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.node = &AttrExpr{Obj: yyDollar[1].node}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = &AttrExpr{Obj: yyDollar[1].node, Attr: &Identifier{Name: yyDollar[3].item.Val}}
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.node = &AttrExpr{Obj: yyDollar[1].node, Attr: yyDollar[3].node}
		}
	}
	goto yystack /* stack new state and value */
}
