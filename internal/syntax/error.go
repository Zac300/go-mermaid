// Package syntax defines the positional error type shared by the lexer and
// parser. It is a dependency-free leaf so both stages (and the public
// package, via a type alias) can use one error type for source positions.
package syntax

import "fmt"

// Error reports a lexing or parsing failure with its source position.
type Error struct {
	Line int // 1-based line number
	Col  int // 1-based column (byte offset within the line)
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("line %d col %d: %s", e.Line, e.Col, e.Msg)
}

// Errorf builds an *Error with a formatted message.
func Errorf(line, col int, format string, args ...any) *Error {
	return &Error{Line: line, Col: col, Msg: fmt.Sprintf(format, args...)}
}
