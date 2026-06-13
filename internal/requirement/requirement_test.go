package requirement

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a requirement diagram", t, func() {
		src := "requirementDiagram\nrequirement test_req {\nid: 1\ntext: the test\nrisk: high\n}\nelement test_entity {\ntype: simulation\n}\ntest_entity - satisfies -> test_req"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then requirement and element nodes parse with fields", func() {
				So(err, ShouldBeNil)
				So(len(d.Nodes), ShouldEqual, 2)
				req := d.node("test_req")
				So(req.Kind, ShouldEqual, "requirement")
				So(req.Fields["risk"], ShouldEqual, "high")
				So(d.node("test_entity").IsElement, ShouldBeTrue)
			})

			Convey("Then the relationship is typed", func() {
				So(len(d.Rels), ShouldEqual, 1)
				So(d.Rels[0].From, ShouldEqual, "test_entity")
				So(d.Rels[0].To, ShouldEqual, "test_req")
				So(d.Rels[0].Type, ShouldEqual, "satisfies")
			})
		})
	})

	Convey("Given no header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("requirement x {")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a requirement diagram, when rendering", t, func() {
		out, err := Render("requirementDiagram\nrequirement r {\nid: 1\n}\nelement e {\ntype: test\n}\ne - satisfies -> r",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws nodes and a typed relationship", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, ">r<")
			So(svg, ShouldContainSubstring, "«satisfies»")
		})
	})
}
