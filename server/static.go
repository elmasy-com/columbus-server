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

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	c.String(http.StatusOK, "%s", openApiYaml)
}
