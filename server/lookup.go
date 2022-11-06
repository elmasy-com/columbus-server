package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/db"
	"github.com/elmasy-com/elnet/domain"
	"github.com/gin-gonic/gin"
)

func LookupGet(c *gin.Context) {

	// Block blacklisted IPs
	if blacklist.IsBlocked(c.ClientIP()) {
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, fault.ErrBlocked.Err)
		} else {
			c.JSON(http.StatusForbidden, fault.ErrBlocked)
		}
		return
	}

	var err error

	d := c.Param("domain")

	if !domain.IsValid(d) {
		c.Error(fault.ErrInvalidDomain)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Error())
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	var full bool
	getFull := c.DefaultQuery("full", "false")

	switch getFull {
	case "true":
		full = true
	case "false":
		// Just to check
		full = false
	default:
		err = fmt.Errorf("invalid value for full: %s", getFull)
		c.Error(err)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	subs, err := db.Lookup(d, full)
	if err != nil {
		c.Error(err)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if len(subs) == 0 {
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrNotFound.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrNotFound)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, strings.Join(subs, "\n"))
	} else {
		c.JSON(http.StatusOK, subs)
	}
}
