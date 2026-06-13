package lexer

import "testing"

// shapeText lexes src and returns the first Text token value, or "" if none.
func shapeText(t *testing.T, src string) string {
	t.Helper()
	toks, err := Lex(src)
	if err != nil {
		t.Fatalf("Lex(%q): %v", src, err)
	}
	for _, tk := range toks {
		if tk.Kind == Text {
			return tk.Val
		}
	}
	return ""
}

func TestShapeOpeners(t *testing.T) {
	cases := map[string]string{
		"graph TD\nA[rect]":       "rect",
		"graph TD\nA(round)":      "round",
		"graph TD\nA([stad])":     "stad",
		"graph TD\nA((circ))":     "circ",
		"graph TD\nA{diamond}":    "diamond",
		"graph TD\nA[\"quoted\"]": "quoted",
	}
	for src, want := range cases {
		if got := shapeText(t, src); got != want {
			t.Errorf("Lex(%q) text = %q, want %q", src, got, want)
		}
	}
}

func TestUnterminated(t *testing.T) {
	for _, src := range []string{
		"graph TD\nA[oops",       // unterminated shape at EOF
		"graph TD\nA[oops\nB",    // shape broken by newline
		"graph TD\nA -->|lab",    // unterminated label at EOF
		"graph TD\nA -->|lab\nB", // label broken by newline
		"graph TD\n#",            // unexpected character
	} {
		if _, err := Lex(src); err == nil {
			t.Errorf("Lex(%q) expected error, got nil", src)
		}
	}
}

func TestPositionTracking(t *testing.T) {
	toks, err := Lex("graph TD\nAB --> C")
	if err != nil {
		t.Fatal(err)
	}
	// Second line's first ident "AB" should be at line 2, col 1.
	var ab Token
	for _, tk := range toks {
		if tk.Kind == Ident && tk.Val == "AB" {
			ab = tk
			break
		}
	}
	if ab.Line != 2 || ab.Col != 1 {
		t.Errorf("AB at line %d col %d, want line 2 col 1", ab.Line, ab.Col)
	}
}
