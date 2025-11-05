package infrastructure

import "strings"

// normalizeModel returns a canonical model ID for SDK usage (strips -latest).
func normalizeModel(name string) string {
	n := strings.TrimSpace(strings.ToLower(name))
	if n == "" {
		// Safer default that is broadly available and cost-effective
		return "gemini-1.5-flash-8b"
	}

	switch n {
	case "gemini-1.5-pro", "gemini-pro", "pro":
		return "gemini-1.5-pro"
	case "gemini-1.5-flash", "gemini-flash", "flash":
		return "gemini-1.5-flash"
	case "gemini-1.5-flash-8b", "flash-8b", "8b":
		return "gemini-1.5-flash-8b"
	}

	// Strip SDK-incompatible suffixes like -latest
	n = strings.TrimSuffix(n, "-latest")
	return n
}

// normalizeModelV1Beta returns a v1beta-friendly model name (ensures -latest suffix).
// Use this when calling the REST v1beta endpoint.
func normalizeModelV1Beta(name string) string {
	base := normalizeModel(name)
	if strings.HasSuffix(base, "-latest") {
		return base
	}
	return base + "-latest"
}
