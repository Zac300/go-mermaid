package parser

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPreprocess(t *testing.T) {
	Convey("Given classDef, class, style, and inline directives", t, func() {
		src := strings.Join([]string{
			"graph TD",
			"A --> B",
			"C:::warn --> B",
			"classDef warn fill:#fdd,stroke:#f00,color:#900",
			"class A warn",
			"style B fill:#dfd,stroke:#393",
		}, "\n")

		Convey("When preprocessing", func() {
			cleaned, styles := Preprocess(src)

			Convey("Then styling directives are removed from the source", func() {
				So(cleaned, ShouldNotContainSubstring, "classDef")
				So(cleaned, ShouldNotContainSubstring, "style B")
				So(cleaned, ShouldNotContainSubstring, ":::")
				So(cleaned, ShouldContainSubstring, "A --> B")
				So(cleaned, ShouldContainSubstring, "C --> B")
			})

			Convey("Then class assignments resolve regardless of order", func() {
				So(styles["A"].Fill, ShouldEqual, "#fdd")
				So(styles["A"].Color, ShouldEqual, "#900")
				So(styles["C"].Stroke, ShouldEqual, "#f00") // via inline :::
			})

			Convey("Then direct style overrides apply", func() {
				So(styles["B"].Fill, ShouldEqual, "#dfd")
				So(styles["B"].Stroke, ShouldEqual, "#393")
			})
		})
	})

	Convey("Given source with no styling", t, func() {
		Convey("When preprocessing", func() {
			cleaned, styles := Preprocess("graph TD\nA --> B")

			Convey("Then the source is unchanged and styles are empty", func() {
				So(cleaned, ShouldEqual, "graph TD\nA --> B")
				So(len(styles), ShouldEqual, 0)
			})
		})
	})
}
