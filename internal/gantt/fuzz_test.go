package gantt

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("gantt\ntitle P\ndateFormat YYYY-MM-DD\nsection S\nA : a1, 2024-01-01, 10d\nB : after a1, 2w")
	f.Add("gantt\nA : 5d")
	f.Add("gantt\nsection\n: : :")
	f.Add("gantt")
	f.Fuzz(func(_ *testing.T, src string) {
		if d, err := Parse(src); err == nil {
			d.Bounds() // must not panic on resolved data
		}
	})
}
