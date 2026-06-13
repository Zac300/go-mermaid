package render

// palette holds the colors used to render a diagram.
type palette struct {
	Background string
	NodeFill   string
	NodeStroke string
	Text       string
	Edge       string
}

// themes maps theme names to palettes. Unknown names fall back to default.
var themes = map[string]palette{
	"default": {
		Background: "#ffffff",
		NodeFill:   "#ECECFF",
		NodeStroke: "#9370DB",
		Text:       "#333333",
		Edge:       "#333333",
	},
	"dark": {
		Background: "#1e1e1e",
		NodeFill:   "#2b2b40",
		NodeStroke: "#8888bb",
		Text:       "#e6e6e6",
		Edge:       "#bbbbbb",
	},
	"neutral": {
		Background: "#ffffff",
		NodeFill:   "#eeeeee",
		NodeStroke: "#999999",
		Text:       "#222222",
		Edge:       "#555555",
	},
}

func paletteFor(name string) palette {
	if p, ok := themes[name]; ok {
		return p
	}
	return themes["default"]
}
