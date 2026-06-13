package quadrant

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("quadrantChart\ntitle T\nx-axis A --> B\ny-axis C --> D\nquadrant-1 Q\nP: [0.3, 0.6]")
	f.Add("quadrantChart\nP: []")
	f.Add("quadrantChart\nP: [1]")
	f.Add("quadrantChart")
	f.Fuzz(func(t *testing.T, src string) {
		_, _ = Parse(src)
	})
}
