package xychart

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given an xychart with bar and line series", t, func() {
		src := "xychart-beta\ntitle \"Sales\"\nx-axis [jan, feb, mar]\ny-axis \"Rev\" 0 --> 10000\nbar [5000, 6000, 7500]\nline [4000, 6000, 9000]"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then title, axes, and series parse", func() {
				So(err, ShouldBeNil)
				So(d.Title, ShouldEqual, "Sales")
				So(d.XCats, ShouldResemble, []string{"jan", "feb", "mar"})
				So(d.YMin, ShouldEqual, 0)
				So(d.YMax, ShouldEqual, 10000)
				So(d.HasYRange, ShouldBeTrue)
				So(len(d.Series), ShouldEqual, 2)
				So(d.Series[0].Kind, ShouldEqual, "bar")
				So(d.Series[1].Values[2], ShouldEqual, 9000)
			})
		})
	})

	Convey("Given no y range", t, func() {
		Convey("When parsing and computing bounds", func() {
			d, err := Parse("xychart-beta\nbar [1, 2, 8]")

			Convey("Then bounds derive 0..max", func() {
				So(err, ShouldBeNil)
				lo, hi := d.Bounds()
				So(lo, ShouldEqual, 0)
				So(hi, ShouldEqual, 8)
			})
		})
	})

	Convey("Given no series", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("xychart-beta\ntitle \"x\"")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given an xychart, when rendering", t, func() {
		out, err := Render("xychart-beta\ntitle \"S\"\nx-axis [a, b]\ny-axis 0 --> 10\nbar [5, 8]\nline [3, 9]",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws bars, a line, and axis labels", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "<rect")
			So(svg, ShouldContainSubstring, "<path")
			So(svg, ShouldContainSubstring, ">a<")
		})
	})
}
