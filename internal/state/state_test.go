package state

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a state diagram with start, transitions, and end", t, func() {
		src := "stateDiagram-v2\n[*] --> Still\nStill --> Moving : go\nMoving --> [*]"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then pseudostates and transitions are captured", func() {
				So(err, ShouldBeNil)
				So(d.state(startID), ShouldNotBeNil)
				So(d.state(startID).Start, ShouldBeTrue)
				So(d.state(endID).End, ShouldBeTrue)
				So(len(d.Transitions), ShouldEqual, 3)
			})

			Convey("Then transition labels are captured", func() {
				So(d.Transitions[1].Label, ShouldEqual, "go")
				So(d.Transitions[1].From, ShouldEqual, "Still")
				So(d.Transitions[1].To, ShouldEqual, "Moving")
			})
		})
	})

	Convey("Given a state description", t, func() {
		Convey("When parsing", func() {
			d, err := Parse("stateDiagram-v2\nStill : The idle state\nStill --> [*]")

			Convey("Then the description becomes the state's label", func() {
				So(err, ShouldBeNil)
				So(d.state("Still").Label, ShouldEqual, "The idle state")
			})
		})
	})

	Convey("Given a composite state block", t, func() {
		Convey("When parsing", func() {
			d, err := Parse("stateDiagram-v2\nstate Active {\n[*] --> Idle\n}\nActive --> [*]")

			Convey("Then the composite state exists and the body is skipped", func() {
				So(err, ShouldBeNil)
				So(d.state("Active"), ShouldNotBeNil)
			})
		})
	})

	Convey("Given source without the header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("[*] --> A")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a state diagram", t, func() {
		Convey("When rendering", func() {
			out, err := Render("stateDiagram-v2\n[*] --> Still\nStill --> [*]",
				RenderOptions{Theme: "default", FontFace: "sans-serif", FontSize: 14, Padding: 16})
			svg := string(out)

			Convey("Then it draws pseudostate circles, a state box, and arrows", func() {
				So(err, ShouldBeNil)
				So(svg, ShouldStartWith, "<svg")
				So(svg, ShouldContainSubstring, "<circle")
				So(svg, ShouldContainSubstring, ">Still<")
				So(svg, ShouldContainSubstring, "marker-end=\"url(#st-arrow)\"")
			})
		})
	})
}
