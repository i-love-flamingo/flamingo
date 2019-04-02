package domain

type (
	// TokenExtras for openID Connect
	TokenExtras map[string]string
)

// Add a value to a key
func (te *TokenExtras) Add(key string, value string) {
	(*te)[key] = value
}

// Get a calue from a key
func (te *TokenExtras) Get(key string) (string, bool) {
	value, ok := (*te)[key]
	return value, ok
}
