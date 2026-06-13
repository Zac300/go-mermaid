package mindmap

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("mindmap\n  root((R))\n    A\n      B\n    C")
	f.Add("mindmap\nroot")
	f.Add("mindmap\n      deep\n  shallow")
	f.Add("mindmap")
	f.Fuzz(func(t *testing.T, src string) {
		_, _ = Parse(src)
	})
}
