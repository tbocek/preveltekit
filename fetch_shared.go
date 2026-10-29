package preveltekit

// FetchError provides detailed HTTP error information
type FetchError struct {
	Status     int
	StatusText string
	URL        string
}

func (e *FetchError) Error() string {
	s := "fetch " + e.URL + ": "
	if e.StatusText != "" {
		return s + itoa(e.Status) + " " + e.StatusText
	}
	return s + "HTTP " + itoa(e.Status)
}

// Get fetches JSON from a URL and decodes it into a typed struct.
// Must be called from a goroutine.
func Get[T any](url string) (T, error) {
	return request[T]("GET", url)
}

// Post sends a POST request with JSON body and decodes the response.
// Must be called from a goroutine.
func Post[T any](url string, body any) (T, error) {
	return requestBody[T]("POST", url, body)
}

// Put sends a PUT request with JSON body and decodes the response.
// Must be called from a goroutine.
func Put[T any](url string, body any) (T, error) {
	return requestBody[T]("PUT", url, body)
}

// Patch sends a PATCH request with JSON body and decodes the response.
// Must be called from a goroutine.
func Patch[T any](url string, body any) (T, error) {
	return requestBody[T]("PATCH", url, body)
}

// Delete sends a DELETE request and decodes the response.
// Must be called from a goroutine.
func Delete[T any](url string) (T, error) {
	return request[T]("DELETE", url)
}
