package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/db"

	"github.com/gin-gonic/gin"
)

func InsertPut(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	_, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int

		switch {
		case errors.Is(err, fault.ErrMissingAPIKey):
			// X-Api-Key header is missing
			code = http.StatusUnauthorized

		case errors.Is(err, fault.ErrUserNotFound):
			// X-Api-Key is invalid
			code = http.StatusUnauthorized
			err = fault.ErrInvalidAPIKey

			blacklist.Block(c.ClientIP())

		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	err = db.Insert(c.Param("domain"))
	if err != nil {

		respCode := 0

		switch {
		case errors.Is(err, fault.ErrInvalidDomain):
			respCode = http.StatusBadRequest
		default:
			respCode = http.StatusInternalServerError
		}

		c.Error(err)
		c.JSON(respCode, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
