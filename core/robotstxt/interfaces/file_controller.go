package interfaces

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// FileController implementing FileController
	FileController struct {
		responder           *web.Responder
		logger              flamingo.Logger
		robotsTxtFilepath   string
		securityTxtFilepath string
		humansTxtFilepath   string
	}
)

// Inject configuration
func (d *FileController) Inject(
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
func (d *FileController) GetRobotsTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.robotsTxtFilepath)
}

// GetSecurityTxt returns /.well-known/security.txt
func (d *FileController) GetSecurityTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.securityTxtFilepath)
}

// GetHumansTxt returns /humans.txt
func (d *FileController) GetHumansTxt(ctx context.Context, _ *web.Request) web.Result {
	return d.serveFile(ctx, d.humansTxtFilepath)
}

func (d *FileController) serveFile(ctx context.Context, filePath string) web.Result {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		d.logger.WithContext(ctx).Error(err)

		d.responder.ServerError(err)
	}

	return d.responder.HTTP(http.StatusOK, bytes.NewReader(fileContent))
}
