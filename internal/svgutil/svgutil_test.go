package svgutil

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitLines(t *testing.T) {
	Convey("Given label text with line breaks", t, func() {
		cases := map[string][]string{
			"one<br>two":     {"one", "two"},
			"a<br/>b<br />c": {"a", "b", "c"},
			"x\\ny":          {"x", "y"},
			"single":         {"single"},
		}
		for in, want := range cases {
			in, want := in, want
			Convey("When splitting "+in, func() {
				Convey("Then the lines match", func() {
					So(SplitLines(in), ShouldResemble, want)
				})
			})
		}
	})
}

func TestMultilineText(t *testing.T) {
	Convey("Given two lines", t, func() {
		var b strings.Builder
		MultilineText(&b, []string{"a", "b"}, 50, 30, 16, "#000", ` font-weight="bold"`)
		out := b.String()

		Convey("When rendered", func() {
			Convey("Then it produces a text element with two tspans", func() {
				So(out, ShouldStartWith, "<text")
				So(strings.Count(out, "<tspan"), ShouldEqual, 2)
				So(out, ShouldContainSubstring, `font-weight="bold"`)
			})
		})
	})
}

func TestTextWidth(t *testing.T) {
	Convey("Given text at a font size", t, func() {
		Convey("When measuring", func() {
			Convey("Then wide glyphs are wider than narrow ones", func() {
				So(TextWidth("WWWW", 14), ShouldBeGreaterThan, TextWidth("iiii", 14))
			})
			Convey("Then width scales with font size", func() {
				So(TextWidth("hello", 28), ShouldAlmostEqual, TextWidth("hello", 14)*2, 0.01)
			})
			Convey("Then empty text has zero width", func() {
				So(TextWidth("", 14), ShouldEqual, 0)
			})
			Convey("Then non-ASCII falls back to an estimate", func() {
				So(TextWidth("日本語", 14), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestNum(t *testing.T) {
	Convey("Given numbers", t, func() {
		Convey("When formatted", func() {
			Convey("Then trailing zeros are trimmed", func() {
				So(Num(1.200), ShouldEqual, "1.2")
				So(Num(3), ShouldEqual, "3")
			})
		})
	})
}
