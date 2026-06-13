// Package class parses and renders Mermaid class diagrams to SVG. It reuses
// the shared layered layout engine for positioning and edge routing, drawing
// UML class boxes (name / attributes / methods compartments) and relationship
// markers (inheritance, composition, aggregation, association, dependency).
package class

// Class is a UML class with attribute and method members.
type Class struct {
	Name       string
	Attributes []string
	Methods    []string
}

// headKind is a relationship line-end decoration.
type headKind int

const (
	headNone headKind = iota
	headArrow
	headTriangle      // inheritance / realization (hollow triangle)
	headDiamondFilled // composition
	headDiamondHollow // aggregation
)

// Relation is a relationship between two classes.
type Relation struct {
	From   string
	To     string
	Label  string
	Dashed bool
	Left   headKind // decoration at the From end
	Right  headKind // decoration at the To end
}

// Diagram is a parsed class diagram.
type Diagram struct {
	Classes   []*Class
	Relations []*Relation
}

func (d *Diagram) class(name string) *Class {
	for _, c := range d.Classes {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (d *Diagram) ensureClass(name string) *Class {
	if c := d.class(name); c != nil {
		return c
	}
	c := &Class{Name: name}
	d.Classes = append(d.Classes, c)
	return c
}
