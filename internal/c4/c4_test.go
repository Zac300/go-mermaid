package c4

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a C4 context diagram", t, func() {
		src := "C4Context\ntitle System Context\nPerson(custA, \"Customer\", \"A bank customer\")\nSystem(sysA, \"Banking\", \"Lets customers view accounts\")\nRel(custA, sysA, \"Uses\")"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then title, elements, and relationships parse", func() {
				So(err, ShouldBeNil)
				So(d.Title, ShouldEqual, "System Context")
				So(len(d.Elements), ShouldEqual, 2)
				So(d.element("custA").Kind, ShouldEqual, "Person")
				So(d.element("custA").Label, ShouldEqual, "Customer")
				So(d.element("custA").Descr, ShouldEqual, "A bank customer")
				So(len(d.Rels), ShouldEqual, 1)
				So(d.Rels[0].From, ShouldEqual, "custA")
				So(d.Rels[0].Label, ShouldEqual, "Uses")
			})
		})
	})

	Convey("Given a non-C4 header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("graph TD")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a C4 diagram, when rendering", t, func() {
		out, err := Render("C4Context\nPerson(a, \"A\", \"desc\")\nSystem(b, \"B\", \"desc\")\nRel(a, b, \"uses\")",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws element boxes with C4 colors and a relationship", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "#08427b") // person
			So(svg, ShouldContainSubstring, "#1168bd") // system
			So(svg, ShouldContainSubstring, ">uses<")
		})
	})
}
