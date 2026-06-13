// Package render turns a laid-out graph into SVG bytes. Output is
// deterministic so it can be compared against golden files.
package render

import (
	"fmt"
	"strings"

	"github.com/Zac300/go-mermaid/internal/domain"
	"github.com/Zac300/go-mermaid/internal/layout"
	"github.com/Zac300/go-mermaid/internal/svgutil"
	"github.com/Zac300/go-mermaid/internal/theme"
)

// Options controls SVG appearance.
type Options struct {
	Theme    string
	FontFace string
	FontSize float64
	Padding  float64
	Title    string
}

// writeTitle draws a centered, bold diagram title at (x, y) if non-empty.
func writeTitle(b *strings.Builder, title string, x, y float64, pal theme.Palette) {
	if title == "" {
		return
	}
	fmt.Fprintf(b, `  <text x="%s" y="%s" fill="%s" text-anchor="middle" font-weight="bold">%s</text>`,
		num(x), num(y), pal.Text, esc(title))
	b.WriteByte('\n')
}

// SVG renders a laid-out graph to an SVG document.
func SVG(res *layout.Result, opts Options) ([]byte, error) {
	pal := theme.For(opts.Theme)
	pad := opts.Padding
	titleH := svgutil.TitleHeight(opts.Title, opts.FontSize)

	contentW := res.Width
	if tw := opts.FontSize * 0.6 * float64(len([]rune(opts.Title))); tw > contentW {
		contentW = tw
	}
	w := contentW + pad*2
	h := res.Height + titleH + pad*2

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 %s %s" font-family="%s" font-size="%s">`,
		num(w), num(h), num(w), num(h), esc(opts.FontFace), num(opts.FontSize))
	b.WriteByte('\n')

	// Arrowhead marker.
	fmt.Fprintf(&b, `  <defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0,0 L10,5 L0,10 z" fill="%s"/></marker></defs>`, pal.Edge)
	b.WriteByte('\n')

	fmt.Fprintf(&b, `  <rect width="100%%" height="100%%" fill="%s"/>`, pal.Background)
	b.WriteByte('\n')

	writeTitle(&b, opts.Title, w/2, pad+opts.FontSize, pal)

	fmt.Fprintf(&b, `  <g transform="translate(%s,%s)">`, num(pad), num(pad+titleH))
	b.WriteByte('\n')

	for _, e := range res.Graph.Edges {
		writeEdge(&b, e, pal)
	}
	for _, n := range res.Graph.Nodes {
		writeNode(&b, n, pal, opts)
	}

	b.WriteString("  </g>\n</svg>\n")
	return []byte(b.String()), nil
}

func writeEdge(b *strings.Builder, e *domain.Edge, pal theme.Palette) {
	if len(e.Points) < 2 {
		return
	}
	var d strings.Builder
	for i, p := range e.Points {
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&d, "%s%s,%s ", cmd, num(p.X), num(p.Y))
	}
	dash := ""
	if e.Arrow == domain.ArrowDotted {
		dash = ` stroke-dasharray="4,4"`
	}
	width := "1.5"
	if e.Arrow == domain.ArrowThick {
		width = "3"
	}
	marker := ` marker-end="url(#arrow)"`
	if e.Arrow == domain.ArrowOpen {
		marker = ""
	}
	fmt.Fprintf(b, `    <path d="%s" fill="none" stroke="%s" stroke-width="%s"%s%s/>`,
		strings.TrimSpace(d.String()), pal.Edge, width, dash, marker)
	b.WriteByte('\n')

	if e.Label != "" {
		// Midpoint of the polyline's endpoints, not Points[len/2] (which is
		// the target endpoint for a two-point line).
		first, last := e.Points[0], e.Points[len(e.Points)-1]
		midX, midY := (first.X+last.X)/2, (first.Y+last.Y)/2
		fmt.Fprintf(b, `    <text x="%s" y="%s" fill="%s" text-anchor="middle" dy="-2">%s</text>`,
			num(midX), num(midY), pal.Text, esc(e.Label))
		b.WriteByte('\n')
	}
}

func writeNode(b *strings.Builder, n *domain.Node, pal theme.Palette, opts Options) {
	x, y, w, h := n.Pos.X, n.Pos.Y, n.Size.W, n.Size.H
	switch n.Shape {
	case domain.ShapeRound, domain.ShapeStadium:
		rx := h / 2
		if n.Shape == domain.ShapeRound {
			rx = 6
		}
		fmt.Fprintf(b, `    <rect x="%s" y="%s" width="%s" height="%s" rx="%s" fill="%s" stroke="%s"/>`,
			num(x), num(y), num(w), num(h), num(rx), pal.NodeFill, pal.NodeStroke)
	case domain.ShapeCircle:
		fmt.Fprintf(b, `    <circle cx="%s" cy="%s" r="%s" fill="%s" stroke="%s"/>`,
			num(x+w/2), num(y+h/2), num(w/2), pal.NodeFill, pal.NodeStroke)
	case domain.ShapeDiamond:
		cx, cy := x+w/2, y+h/2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s",
			num(cx), num(y), num(x+w), num(cy), num(cx), num(y+h), num(x), num(cy))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeHexagon:
		k := h / 2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s %s,%s %s,%s",
			num(x), num(y+h/2), num(x+k), num(y), num(x+w-k), num(y),
			num(x+w), num(y+h/2), num(x+w-k), num(y+h), num(x+k), num(y+h))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeParallelogram:
		k := h / 2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s",
			num(x+k), num(y), num(x+w), num(y), num(x+w-k), num(y+h), num(x), num(y+h))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeParallelogramAlt:
		k := h / 2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s",
			num(x), num(y), num(x+w-k), num(y), num(x+w), num(y+h), num(x+k), num(y+h))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeTrapezoid:
		k := h / 2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s",
			num(x+k), num(y), num(x+w-k), num(y), num(x+w), num(y+h), num(x), num(y+h))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeTrapezoidAlt:
		k := h / 2
		pts := fmt.Sprintf("%s,%s %s,%s %s,%s %s,%s",
			num(x), num(y), num(x+w), num(y), num(x+w-k), num(y+h), num(x+k), num(y+h))
		fmt.Fprintf(b, `    <polygon points="%s" fill="%s" stroke="%s"/>`, pts, pal.NodeFill, pal.NodeStroke)
	case domain.ShapeCylinder:
		ry := h * 0.12
		fmt.Fprintf(b, `    <path d="M%s,%s L%s,%s A%s,%s 0 0 0 %s,%s A%s,%s 0 0 0 %s,%s Z" fill="%s" stroke="%s"/>`,
			num(x), num(y+ry), num(x), num(y+h-ry),
			num(w/2), num(ry), num(x+w), num(y+h-ry),
			num(w/2), num(ry), num(x), num(y+ry),
			pal.NodeFill, pal.NodeStroke)
		fmt.Fprintf(b, `<path d="M%s,%s A%s,%s 0 0 0 %s,%s" fill="none" stroke="%s"/>`,
			num(x), num(y+ry), num(w/2), num(ry), num(x+w), num(y+ry), pal.NodeStroke)
	case domain.ShapeSubroutine:
		fmt.Fprintf(b, `    <rect x="%s" y="%s" width="%s" height="%s" fill="%s" stroke="%s"/>`,
			num(x), num(y), num(w), num(h), pal.NodeFill, pal.NodeStroke)
		fmt.Fprintf(b, `<line x1="%s" y1="%s" x2="%s" y2="%s" stroke="%s"/><line x1="%s" y1="%s" x2="%s" y2="%s" stroke="%s"/>`,
			num(x+6), num(y), num(x+6), num(y+h), pal.NodeStroke,
			num(x+w-6), num(y), num(x+w-6), num(y+h), pal.NodeStroke)
	default: // rect
		fmt.Fprintf(b, `    <rect x="%s" y="%s" width="%s" height="%s" fill="%s" stroke="%s"/>`,
			num(x), num(y), num(w), num(h), pal.NodeFill, pal.NodeStroke)
	}
	b.WriteByte('\n')

	label := n.Label
	if label == "" {
		label = n.ID
	}
	b.WriteString("    ")
	svgutil.MultilineText(b, svgutil.SplitLines(label), x+w/2, y+h/2+opts.FontSize*0.35, opts.FontSize+2, pal.Text, "")
	b.WriteByte('\n')
}
