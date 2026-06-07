package scanner

import "strings"

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimLeft(v, "^~>=< ")
	v = strings.TrimPrefix(v, "v")
	v = strings.Split(v, "+")[0]
	v = strings.Split(v, "-")[0]
	return strings.TrimSpace(v)
}

func splitGroupArtifact(name string) []string {
	parts := strings.Split(name, ":")
	return parts
}

func stringsReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}
