package state

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDescriptionContainingArrow(t *testing.T) {
	Convey("Given a state description whose text contains '-->'", t, func() {
		d, err := Parse("stateDiagram-v2\nIdle : waits for --> signal")

		Convey("When parsing", func() {
			Convey("Then it is a description, not a transition", func() {
				So(err, ShouldBeNil)
				So(len(d.Transitions), ShouldEqual, 0)
				var idle *State
				for _, s := range d.States {
					if s.ID == "Idle" {
						idle = s
					}
				}
				So(idle, ShouldNotBeNil)
				So(idle.Label, ShouldEqual, "waits for --> signal")
			})
		})

		Convey("When parsing a real labelled transition", func() {
			d2, err := Parse("stateDiagram-v2\nA --> B : go")
			Convey("Then the transition and its label are captured", func() {
				So(err, ShouldBeNil)
				So(len(d2.Transitions), ShouldEqual, 1)
				So(d2.Transitions[0].Label, ShouldEqual, "go")
			})
		})
	})
}
