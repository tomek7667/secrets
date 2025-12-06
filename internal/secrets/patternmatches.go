package secrets

import "strings"

const WildCardChar = "*"

func PatternMatches(key, pattern string) bool {
	if pattern == "" {
		return key == ""
	}
	if pattern == WildCardChar {
		return true
	}

	// No wildcard -> exact match
	if !strings.Contains(pattern, WildCardChar) {
		return key == pattern
	}

	parts := strings.Split(pattern, WildCardChar)
	s := key
	idx := 0

	// If pattern doesn't start with "*", first part must be a prefix
	if parts[0] != "" {
		if !strings.HasPrefix(s, parts[0]) {
			return false
		}
		idx = len(parts[0])
	}

	// Match middle parts in order
	for i := 1; i < len(parts)-1; i++ {
		part := parts[i]
		if part == "" {
			continue
		}
		pos := strings.Index(s[idx:], part)
		if pos == -1 {
			return false
		}
		idx += pos + len(part)
	}

	// If pattern doesn't end with "*", last part must be a suffix
	last := parts[len(parts)-1]
	if last != "" {
		if len(s) < len(last) || !strings.HasSuffix(s, last) {
			return false
		}
		if idx > len(s)-len(last) {
			return false
		}
	}

	return true
}
