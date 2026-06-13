package mermaid

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseFrontmatter(t *testing.T) {
	Convey("Given source with a title front-matter block", t, func() {
		src := "---\ntitle: My Diagram\n---\ngraph TD\nA-->B"

		Convey("When parsing front-matter", func() {
			title, body := parseFrontmatter(src)

			Convey("Then the title is extracted and stripped from the body", func() {
				So(title, ShouldEqual, "My Diagram")
				So(body, ShouldStartWith, "graph TD")
				So(body, ShouldNotContainSubstring, "title:")
			})
		})
	})

	Convey("Given source without front-matter", t, func() {
		Convey("When parsing", func() {
			title, body := parseFrontmatter("graph TD\nA-->B")

			Convey("Then the title is empty and body is unchanged", func() {
				So(title, ShouldEqual, "")
				So(body, ShouldEqual, "graph TD\nA-->B")
			})
		})
	})

	Convey("Given a quoted title", t, func() {
		Convey("When parsing", func() {
			title, _ := parseFrontmatter("---\ntitle: \"Quoted Title\"\n---\ngraph TD\nA-->B")

			Convey("Then surrounding quotes are removed", func() {
				So(title, ShouldEqual, "Quoted Title")
			})
		})
	})

	Convey("Given unterminated front-matter", t, func() {
		Convey("When parsing", func() {
			title, body := parseFrontmatter("---\ntitle: X\ngraph TD")

			Convey("Then the whole source is treated as body", func() {
				So(title, ShouldEqual, "")
				So(body, ShouldContainSubstring, "---")
			})
		})
	})
}
