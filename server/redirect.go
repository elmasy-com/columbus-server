package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Redirect(c *gin.Context) {

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/api%s", c.Request.RequestURI))
}
