package mindmap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a mindmap with nested nodes", t, func() {
		src := "mindmap\n  root((Ideas))\n    Origins\n      History\n    Tools\n      Pen"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then the root and hierarchy are built", func() {
				So(err, ShouldBeNil)
				So(d.Root.Text, ShouldEqual, "Ideas") // shape markers stripped
				So(len(d.Root.Children), ShouldEqual, 2)
				So(d.Root.Children[0].Text, ShouldEqual, "Origins")
				So(d.Root.Children[0].Children[0].Text, ShouldEqual, "History")
			})

			Convey("Then depths increase with indentation", func() {
				So(d.Root.Depth, ShouldEqual, 0)
				So(d.Root.Children[0].Depth, ShouldEqual, 1)
				So(d.Root.Children[0].Children[0].Depth, ShouldEqual, 2)
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("  root")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a mindmap, when rendering", t, func() {
		out, err := Render("mindmap\n  root((R))\n    A\n    B",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws nodes and connector paths", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, ">R<")
			So(svg, ShouldContainSubstring, ">A<")
			So(svg, ShouldContainSubstring, "<path")
		})
	})
}
