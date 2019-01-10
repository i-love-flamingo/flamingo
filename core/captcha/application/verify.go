package application

type (
	Verifier struct {
		generator *Generator
	}
)

// Inject dependencies
func (v *Verifier) Inject(
	g *Generator,
) {
	v.generator = g
}

func (v *Verifier) Verify(hash, solution string) bool {
	captcha, err := v.generator.NewCaptchaByHash(hash)
	if err != nil {
		return false
	}

	if solution != captcha.Solution {
		return false
	}

	return true
}
