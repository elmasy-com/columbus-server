package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/elmasy-com/columbus-server/db"
	"github.com/elmasy-com/columbus-server/fault"
	"github.com/elmasy-com/elnet/domain"
	"github.com/gin-gonic/gin"
)

func LookupGet(c *gin.Context) {

	var err error

	d := c.Param("domain")

	subs, err := db.Lookup(d)
	if err != nil {

		c.Error(err)

		respCode := 0

		switch {
		case errors.Is(err, fault.ErrInvalidDomain):
			respCode = http.StatusBadRequest
		default:
			respCode = http.StatusInternalServerError
			err = fmt.Errorf("internal server error")
		}

		if c.GetHeader("Accept") == "text/plain" {
			c.String(respCode, err.Error())
		} else {
			c.JSON(respCode, gin.H{"error": err.Error()})
		}
		return
	}

	if len(subs) == 0 {

		c.Error(fault.ErrNotFound)

		_, err = db.InsertNotFound(d)
		if err != nil {
			c.Error(fmt.Errorf("failed to insert notFound: %w", err))
		}

		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrNotFound.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrNotFound)
		}
		return
	}

	_, err = db.InsertTopList(d)
	if err != nil {
		c.Error(fmt.Errorf("failed to insert topList: %w", err))
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, strings.Join(subs, "\n"))
	} else {
		c.JSON(http.StatusOK, subs)
	}
}

func TLDGet(c *gin.Context) {

	dom := c.Param("domain")

	if !domain.IsValidSLD(dom) {

		c.Error(fault.ErrInvalidDomain)

		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Error())
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	dom = domain.Clean(dom)

	tlds, err := db.TLD(dom)
	if err != nil {

		c.Error(err)

		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusInternalServerError, "internal server error")
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if len(tlds) == 0 {
		c.Error(fault.ErrNotFound)
		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrNotFound.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrNotFound)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, strings.Join(tlds, "\n"))
	} else {
		c.JSON(http.StatusOK, tlds)
	}
}

func StartsGet(c *gin.Context) {

	dom := c.Param("domain")

	if len(dom) < 5 {

		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusBadRequest, fault.ErrInvalidDomain.Error())
		} else {
			c.JSON(http.StatusBadRequest, fault.ErrInvalidDomain)
		}
		return
	}

	domains, err := db.Starts(dom)
	if err != nil {

		c.Error(err)
		code := 0

		if errors.Is(err, fault.ErrInvalidDomain) {
			code = http.StatusBadRequest
		} else {
			code = http.StatusInternalServerError
			err = fmt.Errorf("internal server error")
		}

		if c.GetHeader("Accept") == "text/plain" {
			c.String(code, err.Error())
		} else {
			c.JSON(code, err)
		}
		return
	}

	if len(domains) == 0 {

		c.Error(fault.ErrNotFound)

		if c.GetHeader("Accept") == "text/plain" {
			c.String(http.StatusNotFound, fault.ErrNotFound.Err)
		} else {
			c.JSON(http.StatusNotFound, fault.ErrNotFound)
		}
		return
	}

	if c.GetHeader("Accept") == "text/plain" {
		c.String(http.StatusOK, strings.Join(domains, "\n"))
	} else {
		c.JSON(http.StatusOK, domains)
	}
}
