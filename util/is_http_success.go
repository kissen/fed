package util

// Returns whether status represents a "successful" response, i.e.
// whether status is in the range [200,299].
func IsHTTPSuccess(status int) bool {
	return status >= 200 && status <= 299
}
