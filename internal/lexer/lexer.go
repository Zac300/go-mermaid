// Package lexer turns Mermaid flowchart source into a flat token stream.
// It captures shape and edge-label text as single Text tokens so the
// parser can stay grammar-focused.
package lexer

import (
	"strings"
	"unicode"

	"github.com/Zac300/go-mermaid/internal/syntax"
)

// Lex scans src and returns its tokens, always terminated by an EOF token.
func Lex(src string) ([]Token, error) {
	l := &lexer{src: []rune(src), line: 1, col: 1}
	return l.run()
}

type lexer struct {
	src  []rune
	pos  int
	line int
	col  int
}

var keywords = map[string]bool{
	"graph": true, "flowchart": true, "subgraph": true, "end": true,
}

func (l *lexer) run() ([]Token, error) {
	var toks []Token
	for l.pos < len(l.src) {
		r := l.src[l.pos]
		switch {
		case r == '\n':
			toks = append(toks, l.emit(Newline, "\n"))
			l.advance()
		case r == '\r':
			l.advance()
		case r == ' ' || r == '\t':
			l.advance()
		case r == '%' && l.peek(1) == '%':
			l.skipComment() // %% line comment
		case r == '|':
			tok, err := l.lexPipeLabel()
			if err != nil {
				return nil, err
			}
			toks = append(toks, tok...)
		case isShapeOpen(r):
			tok, err := l.lexShape()
			if err != nil {
				return nil, err
			}
			toks = append(toks, tok...)
		case isConnectorRune(r):
			toks = append(toks, l.lexArrow())
		case isIdentRune(r):
			toks = append(toks, l.lexIdent())
		default:
			return nil, syntax.Errorf(l.line, l.col, "unexpected character %q", string(r))
		}
	}
	toks = append(toks, Token{Kind: EOF, Line: l.line, Col: l.col})
	return toks, nil
}

func (l *lexer) lexIdent() Token {
	start, col := l.pos, l.col
	for l.pos < len(l.src) && isIdentRune(l.src[l.pos]) {
		l.advance()
	}
	val := string(l.src[start:l.pos])
	k := Ident
	if keywords[val] {
		k = Keyword
	}
	return Token{Kind: k, Val: val, Line: l.line, Col: col}
}

// lexArrow consumes a maximal run of connector runes (e.g. -->, -.->, ==>).
func (l *lexer) lexArrow() Token {
	start, col := l.pos, l.col
	for l.pos < len(l.src) && isConnectorRune(l.src[l.pos]) {
		l.advance()
	}
	return Token{Kind: Arrow, Val: string(l.src[start:l.pos]), Line: l.line, Col: col}
}

// lexShape captures a node shape: opener, inner text, closer.
func (l *lexer) lexShape() ([]Token, error) {
	col := l.col
	open := l.readOpener()
	close := shapeCloser(open)
	start := l.pos
	for l.pos < len(l.src) && !l.hasPrefix(close) {
		if l.src[l.pos] == '\n' {
			return nil, syntax.Errorf(l.line, col, "unterminated shape %q", open)
		}
		l.advance()
	}
	if l.pos >= len(l.src) {
		return nil, syntax.Errorf(l.line, col, "unterminated shape %q", open)
	}
	text := strings.TrimSpace(string(l.src[start:l.pos]))
	textCol := col + len([]rune(open))
	for range close {
		l.advance()
	}
	return []Token{
		{Kind: ShapeOpen, Val: open, Line: l.line, Col: col},
		{Kind: Text, Val: stripQuotes(text), Line: l.line, Col: textCol},
		{Kind: ShapeClose, Val: close, Line: l.line, Col: l.col},
	}, nil
}

// lexPipeLabel captures |text| used for inline edge labels.
func (l *lexer) lexPipeLabel() ([]Token, error) {
	col := l.col
	l.advance() // opening |
	start := l.pos
	for l.pos < len(l.src) && l.src[l.pos] != '|' {
		if l.src[l.pos] == '\n' {
			return nil, syntax.Errorf(l.line, col, "unterminated label")
		}
		l.advance()
	}
	if l.pos >= len(l.src) {
		return nil, syntax.Errorf(l.line, col, "unterminated label")
	}
	text := strings.TrimSpace(string(l.src[start:l.pos]))
	l.advance() // closing |
	return []Token{
		{Kind: Pipe, Val: "|", Line: l.line, Col: col},
		{Kind: Text, Val: stripQuotes(text), Line: l.line, Col: col + 1},
		{Kind: Pipe, Val: "|", Line: l.line, Col: l.col - 1},
	}, nil
}

// readOpener consumes the longest valid opening delimiter at pos.
func (l *lexer) readOpener() string {
	for _, o := range []string{"([", "((", "[", "(", "{"} {
		if l.hasPrefix(o) {
			for range o {
				l.advance()
			}
			return o
		}
	}
	o := string(l.src[l.pos])
	l.advance()
	return o
}

func (l *lexer) skipComment() {
	for l.pos < len(l.src) && l.src[l.pos] != '\n' {
		l.advance()
	}
}

func (l *lexer) hasPrefix(s string) bool {
	rs := []rune(s)
	if l.pos+len(rs) > len(l.src) {
		return false
	}
	for i, r := range rs {
		if l.src[l.pos+i] != r {
			return false
		}
	}
	return true
}

func (l *lexer) peek(n int) rune {
	if l.pos+n >= len(l.src) {
		return 0
	}
	return l.src[l.pos+n]
}

func (l *lexer) advance() {
	if l.pos < len(l.src) && l.src[l.pos] == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

func (l *lexer) emit(k Kind, v string) Token {
	return Token{Kind: k, Val: v, Line: l.line, Col: l.col}
}

func shapeCloser(open string) string {
	switch open {
	case "([":
		return "])"
	case "((":
		return "))"
	case "[":
		return "]"
	case "(":
		return ")"
	case "{":
		return "}"
	}
	return ""
}

func isShapeOpen(r rune) bool { return r == '[' || r == '(' || r == '{' }

func isConnectorRune(r rune) bool {
	switch r {
	case '-', '.', '=', '>', '<', '~':
		return true
	}
	return false
}

func isIdentRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
