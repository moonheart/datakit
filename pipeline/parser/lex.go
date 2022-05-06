// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Item struct {
	Typ ItemType
	Pos Pos
	Val string
}

func (i *Item) PositionRange() *PositionRange {
	return &PositionRange{
		Start: i.Pos,
		End:   i.Pos + Pos(len(i.Val)),
	}
}

func (i Item) lexStr() string {
	return fmt.Sprintf("% 06d %02d %s", i.Typ, i.Pos, i.String())
}

func (i Item) String() string {
	switch {
	case i.Typ == EOF:
		return "EOF"
	case i.Typ == ERROR:
		return i.Val
	case i.Typ == ID:
		return fmt.Sprintf("%q", i.Val)
	case i.Typ.IsKeyword():
		return fmt.Sprintf("<%s>", i.Val)
	case i.Typ.IsOperator():
		return fmt.Sprintf("<op:'%s'>", i.Val)
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

func (i ItemType) IsOperator() bool { return i > operatorsStart && i < operatorsEnd }
func (i ItemType) IsKeyword() bool  { return i > keywordsStart && i < keywordsEnd }

type ItemType int

const (
	eof         = -1
	lineComment = "#"
	Digits      = "0123456789"
	HexDigits   = "0123456789abcdefABCDEF"
)

var (
	keywords = map[string]ItemType{
		// Keywords.
		"if":         IF,
		"elif":       ELIF,
		"else":       ELSE,
		"false":      FALSE,
		"identifier": IDENTIFIER,
		"nil":        NIL,
		"null":       NULL,
		"re":         RE,
		"true":       TRUE,
		"jp":         JP,
	}

	ItemTypeStr = map[ItemType]string{
		LEFT_PAREN:    "(",
		RIGHT_PAREN:   ")",
		LEFT_BRACKET:  "[",
		RIGHT_BRACKET: "]",
		LEFT_BRACE:    "{",
		RIGHT_BRACE:   "}",
		COMMA:         ",",
		EQ:            "=",
		EQEQ:          "==",
		SEMICOLON:     ";",
		DOT:           ".",
		SPACE:         "<space>",

		SUB: "-",
		ADD: "+",
		MUL: "*",
		MOD: "%",
		DIV: "/",
		NEQ: "!=",
		LTE: "<=",
		LT:  "<",
		GTE: ">=",
		GT:  ">",
		POW: "^",
		AND: "&&",
		OR:  "||",
	}
)

func init() { //nolint:gochecknoinits
	// Add keywords to Item type strings.
	for s, ty := range keywords {
		ItemTypeStr[ty] = s
	}
	// Special numbers.
	keywords["inf"] = NUMBER
	keywords["nan"] = NUMBER
}

func (i ItemType) String() string {
	if s, ok := ItemTypeStr[i]; ok {
		return s
	}
	return fmt.Sprintf("<Item %d>", i)
}

func (i Item) desc() string {
	if _, ok := ItemTypeStr[i.Typ]; ok {
		return i.String()
	}
	if i.Typ == EOF {
		return i.Typ.desc()
	}
	return fmt.Sprintf("%s %s", i.Typ.desc(), i)
}

func (i ItemType) desc() string {
	switch i {
	case ERROR:
		return "error"
	case EOF:
		return "end of input"
	case COMMENT:
		return "comment"
	case ID:
		return "id"
	case STRING:
		return "string"
	case NUMBER:
		return "number"
	}
	return fmt.Sprintf("%q", i)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Pos is the position in a string.
// Negative numbers indicate undefined positions.
type Pos int

// Lexer holds the state of the scanner.
type Lexer struct {
	input       string  // The string being scanned.
	state       stateFn // The next lexing function to enter.
	pos         Pos     // Current position in the input.
	start       Pos     // Start position of this Item.
	width       Pos     // Width of last rune read from input.
	lastPos     Pos     // Position of most recent Item returned by NextItem.
	itemp       *Item   // Pointer to where the next scanned item should be placed.
	scannedItem bool    // Set to true every time an item is scanned.

	parenDepth   int // nested depth of () exprs.
	braceDepth   int // nested depth of {} exprs.
	bracketDepth int // nested depth of [] exprs.

	stringOpen    rune // Quote rune of the string currently being read.
	backquoteOpen rune // backquote keyworkds and utf8 characters

	// seriesDesc is set when a series description for the testing
	// language is lexed.
	// seriesDesc bool
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input: input,
		state: lexStatements,
	}
	return l
}

// Lexer entry.
func lexStatements(l *Lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], lineComment) {
		return lexLineComment
	}

	switch r := l.next(); {
	case r == ',':
		l.emit(COMMA)

	case isSpace(r):
		return lexSpace

	case r == '*':
		l.emit(MUL)

	case r == '/':
		l.emit(DIV)

	case r == '%':
		l.emit(MOD)

	case r == '+':
		l.emit(ADD)

	case r == '-':
		l.emit(SUB)

	case r == '^':
		l.emit(POW)

	case r == '=':
		if t := l.peek(); t == '=' {
			l.next()
			l.emit(EQEQ)
		} else {
			l.emit(EQ)
		}

	case r == ';':
		l.emit(SEMICOLON)

	case r == '.':
		l.emit(DOT)

	case r == '|':
		if t := l.peek(); t == '|' {
			l.next()
			l.emit(OR)
		} else {
			// TODO: add bit-or operator
			return l.errorf("unexpected character `%q' after `!'", r)
		}

	case r == '&':
		if t := l.peek(); t == '&' {
			l.next()
			l.emit(AND)
		} else {
			// TODO: add bit-and operator
			return l.errorf("unexpected character `%q' after `!'", r)
		}

	case r == '!':
		switch nr := l.next(); {
		case nr == '=':
			l.emit(NEQ)
		default:
			return l.errorf("unexpected character `%q' after `!'", nr)
		}

	case r == '<':
		if t := l.peek(); t == '=' {
			l.next()
			l.emit(LTE)
		} else {
			l.emit(LT)
		}

	case r == '>':
		if t := l.peek(); t == '=' {
			l.next()
			l.emit(GTE)
		} else {
			l.emit(GT)
		}

	case isDigit(r) || (r == '.' && isDigit(l.peek())):
		l.backup()
		return lexNumberOrDuration

	case r == '"':
		if t1 := l.peek(); t1 == '"' {
			l.next()
			if t2 := l.peek(); t2 == '"' {
				l.next()
				return lexMultilineString
			} else {
				l.emit(STRING)
			}
		} else {
			l.stringOpen = r
			return lexString
		}

	case r == '\'':
		if t := l.peek(); t == '\'' {
			l.next()
			if t := l.peek(); t == '\'' {
				l.next()
				return lexMultilineString
			} else {
				l.emit(STRING)
			}
		} else {
			l.stringOpen = r
			return lexString
		}

	case r == '`':
		l.backquoteOpen = r
		return lexRawString

	case isAlpha(r) || isUTF8(r):
		l.backup()
		return lexKeywordOrIdentifier

	case r == '(':
		l.emit(LEFT_PAREN)
		l.parenDepth++
		return lexStatements

	case r == ')':
		l.emit(RIGHT_PAREN)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right parenthesis %q", r)
		}
		return lexStatements

	case r == '{':
		l.emit(LEFT_BRACE)
		l.braceDepth++
		return lexStatements

	case r == '}':
		l.emit(RIGHT_BRACE)
		l.braceDepth--
		if l.braceDepth < 0 {
			return l.errorf("unexpected right parenthesis %q", r)
		}
		return lexStatements

	case r == '[':
		l.bracketDepth++
		l.emit(LEFT_BRACKET)

	case r == ']':
		l.bracketDepth--
		l.emit(RIGHT_BRACKET)

	case r == eof:
		switch {
		case l.parenDepth != 0:
			return l.errorf("unclosed left parenthesis")
		case l.bracketDepth != 0:
			return l.errorf("unclosed left bracket")
		case l.braceDepth != 0:
			return l.errorf("unclosed left brace")
		}

		l.emit(EOF)
		return nil

	default:
		return l.errorf("unexpected character: %q", r)
	}
	return lexStatements
}

// Other state functions
// scan alphanumberic identifier, maybe keyword.
func lexKeywordOrIdentifier(l *Lexer) stateFn {
__goon:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r), isUTF8(r):
			// absorb
		default:
			l.backup()
			word := l.input[l.start:l.pos]

			if kw, ok := keywords[strings.ToLower(word)]; ok {
				// log.Debugf("emit keyword: %s", kw)
				l.emit(kw)
			} else {
				l.emit(ID)
			}

			break __goon
		}
	}

	return lexStatements
}

func lexSpace(l *Lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}

	l.ignore()
	return lexStatements
}

func lexNumberOrDuration(l *Lexer) stateFn {
	if l.scanNumber() {
		l.emit(NUMBER)
		return lexStatements
	}

	return l.errorf("bad duration: %q", l.cur())
}

func lexRawString(l *Lexer) stateFn {
__goon:
	for {
		switch l.next() {
		case utf8.RuneError:
			l.errorf("invalid UTF-8 rune")
			return lexRawString
		case eof:
			l.errorf("unterminated raw string")
			return lexRawString
		case l.backquoteOpen:
			break __goon
		}
	}

	l.emit(QUOTED_STRING)
	return lexStatements
}

func lexLineComment(l *Lexer) stateFn {
	l.pos += Pos(len(lineComment))
	for r := l.next(); !isEOL(r) && r != eof; {
		r = l.next()
	}
	l.backup()
	l.emit(COMMENT)
	return lexStatements
}

func lexEscape(l *Lexer) stateFn {
	ch := l.next()
	var n int
	var base, max uint32

	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', l.stringOpen, l.backquoteOpen:
		return lexString
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x', 'X':
		ch = l.next()
		n, base, max = 2, 16, 255
	case 'u':
		ch = l.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		ch = l.next()
		n, base, max = 8, 16, unicode.MaxRune
	case eof:
		l.errorf("escape squence not terminated")
		return lexString
	default:
		l.errorf("unknown escape sequence %#U", ch)
		return lexString
	}

	log.Debugf("n: %d, base: %d, max: %d", n, base, max)

	var x uint32
	for n > 0 {
		d := uint32(digitVal(ch))
		if d >= base {
			if ch == eof {
				l.errorf("escape sequence not terminated")
			}
			l.errorf("illegal character %#U in escape sequence", ch)
			return lexString
		}

		x = x*base + d
		ch = l.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		l.errorf("escape sequence is an invalid Unicode code point")
	}

	log.Debugf("get number %d", x)
	return lexString
}

func lexMultilineString(l *Lexer) stateFn {
__goon:
	for {
		c := l.next()
		switch c {
		case eof:
			return l.errorf("unterminated quoted string within lexMultilineString")

		case utf8.RuneError:
			l.errorf("invalid UTF-8 rune")

		case '"', '\'':
			if t := l.peek(); t == '"' || t == '\'' {
				l.next()
				if t := l.peek(); t == '"' || t == '\'' {
					l.next()
					break __goon
				}
			}
		default: // pass
		}
	}

	l.emit(MULTILINE_STRING)
	return lexStatements
}

func lexString(l *Lexer) stateFn {
__goon:
	for {
		switch l.next() {
		case '\\':
			return lexEscape

		case utf8.RuneError:
			l.errorf("invalid UTF-8 rune")

		case eof, '\n':
			return l.errorf("unterminated quoted string within lexString")

		case l.stringOpen:
			break __goon
		}
	}

	l.emit(STRING)
	return lexStatements
}

// lexer tool functions.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) emit(t ItemType) {
	*l.itemp = Item{t, l.start, l.input[l.start:l.pos]}

	l.start = l.pos
	l.scannedItem = true
	log.Debugf("emit: %+#v", l.itemp)
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	*l.itemp = Item{ERROR, l.start, fmt.Sprintf(format, args...)}
	l.scannedItem = true

	return nil
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() { l.pos -= l.width }

func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
		/* consume */
	}
	l.backup()
}

func (l *Lexer) NextItem(itemp *Item) {
	l.scannedItem = false
	l.itemp = itemp

	if l.state != nil {
		for !l.scannedItem {
			l.state = l.state(l)
		}
	} else {
		l.emit(EOF)
	}

	l.lastPos = l.itemp.Pos
}

func (l *Lexer) cur() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) scanNumber() bool {
	digs := Digits
	if l.accept("0") && l.accept("xX") {
		digs = HexDigits
	}

	l.acceptRun(digs)
	if l.accept(".") {
		l.acceptRun(digs)
	}

	if l.accept("eE") { // scientific notation
		l.accept("+-")
		l.acceptRun(Digits)
	}

	// next things should not be alphanumberic
	if r := l.peek(); !isAlphaNumeric(r) {
		return true
	}

	return false
}

// helpers.
func isAlphaNumeric(r rune) bool { return isAlpha(r) || isDigit(r) }
func isAlpha(r rune) bool        { return r == '_' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') }
func isUTF8(r rune) bool         { return utf8.RuneLen(r) > 1 }
func isDigit(r rune) bool        { return '0' <= r && r <= '9' }
func isSpace(r rune) bool        { return r == ' ' || r == '\t' || r == '\n' || r == '\r' }
func isEOL(r rune) bool          { return r == '\r' || r == '\n' }

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}
