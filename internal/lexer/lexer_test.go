package lexer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func kinds(toks []Token) []Kind {
	ks := make([]Kind, 0, len(toks))
	for _, t := range toks {
		ks = append(ks, t.Kind)
	}
	return ks
}

func TestLex(t *testing.T) {
	Convey("Given Mermaid flowchart source", t, func() {

		Convey("When lexing a header and a simple edge", func() {
			toks, err := Lex("graph TD\nA --> B")

			Convey("Then it produces the expected token kinds", func() {
				So(err, ShouldBeNil)
				So(kinds(toks), ShouldResemble, []Kind{
					Keyword, Ident, // graph TD
					Newline,
					Ident, Arrow, Ident, // A --> B
					EOF,
				})
			})
		})

		Convey("When lexing shapes and an inline edge label", func() {
			toks, err := Lex("graph LR\nA[Start] -->|go| B((End))")

			Convey("Then shape and pipe text are captured as Text tokens", func() {
				So(err, ShouldBeNil)
				var texts []string
				for _, tk := range toks {
					if tk.Kind == Text {
						texts = append(texts, tk.Val)
					}
				}
				So(texts, ShouldResemble, []string{"Start", "go", "End"})
			})
		})

		Convey("When a shape is left unterminated", func() {
			_, err := Lex("graph TD\nA[Start --> B")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestLexTable(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Kind
	}{
		{"bare nodes", "graph TD\nA", []Kind{Keyword, Ident, Newline, Ident, EOF}},
		{"dotted arrow", "graph TD\nA -.-> B", []Kind{Keyword, Ident, Newline, Ident, Arrow, Ident, EOF}},
		{"comment skipped", "graph TD\n%% note\nA --> B", []Kind{Keyword, Ident, Newline, Newline, Ident, Arrow, Ident, EOF}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			toks, err := Lex(c.in)
			if err != nil {
				t.Fatalf("Lex(%q) error: %v", c.in, err)
			}
			got := kinds(toks)
			if len(got) != len(c.want) {
				t.Fatalf("Lex(%q) = %v, want %v", c.in, got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("Lex(%q)[%d] = %v, want %v", c.in, i, got[i], c.want[i])
				}
			}
		})
	}
}
