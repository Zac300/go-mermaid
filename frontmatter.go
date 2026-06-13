package mermaid

import "strings"

// parseFrontmatter splits an optional leading YAML front-matter block
// (delimited by --- lines) from the diagram body and extracts the title.
// Only the "title" key is interpreted; other keys are ignored. When no
// front-matter is present, title is empty and body is the original source.
func parseFrontmatter(src string) (title, body string) {
	s := strings.TrimLeft(src, "\r\n")
	if !strings.HasPrefix(s, "---") {
		return "", src
	}
	lines := strings.Split(s, "\n")
	if strings.TrimSpace(lines[0]) != "---" {
		return "", src
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			for _, fm := range lines[1:i] {
				k, v, ok := strings.Cut(fm, ":")
				if ok && strings.TrimSpace(k) == "title" {
					title = strings.Trim(strings.TrimSpace(v), `"'`)
				}
			}
			return title, strings.Join(lines[i+1:], "\n")
		}
	}
	// Unterminated front-matter: treat the whole thing as body.
	return "", src
}
