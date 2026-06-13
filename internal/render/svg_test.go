package render

import (
	"strings"
	"testing"

	"github.com/Zac300/go-mermaid/internal/layout"
	"github.com/Zac300/go-mermaid/internal/lexer"
	"github.com/Zac300/go-mermaid/internal/parser"
	. "github.com/smartystreets/goconvey/convey"
)

func laidOut(src string) *layout.Result {
	toks, err := lexer.Lex(src)
	if err != nil {
		panic(err)
	}
	g, err := parser.Parse(toks)
	if err != nil {
		panic(err)
	}
	res, err := layout.Compute(g, layout.Options{NodeSep: 50, RankSep: 50, FontSize: 14})
	if err != nil {
		panic(err)
	}
	return res
}

var opts = Options{Theme: "default", FontFace: "sans-serif", FontSize: 14, Padding: 16}

func TestSVG(t *testing.T) {
	Convey("Given a laid-out flowchart", t, func() {
		res := laidOut("graph TD\nA[Start] --> B((End))")
		out, err := SVG(res, opts)
		svg := string(out)

		Convey("Then it produces a well-formed SVG document", func() {
			So(err, ShouldBeNil)
			So(svg, ShouldStartWith, "<svg")
			So(svg, ShouldContainSubstring, "</svg>")
			So(svg, ShouldContainSubstring, "marker id=\"arrow\"")
			So(svg, ShouldContainSubstring, "<rect")   // rectangle node
			So(svg, ShouldContainSubstring, "<circle") // circle node
			So(svg, ShouldContainSubstring, "<path")   // edge
		})
	})

	Convey("Given each node shape", t, func() {
		cases := map[string]string{
			"round":   "graph TD\nA(R)",
			"stadium": "graph TD\nA([S])",
			"diamond": "graph TD\nA{D}",
		}
		wants := map[string]string{
			"round":   "<rect",
			"stadium": "<rect",
			"diamond": "<polygon",
		}
		for name, src := range cases {
			out, _ := SVG(laidOut(src), opts)
			So(string(out), ShouldContainSubstring, wants[name])
		}
	})

	Convey("Given different arrow styles", t, func() {
		dotted, _ := SVG(laidOut("graph TD\nA -.-> B"), opts)
		So(string(dotted), ShouldContainSubstring, "stroke-dasharray")

		thick, _ := SVG(laidOut("graph TD\nA ==> B"), opts)
		So(string(thick), ShouldContainSubstring, `stroke-width="3"`)

		open, _ := SVG(laidOut("graph TD\nA --- B"), opts)
		So(string(open), ShouldNotContainSubstring, "marker-end")
	})

	Convey("Given an edge label", t, func() {
		out, _ := SVG(laidOut("graph TD\nA -->|go| B"), opts)
		So(string(out), ShouldContainSubstring, ">go<")
	})

	Convey("Given the dark theme", t, func() {
		out, _ := SVG(laidOut("graph TD\nA --> B"), Options{Theme: "dark", FontSize: 14, Padding: 16})
		So(string(out), ShouldContainSubstring, "#1e1e1e")
	})

	Convey("Given an unknown theme", t, func() {
		out, _ := SVG(laidOut("graph TD\nA --> B"), Options{Theme: "nope", FontSize: 14})
		So(string(out), ShouldContainSubstring, "#ffffff") // falls back to default
	})

	Convey("Given a label with XML-special characters", t, func() {
		res := laidOut("graph TD\nA --> B")
		res.Graph.Nodes[0].Label = `a<b & "c"`
		out, _ := SVG(res, opts)
		So(string(out), ShouldContainSubstring, "a&lt;b &amp; &quot;c&quot;")
	})
}

func TestNum(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{12, "12"},
		{12.5, "12.5"},
		{12.25, "12.25"},
		{0, "0"},
		{-0.0, "0"},
		{1.200, "1.2"},
	}
	for _, c := range cases {
		if got := num(c.in); got != c.want {
			t.Errorf("num(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEscPlain(t *testing.T) {
	if got := esc("plain"); strings.Contains(got, "&") {
		t.Errorf("esc(plain) altered text: %q", got)
	}
}
