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
			So(g.NodeByID("A"), ShouldEqual, a)
		})

		Convey("When looking up a missing node", func() {
			So(g.NodeByID("Z"), ShouldBeNil)
		})

		Convey("When computing a node center", func() {
			c := a.Center()
			So(c.X, ShouldEqual, 30) // 10 + 40/2
			So(c.Y, ShouldEqual, 35) // 20 + 30/2
		})
	})
}
