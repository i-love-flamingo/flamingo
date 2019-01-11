package application_test

import (
	"fmt"

	"flamingo.me/flamingo/core/captcha/application"
)

// Example generates a new captcha instance and calls the verify process with a wrong and a right solution
// In reality the solution for the Verify call would be the user input
func Example() {
	generator := &application.Generator{}
	generator.Inject(
		&struct {
			EncryptionPassPhrase string `inject:"config:captcha.encryptionPassPhrase"`
		}{
			EncryptionPassPhrase: "example",
		},
	)

	c := generator.NewCaptchaBySolution("123456")

	verifier := application.Verifier{}
	verifier.Inject(generator)

	wrong := verifier.Verify(c.Hash, "654321")
	right := verifier.Verify(c.Hash, "123456")

	fmt.Println("Wrong:", wrong, "Right:", right)

	// Output: Wrong: false Right: true
}
