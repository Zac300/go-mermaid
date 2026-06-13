package git

import "testing"

func FuzzParse(f *testing.F) {
	f.Add("gitGraph\ncommit\nbranch dev\ncommit\ncheckout main\nmerge dev tag: \"v1\"")
	f.Add("gitGraph\nmerge nothing")
	f.Add("gitGraph:")
	f.Add("gitGraph")
	f.Fuzz(func(t *testing.T, src string) {
		_, _ = Parse(src)
	})
}
