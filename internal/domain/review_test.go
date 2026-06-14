package domain

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPolylineMidpoint(t *testing.T) {
	Convey("Given an L-shaped polyline", t, func() {
		pts := []Point{{X: 0, Y: 0}, {X: 0, Y: 10}, {X: 10, Y: 10}}

		Convey("When taking the midpoint", func() {
			mid := PolylineMidpoint(pts)

			Convey("Then it lands on the path at half the total length, not the chord", func() {
				So(mid.X, ShouldEqual, 0)
				So(mid.Y, ShouldEqual, 10)
			})
		})
	})

	Convey("Given a straight two-point segment", t, func() {
		mid := PolylineMidpoint([]Point{{X: 0, Y: 0}, {X: 10, Y: 0}})
		Convey("Then the midpoint is the centre", func() {
			So(mid.X, ShouldEqual, 5)
			So(mid.Y, ShouldEqual, 0)
		})
	})

	Convey("Given an empty polyline", t, func() {
		Convey("Then it returns the zero point without panicking", func() {
			So(PolylineMidpoint(nil), ShouldResemble, Point{})
		})
	})
}
