package class

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Given a class with a member block", t, func() {
		src := "classDiagram\nclass Animal {\n+int age\n+isMammal() bool\n}"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then members are split into attributes and methods", func() {
				So(err, ShouldBeNil)
				c := d.class("Animal")
				So(c, ShouldNotBeNil)
				So(c.Attributes, ShouldResemble, []string{"+int age"})
				So(c.Methods, ShouldResemble, []string{"+isMammal() bool"})
			})
		})
	})

	Convey("Given shorthand member syntax", t, func() {
		Convey("When parsing", func() {
			d, err := Parse("classDiagram\nAnimal : +String name\nAnimal : +run() void")

			Convey("Then the class accumulates members", func() {
				So(err, ShouldBeNil)
				c := d.class("Animal")
				So(len(c.Attributes), ShouldEqual, 1)
				So(len(c.Methods), ShouldEqual, 1)
			})
		})
	})

	Convey("Given relationship operators", t, func() {
		cases := []struct {
			src    string
			left   headKind
			right  headKind
			dashed bool
		}{
			{"classDiagram\nA <|-- B", headTriangle, headNone, false},
			{"classDiagram\nA --|> B", headNone, headTriangle, false},
			{"classDiagram\nA *-- B", headDiamondFilled, headNone, false},
			{"classDiagram\nA o-- B", headDiamondHollow, headNone, false},
			{"classDiagram\nA --> B", headNone, headArrow, false},
			{"classDiagram\nA ..> B", headNone, headArrow, true},
			{"classDiagram\nA ..|> B", headNone, headTriangle, true},
		}
		for _, c := range cases {
			c := c
			Convey("When parsing "+c.src, func() {
				d, err := Parse(c.src)

				Convey("Then heads and line style are decoded", func() {
					So(err, ShouldBeNil)
					r := d.Relations[0]
					So(r.Left, ShouldEqual, c.left)
					So(r.Right, ShouldEqual, c.right)
					So(r.Dashed, ShouldEqual, c.dashed)
				})
			})
		}
	})

	Convey("Given a relationship with a label", t, func() {
		Convey("When parsing", func() {
			d, err := Parse("classDiagram\nOwner --> Dog : owns")

			Convey("Then the label is captured and both classes exist", func() {
				So(err, ShouldBeNil)
				So(d.Relations[0].Label, ShouldEqual, "owns")
				So(len(d.Classes), ShouldEqual, 2)
			})
		})
	})

	Convey("Given source without the header", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("class Animal")

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a class diagram with inheritance", t, func() {
		src := "classDiagram\nclass Animal {\n+int age\n}\nAnimal <|-- Dog"

		Convey("When rendering", func() {
			out, err := Render(src, RenderOptions{Theme: "default", FontFace: "sans-serif", FontSize: 14, Padding: 16})
			svg := string(out)

			Convey("Then it draws class boxes, members, and the relationship", func() {
				So(err, ShouldBeNil)
				So(svg, ShouldStartWith, "<svg")
				So(svg, ShouldContainSubstring, ">Animal<")
				So(svg, ShouldContainSubstring, ">Dog<")
				So(svg, ShouldContainSubstring, "+int age")
				So(svg, ShouldContainSubstring, "<path")
			})
		})
	})
}
