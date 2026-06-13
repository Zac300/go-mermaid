package mermaid_test

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mermaid "github.com/Zac300/go-mermaid"
	. "github.com/smartystreets/goconvey/convey"
)

// update regenerates golden files: go test -run TestGolden -update
var update = flag.Bool("update", false, "update golden SVG files")

func TestRender(t *testing.T) {
	Convey("Given the public Render API", t, func() {

		Convey("When rendering a minimal flowchart", func() {
			out, err := mermaid.Render("graph TD\nA --> B")

			Convey("Then it returns valid-looking SVG bytes", func() {
				So(err, ShouldBeNil)
				So(string(out), ShouldStartWith, "<svg")
				So(string(out), ShouldContainSubstring, "</svg>")
			})
		})

		Convey("When options are supplied", func() {
			out, err := mermaid.Render("graph TD\nA --> B", mermaid.WithTheme(mermaid.Dark))

			Convey("Then the dark background is applied", func() {
				So(err, ShouldBeNil)
				So(string(out), ShouldContainSubstring, "#1e1e1e")
			})
		})

		Convey("When the source has a syntax error", func() {
			_, err := mermaid.Render("not a diagram {{{")

			Convey("Then ErrParse matches and a ParseError is recoverable", func() {
				So(errors.Is(err, mermaid.ErrParse), ShouldBeTrue)
				var pe *mermaid.ParseError
				So(errors.As(err, &pe), ShouldBeTrue)
			})
		})
	})
}

func TestGolden(t *testing.T) {
	dir := filepath.Join("testdata", "golden")
	inputs, err := filepath.Glob(filepath.Join(dir, "*.mmd"))
	if err != nil {
		t.Fatal(err)
	}
	if len(inputs) == 0 {
		t.Skip("no golden inputs")
	}

	for _, in := range inputs {
		in := in
		name := strings.TrimSuffix(filepath.Base(in), ".mmd")
		t.Run(name, func(t *testing.T) {
			src, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}
			got, err := mermaid.Render(string(src))
			if err != nil {
				t.Fatalf("Render(%s): %v", name, err)
			}
			goldenPath := strings.TrimSuffix(in, ".mmd") + ".svg"
			if *update {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden (run with -update to create): %v", err)
			}
			if string(got) != string(want) {
				t.Errorf("%s: output differs from golden; run: go test -run TestGolden -update", name)
			}
		})
	}
}
