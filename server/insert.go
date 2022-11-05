package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/db"

	"github.com/gin-gonic/gin"
)

func InsertPut(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, "blocked")
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		}
		return
	}

	_, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int
		var err error

		switch {
		case errors.Is(err, db.ErrUserKeyEmpty):
			// X-Api-Key header is missing
			code = http.StatusUnauthorized
			err = fmt.Errorf("missing X-Api-Key")

		case errors.Is(err, db.ErrUserNotFound):
			// X-Api-Key is invalid
			code = http.StatusUnauthorized
			err = fmt.Errorf("invalid X-Api-Key")

			blacklist.Block(c.ClientIP())

		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get user: %w", err)
		}

		c.Error(err)

		if c.GetHeader("Accept") == "text/plain" {
			c.String(code, err.Error())
		} else {
			c.JSON(code, gin.H{"error": err.Error()})
		}

		return
	}

	err = db.Insert(c.Param("domain"))
	if err != nil {

		respCode := 0

		switch {
		case errors.Is(err, db.ErrInvalidDomain):
			respCode = http.StatusBadRequest
			err = db.ErrInvalidDomain
		case strings.Contains(err.Error(), "cannot derive eTLD+1 for domain"):
			respCode = http.StatusBadRequest
			err = fmt.Errorf("domain is a public suffix")
		default:
			respCode = http.StatusInternalServerError
		}

		c.Error(err)

		if c.GetHeader("Accept") == "text/plain" {
			c.String(respCode, err.Error())
		} else {
			c.JSON(respCode, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}
