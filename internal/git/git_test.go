package git

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a gitGraph with a branch and merge", t, func() {
		src := "gitGraph\ncommit\ncommit id: \"a\"\nbranch develop\ncommit\ncheckout main\nmerge develop tag: \"v1\""

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then branches include main and develop", func() {
				So(err, ShouldBeNil)
				So(len(d.Branches), ShouldEqual, 2)
				So(d.Branches[0].Name, ShouldEqual, "main")
				So(d.Branches[1].Name, ShouldEqual, "develop")
			})

			Convey("Then commits and the merge are recorded", func() {
				So(len(d.Commits), ShouldEqual, 4) // 2 + 1 + merge
				So(len(d.Merges), ShouldEqual, 1)
				last := d.Commits[len(d.Commits)-1]
				So(last.Merge, ShouldBeTrue)
				So(last.Tag, ShouldEqual, "v1")
			})

			Convey("Then commit id metadata parses", func() {
				So(d.Commits[1].ID, ShouldEqual, "a")
			})
		})
	})

	Convey("Given the header has a trailing colon", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("gitGraph:\ncommit")

			Convey("Then it parses without error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("commit")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a gitGraph, when rendering", t, func() {
		out, err := Render("gitGraph\ncommit\nbranch dev\ncommit\ncheckout main\nmerge dev",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws lanes, commits, and branch labels", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "<circle")
			So(svg, ShouldContainSubstring, ">main<")
			So(svg, ShouldContainSubstring, ">dev<")
		})
	})
}
