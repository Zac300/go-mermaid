# go-mermaid

[![CI](https://github.com/Zac300/go-mermaid/actions/workflows/ci.yml/badge.svg)](https://github.com/Zac300/go-mermaid/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Zac300/go-mermaid/branch/main/graph/badge.svg)](https://codecov.io/gh/Zac300/go-mermaid)
[![Go Reference](https://pkg.go.dev/badge/github.com/Zac300/go-mermaid.svg)](https://pkg.go.dev/github.com/Zac300/go-mermaid)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zac300/go-mermaid)](https://goreportcard.com/report/github.com/Zac300/go-mermaid)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Render [Mermaid](https://mermaid.js.org) diagrams to SVG in **pure Go** — no
headless browser, no Node.js, no JavaScript runtime. Just a library and a
single static binary.

> **Status:** early (v0). Flowcharts (`graph TD` / `graph LR`) are supported.
> Other diagram types are on the [roadmap](#roadmap). Not affiliated with the
> Mermaid project; this is an independent, compatible renderer.

## Why

Every existing Go path to Mermaid SVG shells out to headless Chrome or a Node
sidecar. That's heavy, slow, and hard to deploy. `go-mermaid` does the parse →
layout → SVG pipeline natively, so you can render diagrams inside any Go service
or CLI with zero external dependencies.

## Install

Library:

```bash
go get github.com/Zac300/go-mermaid
```

CLI:

```bash
go install github.com/Zac300/go-mermaid/cmd/mermaid@latest
```

## Usage

### Library

```go
package main

import (
	"os"

	mermaid "github.com/Zac300/go-mermaid"
)

func main() {
	svg, err := mermaid.Render("graph TD\n  A[Start] --> B{OK?}\n  B -->|yes| C([Done])")
	if err != nil {
		panic(err)
	}
	os.WriteFile("diagram.svg", svg, 0o644)
}
```

With options:

```go
svg, err := mermaid.Render(src,
	mermaid.WithTheme(mermaid.Dark),
	mermaid.WithFont("Inter", 14),
	mermaid.WithPadding(24),
	mermaid.WithSpacing(50, 60),
)
```

### CLI

```bash
mermaid diagram.mmd > diagram.svg
echo "graph LR; A-->B-->C" | mermaid -theme dark -o out.svg
```

## Error handling

`Render` wraps stage-specific sentinels so you can branch on the failure, and
parse errors carry source position:

```go
svg, err := mermaid.Render(src)
if errors.Is(err, mermaid.ErrParse) {
	var pe *mermaid.ParseError
	if errors.As(err, &pe) {
		log.Printf("syntax error at line %d col %d: %s", pe.Line, pe.Col, pe.Msg)
	}
}
```

Sentinels: `ErrParse`, `ErrLayout`, `ErrRender`, `ErrUnsupported`.

## Supported syntax (v0)

| Feature | Example |
| --- | --- |
| Directions | `graph TD`, `TB`, `BT`, `LR`, `RL` |
| Rectangle | `A[Label]` |
| Rounded | `A(Label)` |
| Stadium | `A([Label])` |
| Circle | `A((Label))` |
| Diamond | `A{Label}` |
| Arrow | `A --> B` |
| Open link | `A --- B` |
| Dotted | `A -.-> B` |
| Thick | `A ==> B` |
| Edge label | `A -->\|text\| B` |
| Comments | `%% comment` |

## Roadmap

- [ ] Network-simplex ranking (tighter layouts)
- [ ] Crossing minimization (median/barycenter ordering)
- [ ] Orthogonal/spline edge routing
- [ ] Subgraphs
- [ ] Sequence diagrams
- [ ] Class diagrams
- [ ] PNG output

## Architecture

A small, testable pipeline:

```
source → lexer → parser → domain.Graph → layout → render → SVG
```

The `internal/domain` package holds the pure model; each stage is an adapter
that produces or consumes it. See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE) © Zac Krebbekx
