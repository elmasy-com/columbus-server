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
		c.Error(fault.ErrBlocked)
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

	d = strings.ToLower(d)
	d = domain.GetDomain(d)
	if d == "" {
		c.Error(fmt.Errorf("%w: domain is empty", fault.ErrNotFound))
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrNotFound.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrNotFound)
		}
		return
	}

	subs, err := db.Lookup(d)
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
		c.Error(fault.ErrNotFound)
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
