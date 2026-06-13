package parser

import (
	"testing"

	"github.com/Zac300/go-mermaid/internal/domain"
)

func TestDirections(t *testing.T) {
	cases := map[string]domain.Direction{
		"graph TD\nA": domain.TopBottom,
		"graph TB\nA": domain.TopBottom,
		"graph BT\nA": domain.BottomTop,
		"graph LR\nA": domain.LeftRight,
		"graph RL\nA": domain.RightLeft,
	}
	for src, want := range cases {
		g, err := parse(src)
		if err != nil {
			t.Fatalf("parse(%q): %v", src, err)
		}
		if g.Direction != want {
			t.Errorf("parse(%q) direction = %v, want %v", src, g.Direction, want)
		}
	}
}

func TestArrowKinds(t *testing.T) {
	cases := map[string]domain.Arrow{
		"graph TD\nA --> B":  domain.ArrowNormal,
		"graph TD\nA --- B":  domain.ArrowOpen,
		"graph TD\nA -.-> B": domain.ArrowDotted,
		"graph TD\nA ==> B":  domain.ArrowThick,
	}
	for src, want := range cases {
		g, err := parse(src)
		if err != nil {
			t.Fatalf("parse(%q): %v", src, err)
		}
		if g.Edges[0].Arrow != want {
			t.Errorf("parse(%q) arrow = %v, want %v", src, g.Edges[0].Arrow, want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	bad := []string{
		"graph XY\nA --> B",   // unknown direction
		"graph TD A --> B",    // junk after header
		"graph TD\n--> B",     // missing node id
		"graph TD\nA -->",     // missing target
		"graph TD\nA -->|x B", // unterminated label handled at lex; ensure error
	}
	for _, src := range bad {
		if _, err := parse(src); err == nil {
			t.Errorf("parse(%q) expected error, got nil", src)
		}
	}
}

func TestFlowchartKeyword(t *testing.T) {
	g, err := parse("flowchart LR\nA --> B")
	if err != nil {
		t.Fatalf("flowchart keyword: %v", err)
	}
	if g.Direction != domain.LeftRight {
		t.Errorf("direction = %v, want LeftRight", g.Direction)
	}
}
