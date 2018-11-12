package domain

type (
	TokenExtras map[string]string
)

func (te *TokenExtras) Add(key string, value string) {
	(*te)[key] = value
}

func (te *TokenExtras) Get(key string) (string, bool) {
	value, ok := (*te)[key]
	return value, ok
}
