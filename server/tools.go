package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/elnet/domain"
	"github.com/gin-gonic/gin"
)

// GET /tools/tld/{fqdn}
// Returns the TLD part of a FQDN.
func ToolsTLDGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, fault.ErrBlocked.Err)
		} else {
			c.JSON(http.StatusForbidden, fault.ErrBlocked)
		}
		return
	}

	fqdn := c.Param("fqdn")

	if !domain.IsValid(fqdn) || fqdn == "." {
		c.Error(fault.ErrInvalidDomain)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	fqdn = strings.ToLower(fqdn)

	d := domain.GetTLD(fqdn)
	if d == "" {
		c.Error(fault.ErrNotFound)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrInvalidDomain)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, d)
	} else {
		c.JSON(http.StatusOK, gin.H{"result": d})
	}
}

// GET /tools/domain/{fqdn}
// Returns the domain part of a FQDN.
func ToolsDomainGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, fault.ErrBlocked.Err)
		} else {
			c.JSON(http.StatusForbidden, fault.ErrBlocked)
		}
		return
	}

	fqdn := c.Param("fqdn")

	if !domain.IsValid(fqdn) || fqdn == "." {
		c.Error(fault.ErrInvalidDomain)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	fqdn = strings.ToLower(fqdn)

	d := domain.GetDomain(fqdn)
	if d == "" {
		c.Error(fault.ErrNotFound)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrInvalidDomain)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, d)
	} else {
		c.JSON(http.StatusOK, gin.H{"result": d})
	}
}

// GET /tools/subdomain/{fqdn}
// Returns the subdomain part of a FQDN.
func ToolsSubdomainGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, fault.ErrBlocked.Err)
		} else {
			c.JSON(http.StatusForbidden, fault.ErrBlocked)
		}
		return
	}

	fqdn := c.Param("fqdn")

	if !domain.IsValid(fqdn) || fqdn == "." {
		c.Error(fault.ErrInvalidDomain)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	fqdn = strings.ToLower(fqdn)

	d := domain.GetSub(fqdn)
	if d == "" {
		c.Error(fault.ErrNotFound)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrInvalidDomain.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrInvalidDomain)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, d)
	} else {
		c.JSON(http.StatusOK, gin.H{"result": d})
	}
}

// GET /tools/isvalid/{fqdn}
// Returns wether fqdn is valid.
func ToolsIsValidGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusForbidden, fault.ErrBlocked.Err)
		} else {
			c.JSON(http.StatusForbidden, fault.ErrBlocked)
		}
		return
	}

	fqdn := c.Param("fqdn")
	fqdn = strings.ToLower(fqdn)

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, fmt.Sprintf("%v", domain.IsValid(fqdn)))
	} else {
		c.JSON(http.StatusForbidden, gin.H{"result": fmt.Sprintf("%v", domain.IsValid(fqdn))})
	}
}