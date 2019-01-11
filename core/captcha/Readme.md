# Captcha

This module wraps [github.com/dchest/captcha](https://github.com/dchest/captcha) for usage within flamingo.

A captcha is represented by a base64-encoded encrypted hash of the solution.

## General Usage

Create a global instance of `application.Generator`. The generator holds the encryption key
and must therefor be used for all captchas within the project.

The Generator provides 3 methods to get captchas:

 * `NewCaptcha(length int)`
 * `NewCaptchaByHash(hash string) (*domain.Captcha, error)`
 * `NewCaptchaBySolution(solution string) *domain.Captcha`
 
Note that `NewCaptchaBySolution` generates a different hash each time, even if the solution is the same.

To check a string against a captcha, you can use `Verifier.Verify(hash, solution string) bool`.

See `application/example_test.go` for a full example.

## Usage in templates

This module registers 3 template functions:

 * `{{captcha}}` to generate a new Captcha
 * `{{captchaImage $captcha [true|false]}}` to get the image URL for a given captcha (optional as download)
 * `{{captchaSound $captcha [true|false]}}` to get the audio URL for a given captcha (optional as download)
 
```gotemplate
{{with $captcha:=captcha}}

	<form action="{{url "captcha"}}" method="post">
		<img src="{{captchaImage $captcha}}" alt="solve the captcha">
		<audio id=audio controls="controls" src="{{captchaSound $captcha}}" preload=none>
			You browser doesn't support audio.
		</audio>
		<a href="{{captchaSound $captcha true}}">Download audio</a>
		<input type="hidden" name="captcha_hash" value="{{$captcha.Hash}}">
		<label>
			Solution
			<input type="text" name="captcha" value="{{$captcha.Solution}}">
		</label>
		<button type="submit">Send</button>
	</form>

{{end}}
```
 
## Images and Audios

Images and audios are always generated on the fly by decrypting the hash. 
The module registers the route `"/captcha/*n"` for this purpose and assigns dchest/captcha's HTTPHandler to it.
Refer to https://godoc.org/github.com/dchest/captcha for image and audio generation documentation and language features.
