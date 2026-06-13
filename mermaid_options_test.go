package mermaid_test

import (
	"strconv"
	"strings"
	"testing"

	mermaid "github.com/Zac300/go-mermaid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	const src = "graph TD\nA --> B"

	Convey("Given Render options", t, func() {

		Convey("When setting font face and size", func() {
			out, err := mermaid.Render(src, mermaid.WithFont("Inter", 20))
			So(err, ShouldBeNil)
			So(string(out), ShouldContainSubstring, `font-family="Inter"`)
			So(string(out), ShouldContainSubstring, `font-size="20"`)
		})

		Convey("When setting padding", func() {
			out, err := mermaid.Render(src, mermaid.WithPadding(40))
			So(err, ShouldBeNil)
			So(string(out), ShouldContainSubstring, "translate(40,40)")
		})

		Convey("When increasing rank spacing", func() {
			tight, _ := mermaid.Render(src, mermaid.WithSpacing(50, 20))
			loose, _ := mermaid.Render(src, mermaid.WithSpacing(50, 200))
			So(height(loose), ShouldBeGreaterThan, height(tight))
		})

		Convey("When using the neutral theme", func() {
			out, err := mermaid.Render(src, mermaid.WithTheme(mermaid.Neutral))
			So(err, ShouldBeNil)
			So(string(out), ShouldContainSubstring, "#eeeeee")
		})
	})
}

// height extracts the SVG height attribute as a number.
func height(svg []byte) float64 {
	s := string(svg)
	i := strings.Index(s, `height="`)
	if i < 0 {
		return 0
	}
	s = s[i+len(`height="`):]
	v, _ := strconv.ParseFloat(s[:strings.IndexByte(s, '"')], 64)
	return v
}
