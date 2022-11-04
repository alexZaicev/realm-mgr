package headers

// Carrier is used to get and set headers for a request.
//
// For compliance with HTTP/2, the header keys/names used as inputs should be normalised to
// lowercase. It is the responsibility of the implementation to normalise it to another format if
// required.
type Carrier interface {
	// GetSingle gets a single header from a request.
	//
	// Implementations should return a HeaderNotFound error if the header is not present, or a
	// MultipleHeadersFound error if more than one value for the header are present.
	GetSingle(key string) (string, error)

	// GetMultiple gets multiple header values from a request. An empty slice of values is returned
	// if the header is not found.
	GetMultiple(key string) []string

	// Set sets a single value for a request header, removing any existing values if present.
	Set(key, value string)
}
