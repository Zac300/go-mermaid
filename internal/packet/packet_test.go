package packet

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a packet diagram", t, func() {
		src := "packet-beta\n0-15: \"Source Port\"\n16-31: \"Dest Port\"\n32: \"flag\""

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then bit ranges and labels parse", func() {
				So(err, ShouldBeNil)
				So(len(d.Fields), ShouldEqual, 3)
				So(d.Fields[0].Start, ShouldEqual, 0)
				So(d.Fields[0].End, ShouldEqual, 15)
				So(d.Fields[0].Label, ShouldEqual, "Source Port")
				So(d.Fields[2].Start, ShouldEqual, 32)
				So(d.Fields[2].End, ShouldEqual, 32)
			})
		})
	})

	Convey("Given an invalid range", t, func() {
		cases := map[string]string{
			"bad bits":  "packet-beta\n0-x: \"a\"",
			"reversed":  "packet-beta\n10-2: \"a\"",
			"no colon":  "packet-beta\n0-15 label",
			"no header": "0-15: \"a\"",
		}
		for name, src := range cases {
			src := src
			Convey("When parsing the "+name+" case", func() {
				_, err := Parse(src)
				Convey("Then it returns an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}

func TestRender(t *testing.T) {
	Convey("Given a packet diagram that wraps rows, when rendering", t, func() {
		out, err := Render("packet-beta\n0-31: \"Header\"\n32-63: \"Payload\"",
			RenderOptions{Theme: "default", FontSize: 14, Padding: 16})
		svg := string(out)

		Convey("Then it draws field cells and bit labels", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, ">Header<")
			So(svg, ShouldContainSubstring, ">Payload<")
		})
	})
}
