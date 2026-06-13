package journey

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a journey with sections and tasks", t, func() {
		src := "journey\ntitle My day\nsection Work\nMake tea: 5: Me\nDo work: 1: Me, Cat"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then title, sections, and tasks are captured", func() {
				So(err, ShouldBeNil)
				So(d.Title, ShouldEqual, "My day")
				So(len(d.Sections), ShouldEqual, 1)
				So(len(d.Sections[0].Tasks), ShouldEqual, 2)
			})

			Convey("Then task score and actors parse", func() {
				task := d.Sections[0].Tasks[1]
				So(task.Name, ShouldEqual, "Do work")
				So(task.Score, ShouldEqual, 1)
				So(task.Actors, ShouldResemble, []string{"Me", "Cat"})
			})
		})
	})

	Convey("Given an invalid task line", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("journey\nsection S\nBad task")

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
	Convey("Given a journey, when rendering", t, func() {
		out, err := Render("journey\ntitle Day\nsection Work\nTea: 5: Me\nWork: 1: Me",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws score points and labels", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "<circle")
			So(svg, ShouldContainSubstring, ">Tea<")
			So(svg, ShouldContainSubstring, "#27ae60") // score 5 green
		})
	})
}
