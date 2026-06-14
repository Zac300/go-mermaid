package layout

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelfLoopRanking(t *testing.T) {
	Convey("Given a node with a self-loop that also has a successor", t, func() {
		g := graphFrom("graph TD\nA --> A\nA --> B")

		Convey("When computing the layout", func() {
			res, err := Compute(g, opts)

			Convey("Then the self-loop does not collapse ranks and B sits below A", func() {
				So(err, ShouldBeNil)
				a := g.NodeByID("A")
				b := g.NodeByID("B")
				So(a, ShouldNotBeNil)
				So(b, ShouldNotBeNil)
				So(b.Pos.Y, ShouldBeGreaterThan, a.Pos.Y)
				So(res.Height, ShouldBeGreaterThan, 0)
			})
		})
	})
}
