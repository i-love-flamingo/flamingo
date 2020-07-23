package interfaces

import (
	"bytes"
	"context"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"io/ioutil"
	"net/http"
)

type (
	// FileControllerInterface is the callback HTTP action provider
	FileControllerInterface interface {
		GetRobotsTxt(context.Context, *web.Request) web.Result
		GetSecurityTxt(context.Context, *web.Request) web.Result
		GetHumansTxt(context.Context, *web.Request) web.Result
	}

	// DefaultFileControllerInterface implementing FileControllerInterface
	DefaultFileControllerInterface struct {
		responder           *web.Responder
		logger              flamingo.Logger
		robotsTxtFilepath   string
		securityTxtFilepath string
		humansTxtFilepath   string
	}
)

var (
	_ FileControllerInterface = new(DefaultFileControllerInterface)
)

// Inject configuration
func (d *DefaultFileControllerInterface) Inject(
	responder *web.Responder,
	logger flamingo.Logger,
	config *struct {
		RobotsTxtFilepath   string `inject:"config:core.robotstxt.filepath"`
		SecurityTxtFilepath string `inject:"config:core.securitytxt.filepath"`
		HumansTxtFilepath   string `inject:"config:core.humanstxt.filepath"`
	},
) {
	d.responder = responder
	d.logger = logger.WithField("category", "robotstxt")
	if config != nil {
		d.robotsTxtFilepath = config.RobotsTxtFilepath
		d.securityTxtFilepath = config.SecurityTxtFilepath
		d.humansTxtFilepath = config.HumansTxtFilepath
	}
}

// GetRobotsTxt returns /robots.txt
func (d *DefaultFileControllerInterface) GetRobotsTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.robotsTxtFilepath)
}

// GetSecurityTxt returns /.well-known/security.txt
func (d *DefaultFileControllerInterface) GetSecurityTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.securityTxtFilepath)
}

// GetHumansTxt returns /humans.txt
func (d *DefaultFileControllerInterface) GetHumansTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.humansTxtFilepath)
}

func (d *DefaultFileControllerInterface) serveFile(ctx context.Context, filePath string) web.Result {
	fileContent, err := d.readFile(filePath)
	if err != nil {
		d.logger.WithContext(ctx).Error(err)

		d.responder.ServerError(err)
	}

	return d.responder.HTTP(http.StatusOK, bytes.NewReader(fileContent))
}

func (d *DefaultFileControllerInterface) readFile(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return b, nil
}
