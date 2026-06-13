package render

import (
	"strconv"
	"strings"
)

// num formats a float for SVG output deterministically: rounded to two
// decimals with trailing zeros trimmed (e.g. 12, 12.5, 12.25).
func num(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

// esc escapes text for safe inclusion in SVG/XML content and attributes.
var escaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&#39;",
)

func esc(s string) string { return escaper.Replace(s) }
