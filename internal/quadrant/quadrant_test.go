package quadrant

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a quadrant chart", t, func() {
		src := "quadrantChart\ntitle Reach\nx-axis Low --> High\ny-axis Bottom --> Top\nquadrant-1 Expand\nA: [0.3, 0.6]\nB: [0.45, 0.23]"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then title, axes, quadrant and points parse", func() {
				So(err, ShouldBeNil)
				So(d.Title, ShouldEqual, "Reach")
				So(d.XLeft, ShouldEqual, "Low")
				So(d.XRight, ShouldEqual, "High")
				So(d.YBottom, ShouldEqual, "Bottom")
				So(d.YTop, ShouldEqual, "Top")
				So(d.Quadrant[0], ShouldEqual, "Expand")
				So(len(d.Points), ShouldEqual, 2)
				So(d.Points[0].X, ShouldEqual, 0.3)
				So(d.Points[0].Y, ShouldEqual, 0.6)
			})
		})
	})

	Convey("Given an invalid point", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("quadrantChart\nA: [x, y]")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("title X")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a quadrant chart, when rendering", t, func() {
		out, err := Render("quadrantChart\ntitle T\nquadrant-1 Q1\nA: [0.7, 0.8]",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws quadrants, axes grid, and the point", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "<circle")
			So(svg, ShouldContainSubstring, ">A<")
			So(svg, ShouldContainSubstring, ">Q1<")
		})
	})
}
