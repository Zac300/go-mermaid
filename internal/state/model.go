// Package state parses and renders Mermaid state diagrams (stateDiagram-v2)
// to SVG, reusing the shared layered layout engine. States are rounded boxes;
// [*] start/end pseudostates render as filled and ringed circles.
package state

const (
	startID = "__start__"
	endID   = "__end__"
)

// State is a state node. Start/End mark the [*] pseudostates.
type State struct {
	ID    string
	Label string
	Start bool
	End   bool
}

// Transition is an arrow between two states.
type Transition struct {
	From  string
	To    string
	Label string
}

// Diagram is a parsed state diagram.
type Diagram struct {
	States      []*State
	Transitions []*Transition
}

func (d *Diagram) state(id string) *State {
	for _, s := range d.States {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (d *Diagram) ensureState(id string) *State {
	if s := d.state(id); s != nil {
		return s
	}
	s := &State{ID: id, Label: id}
	d.States = append(d.States, s)
	return s
}

// ensurePseudo returns the shared start or end pseudostate, creating it once.
func (d *Diagram) ensurePseudo(end bool) *State {
	id := startID
	if end {
		id = endID
	}
	if s := d.state(id); s != nil {
		return s
	}
	s := &State{ID: id, Start: !end, End: end}
	d.States = append(d.States, s)
	return s
}
