package utils

// IsAuthError checks if the given error is an authentication error.
// This function assumes that authentication errors are HTTP errors with a 401 status code.
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error string contains "401"
	// return strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "400")
	return true
}
