package radar

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a radar chart with axes and curves", t, func() {
		src := "radar-beta\ntitle Skills\naxis a[\"Speed\"], b[\"Power\"], c[\"Range\"]\ncurve s1[\"Team A\"]{80, 60, 90}\ncurve s2[\"Team B\"]{50, 90, 40}"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then title, axes, and curves parse with labels", func() {
				So(err, ShouldBeNil)
				So(d.Title, ShouldEqual, "Skills")
				So(d.Axes, ShouldResemble, []string{"Speed", "Power", "Range"})
				So(len(d.Curves), ShouldEqual, 2)
				So(d.Curves[0].Name, ShouldEqual, "Team A")
				So(d.Curves[0].Values, ShouldResemble, []float64{80, 60, 90})
				So(d.Max(), ShouldEqual, 90)
			})
		})
	})

	Convey("Given no axes", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("radar-beta\ntitle X")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("axis a, b")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a radar chart, when rendering", t, func() {
		out, err := Render("radar-beta\naxis a, b, c\ncurve s1[\"A\"]{1, 2, 3}",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws grid rings, axis labels, and a curve polygon", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, ">a<")
			So(svg, ShouldContainSubstring, "fill-opacity=\"0.2\"")
		})
	})
}
