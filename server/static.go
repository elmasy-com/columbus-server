package server

import (
	_ "embed"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openApiYaml []byte

func StaticOpenApiYamlGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.String(http.StatusOK, string(openApiYaml))
}

//go:embed dist/index.html
var indexHtml []byte

func StaticIndexHtmlGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.Data(http.StatusOK, "text/html", indexHtml)
}

//go:embed dist/assets/favicon.d5f09fd4.ico
var faviconIco []byte

func StaticFaviconIcoGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.Data(http.StatusOK, "image/vnd.microsoft.icon", faviconIco)
}

//go:embed dist/assets/index.20b46c90.css
var indexCss []byte

func StaticIndexCssGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.Data(http.StatusOK, "text/css", indexCss)
}

//go:embed dist/assets/index.d727ceaf.js
var indexJs []byte

func StaticIndexJsGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.Data(http.StatusOK, "text/javascript", indexJs)
}

//go:embed dist/assets/logo_white.66566ab4.svg
var logoWhiteSvg []byte

func StaticLogoWhiteSvgGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.Data(http.StatusOK, "image/svg+xml", logoWhiteSvg)
}
