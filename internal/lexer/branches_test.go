package lexer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// firstText returns the first Text token value, or "" if none.
func firstText(toks []Token) string {
	for _, tk := range toks {
		if tk.Kind == Text {
			return tk.Val
		}
	}
	return ""
}

func TestShapeOpeners(t *testing.T) {
	Convey("Given each node shape syntax", t, func() {
		cases := map[string]string{
			"graph TD\nA[rect]":       "rect",
			"graph TD\nA(round)":      "round",
			"graph TD\nA([stad])":     "stad",
			"graph TD\nA((circ))":     "circ",
			"graph TD\nA{diamond}":    "diamond",
			"graph TD\nA[\"quoted\"]": "quoted",
		}
		for src, want := range cases {
			src, want := src, want
			Convey("When lexing "+src, func() {
				toks, err := Lex(src)

				Convey("Then the inner label text is captured", func() {
					So(err, ShouldBeNil)
					So(firstText(toks), ShouldEqual, want)
				})
			})
		}
	})
}

func TestUnterminated(t *testing.T) {
	Convey("Given malformed source", t, func() {
		cases := map[string]string{
			"shape at EOF":     "graph TD\nA[oops",
			"shape at newline": "graph TD\nA[oops\nB",
			"label at EOF":     "graph TD\nA -->|lab",
			"label at newline": "graph TD\nA -->|lab\nB",
			"bad character":    "graph TD\n#",
		}
		for name, src := range cases {
			src := src
			Convey("When lexing the "+name+" case", func() {
				_, err := Lex(src)

				Convey("Then it returns an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}

func TestPositionTracking(t *testing.T) {
	Convey("Given an identifier on the second line", t, func() {
		toks, err := Lex("graph TD\nAB --> C")

		Convey("When locating the AB token", func() {
			So(err, ShouldBeNil)
			var ab Token
			for _, tk := range toks {
				if tk.Kind == Ident && tk.Val == "AB" {
					ab = tk
					break
				}
			}

			Convey("Then its position is line 2, col 1", func() {
				So(ab.Line, ShouldEqual, 2)
				So(ab.Col, ShouldEqual, 1)
			})
		})
	})
}
