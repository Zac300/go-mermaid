package mermaid

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExtractA11y(t *testing.T) {
	Convey("Given body with accTitle and accDescr", t, func() {
		body := "graph TD\naccTitle: The Title\naccDescr: The description\nA-->B"

		Convey("When extracting accessibility directives", func() {
			at, ad, cleaned := extractA11y(body)

			Convey("Then they are captured and removed from the body", func() {
				So(at, ShouldEqual, "The Title")
				So(ad, ShouldEqual, "The description")
				So(cleaned, ShouldNotContainSubstring, "accTitle")
				So(cleaned, ShouldContainSubstring, "A-->B")
			})
		})
	})
}

func TestInjectA11y(t *testing.T) {
	Convey("Given an SVG document", t, func() {
		svg := []byte(`<svg xmlns="x" width="10" height="10">` + "\n" + `<rect/>` + "\n" + `</svg>`)

		Convey("When injecting a title and description", func() {
			out := string(injectA11y(svg, "T", "D"))

			Convey("Then role, title, and desc are present", func() {
				So(out, ShouldContainSubstring, `role="img"`)
				So(out, ShouldContainSubstring, "<title>T</title>")
				So(out, ShouldContainSubstring, "<desc>D</desc>")
			})
		})

		Convey("When there is no title or desc", func() {
			out := string(injectA11y(svg, "", ""))

			Convey("Then only role is added", func() {
				So(out, ShouldContainSubstring, `role="img"`)
				So(out, ShouldNotContainSubstring, "<title>")
			})
		})
	})
}
