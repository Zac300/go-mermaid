// Package svgutil holds small formatting helpers shared by the diagram
// renderers: deterministic number formatting and XML text escaping.
package svgutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Num formats a float for SVG output deterministically: rounded to two
// decimals with trailing zeros trimmed (e.g. 12, 12.5, 12.25).
func Num(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

var escaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&#39;",
)

// Esc escapes text for safe inclusion in SVG/XML content and attributes.
func Esc(s string) string { return escaper.Replace(s) }

// SplitLines splits label text on <br>, <br/>, <br /> and literal or escaped
// newlines, returning at least one line.
func SplitLines(s string) []string {
	for _, sep := range []string{"<br/>", "<br />", "<br>"} {
		s = strings.ReplaceAll(s, sep, "\n")
	}
	s = strings.ReplaceAll(s, "\\n", "\n")
	return strings.Split(s, "\n")
}

// MultilineText writes a centered <text> element whose lines are vertically
// centered around center. extraAttrs is inserted into the <text> tag (e.g.
// ` font-weight="bold"`).
func MultilineText(b *strings.Builder, lines []string, cx, center, lineH float64, fill, extraAttrs string) {
	n := len(lines)
	if n == 0 {
		return
	}
	fmt.Fprintf(b, `<text fill="%s" text-anchor="middle"%s>`, fill, extraAttrs)
	base := center - lineH*float64(n-1)/2
	for i, ln := range lines {
		fmt.Fprintf(b, `<tspan x="%s" y="%s">%s</tspan>`, Num(cx), Num(base+float64(i)*lineH), Esc(ln))
	}
	b.WriteString("</text>")
}

// TitleHeight is the vertical space reserved for a diagram title, or 0 when
// there is no title.
func TitleHeight(title string, fontSize float64) float64 {
	if title == "" {
		return 0
	}
	return fontSize*1.4 + 12
}
