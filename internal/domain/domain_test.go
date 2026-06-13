package domain

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphAndNode(t *testing.T) {
	Convey("Given a graph with two nodes", t, func() {
		a := &Node{ID: "A", Pos: Point{X: 10, Y: 20}, Size: Size{W: 40, H: 30}}
		b := &Node{ID: "B"}
		g := &Graph{Nodes: []*Node{a, b}}

		Convey("When looking up an existing node", func() {
			got := g.NodeByID("A")

			Convey("Then it returns that node", func() {
				So(got, ShouldEqual, a)
			})
		})

		Convey("When looking up a missing node", func() {
			got := g.NodeByID("Z")

			Convey("Then it returns nil", func() {
				So(got, ShouldBeNil)
			})
		})

		Convey("When computing a node center", func() {
			c := a.Center()

			Convey("Then it is the midpoint of the box", func() {
				So(c.X, ShouldEqual, 30) // 10 + 40/2
				So(c.Y, ShouldEqual, 35) // 20 + 30/2
			})
		})
	})
}
