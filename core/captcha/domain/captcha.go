package domain

type (
	// Captcha contains the captcha solution as clear text and its encrypted Hash
	// The hash is base64 encoded and save to use un URLs
	Captcha struct {
		Solution string
		Hash     string
	}
)
