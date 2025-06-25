package api

// derefString safely dereferences a string pointer and returns its value.
// If the pointer is nil, it returns an empty string.
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// derefStringWithDefault safely dereferences a string pointer.
// If the pointer is nil, it returns the provided default value.
func derefStringWithDefault(s *string, defaultValue string) string {
	if s != nil {
		return *s
	}
	return defaultValue
}
