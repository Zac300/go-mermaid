package kanban

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a kanban board", t, func() {
		src := "kanban\n  Todo\n    [Task 1]\n    [Task 2]\n  Done\n    [Task 3]"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then columns and cards parse by indentation", func() {
				So(err, ShouldBeNil)
				So(len(d.Columns), ShouldEqual, 2)
				So(d.Columns[0].Title, ShouldEqual, "Todo")
				So(len(d.Columns[0].Cards), ShouldEqual, 2)
				So(d.Columns[0].Cards[0].Text, ShouldEqual, "Task 1")
				So(d.Columns[1].Title, ShouldEqual, "Done")
				So(len(d.Columns[1].Cards), ShouldEqual, 1)
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("  Todo")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a kanban board, when rendering", t, func() {
		out, err := Render("kanban\n  Todo\n    [A]\n  Done\n    [B]",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws column headers and cards", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, ">Todo<")
			So(svg, ShouldContainSubstring, ">A<")
		})
	})
}
