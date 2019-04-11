package domain

import "encoding/gob"

type (
	// Idendity information
	Idendity interface {
		User() User
	}

	// User is a basic authenticated user
	User struct {
		Sub          string
		Name         string
		Email        string
		Salutation   string
		FirstName    string
		LastName     string
		Street       string
		ZipCode      string
		City         string
		DateOfBirth  string
		Country      string
		CustomFields map[string]string
		Groups       []string
	}
)


func init() {
	gob.Register(User{})
}

// Get a custom field by the name
func (u User) Get(name string) string {
	if u.CustomFields == nil {
		return ""
	}
	return u.CustomFields[name]
}