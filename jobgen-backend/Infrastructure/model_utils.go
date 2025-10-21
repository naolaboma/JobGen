package infrastructure

import "strings"

// normalizeModel resolves common aliases and removes unsupported suffixes for v1 API
func normalizeModel(name string) string {
	if name == "" {
		return "gemini-1.5-pro"
	}
	n := strings.TrimSpace(strings.ToLower(name))

	switch n {
	case "gemini-1.5-pro", "gemini-pro", "pro":
		return "gemini-1.5-pro"
	case "gemini-1.5-flash", "gemini-flash", "flash":
		return "gemini-1.5-flash"
	case "gemini-1.5-flash-8b", "flash-8b":
		return "gemini-1.5-flash-8b"
	}

	// If a -latest suffix was provided (e.g., from other SDKs), strip it for v1 REST/SDK
	n = strings.TrimSuffix(n, "-latest")

	return n
}
