// Command mermaid renders a Mermaid diagram to SVG.
//
// Usage:
//
//	mermaid [flags] [input.mmd]
//
// With no input file (or "-"), source is read from stdin. SVG is written
// to stdout unless -o is given.
//
//	mermaid diagram.mmd > diagram.svg
//	echo "graph TD; A-->B" | mermaid -theme dark -o out.svg
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	mermaid "github.com/Zac300/go-mermaid"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "mermaid:", err)
		os.Exit(1)
	}
}

func run() error {
	theme := flag.String("theme", "default", "color theme: default, dark, neutral")
	out := flag.String("o", "", "output file (default stdout)")
	padding := flag.Float64("padding", 16, "outer padding in pixels")
	flag.Parse()

	src, err := readInput(flag.Arg(0))
	if err != nil {
		return err
	}

	svg, err := mermaid.Render(string(src),
		mermaid.WithTheme(mermaid.Theme(*theme)),
		mermaid.WithPadding(*padding),
	)
	if err != nil {
		return err
	}

	if *out == "" || *out == "-" {
		_, err = os.Stdout.Write(svg)
		return err
	}
	return os.WriteFile(*out, svg, 0o644)
}

func readInput(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}
