package parser

import (
	"testing"

	"github.com/Zac300/go-mermaid/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDirections(t *testing.T) {
	Convey("Given each direction keyword", t, func() {
		cases := map[string]domain.Direction{
			"graph TD\nA": domain.TopBottom,
			"graph TB\nA": domain.TopBottom,
			"graph BT\nA": domain.BottomTop,
			"graph LR\nA": domain.LeftRight,
			"graph RL\nA": domain.RightLeft,
		}
		for src, want := range cases {
			src, want := src, want
			Convey("When parsing "+src, func() {
				g, err := parse(src)

				Convey("Then the graph direction matches", func() {
					So(err, ShouldBeNil)
					So(g.Direction, ShouldEqual, want)
				})
			})
		}
	})
}

func TestArrowKinds(t *testing.T) {
	Convey("Given each arrow style", t, func() {
		cases := map[string]domain.Arrow{
			"graph TD\nA --> B":  domain.ArrowNormal,
			"graph TD\nA --- B":  domain.ArrowOpen,
			"graph TD\nA -.-> B": domain.ArrowDotted,
			"graph TD\nA ==> B":  domain.ArrowThick,
		}
		for src, want := range cases {
			src, want := src, want
			Convey("When parsing "+src, func() {
				g, err := parse(src)

				Convey("Then the edge arrow kind matches", func() {
					So(err, ShouldBeNil)
					So(g.Edges[0].Arrow, ShouldEqual, want)
				})
			})
		}
	})
}

func TestParseErrors(t *testing.T) {
	Convey("Given invalid source", t, func() {
		cases := map[string]string{
			"unknown direction":  "graph XY\nA --> B",
			"junk after header":  "graph TD A --> B",
			"missing node id":    "graph TD\n--> B",
			"missing target":     "graph TD\nA -->",
			"unterminated label": "graph TD\nA -->|x B",
		}
		for name, src := range cases {
			src := src
			Convey("When parsing the "+name+" case", func() {
				_, err := parse(src)

				Convey("Then it returns an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}

func TestFlowchartKeyword(t *testing.T) {
	Convey("Given the 'flowchart' keyword", t, func() {
		g, err := parse("flowchart LR\nA --> B")

		Convey("When parsed", func() {
			Convey("Then it behaves like 'graph' with a direction", func() {
				So(err, ShouldBeNil)
				So(g.Direction, ShouldEqual, domain.LeftRight)
			})
		})
	})
}
