package server

import (
	"embed"
	"fmt"
	"io"
	"mime"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFS embed.FS

func NoRouteHandler(c *gin.Context) {

	p := c.Request.URL.Path
	if p == "/" {
		c.Redirect(301, "/index.html")
		return
	}

	p = "static" + p

	file, err := staticFS.Open(p)
	if err != nil {
		c.Error(err)
		c.Status(404)
		return
	}

	out, err := io.ReadAll(file)
	if err != nil {
		c.Error(err)
		c.Status(404)
		return
	}

	fileStat, err := file.Stat()
	if err != nil {
		c.Error(err)
		c.Status(404)
		return
	}

	i := strings.LastIndexByte(fileStat.Name(), '.')

	if i == -1 {
		// Unknown types will be application/octet-stream
		c.Data(200, "application/octet-stream", out)
		return
	}

	// Used to store the type
	t := ""

	switch ext := fileStat.Name()[i:]; ext {
	case ".ico":
		t = "image/vnd.microsoft.icon"
	case ".js":
		t = "text/javascript"
	case ".css":
		t = "text/css"
	case ".svg":
		t = "image/svg+xml"
	case ".html":
		t = "text/html"
	case ".yaml":
		t = "application/octet-stream"
	default:
		fmt.Printf("Unknown file extension: %s\n", ext)

		t = mime.TypeByExtension(ext)

		if t == "" {
			t = "application/octet-stream"
		}
	}

	c.Data(200, t, out)
}
