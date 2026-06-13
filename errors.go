package mermaid

import "errors"

// Sentinel errors returned by Render, wrapped around the underlying cause.
// Match them with errors.Is. For parse failures, errors.As can additionally
// recover a *parser.ParseError carrying source line/column.
var (
	// ErrParse indicates the source could not be lexed or parsed.
	ErrParse = errors.New("mermaid: parse error")
	// ErrLayout indicates the graph could not be laid out.
	ErrLayout = errors.New("mermaid: layout error")
	// ErrRender indicates the laid-out graph could not be rendered to SVG.
	ErrRender = errors.New("mermaid: render error")
	// ErrUnsupported indicates a diagram type or feature not yet implemented.
	ErrUnsupported = errors.New("mermaid: unsupported feature")
)
