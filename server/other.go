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

// GET /other
// Return a other user based on username.
func OtherGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
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

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		c.Error(fault.ErrNotAdmin)
		c.JSON(http.StatusForbidden, fault.ErrNotAdmin)
		return
	}

	target, err := db.UserGetName(c.Query("username"))
	if err != nil {

		var code int

		switch {
		case errors.Is(err, fault.ErrNameEmpty):
			code = http.StatusBadRequest
			err = fault.ErrUserNameEmpty
		case errors.Is(err, fault.ErrUserNotFound):
			code = http.StatusNotFound
		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get target user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, target)
}

// PATCH /other/key
// Change other user's key.
func OtherKeyPatch(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
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

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		c.Error(fault.ErrNotAdmin)
		c.JSON(http.StatusForbidden, fault.ErrNotAdmin)
		return
	}

	target, err := db.UserGetName(c.Query("username"))
	if err != nil {

		var code int

		switch {
		case errors.Is(err, fault.ErrNameEmpty):
			code = http.StatusBadRequest
			err = fault.ErrUserNameEmpty
		case errors.Is(err, fault.ErrUserNotFound):
			code = http.StatusNotFound
		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get target user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	err = db.UserChangeKey(&target)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// PATCH /other/name
// Change other user's name.
func OtherNamePatch(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
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

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		c.Error(fault.ErrNotAdmin)
		c.JSON(http.StatusForbidden, fault.ErrNotAdmin)
		return
	}

	target, err := db.UserGetName(c.Query("username"))
	if err != nil {

		var code int

		switch {
		case errors.Is(err, fault.ErrNameEmpty):
			code = http.StatusBadRequest
			err = fault.ErrUserNameEmpty
		case errors.Is(err, fault.ErrUserNotFound):
			code = http.StatusNotFound
		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get target user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	if target.Name == c.Query("name") {
		c.Error(fault.ErrSameName)
		c.JSON(http.StatusBadRequest, fault.ErrSameName)
		return
	}

	taken, err := db.IsNameTaken(c.Query("name"))
	if err != nil {

		code := 0

		switch {
		case errors.Is(err, fault.ErrNameEmpty):
			code = http.StatusBadRequest
		default:
			code = http.StatusInternalServerError
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	if taken {
		c.Error(fault.ErrNameTaken)
		c.JSON(http.StatusConflict, fault.ErrNameTaken)
		return
	}

	err = db.UserChangeName(&target, c.Query("name"))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// PATCH /other/admin
// Change other user's admin value.
func OtherAdminPatch(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
		c.Error(fault.ErrBlocked)
		c.JSON(http.StatusForbidden, fault.ErrBlocked)
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
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

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		c.Error(fault.ErrNotAdmin)
		c.JSON(http.StatusForbidden, fault.ErrNotAdmin)
		return
	}

	target, err := db.UserGetName(c.Query("username"))
	if err != nil {

		var code int

		switch {
		case errors.Is(err, fault.ErrNameEmpty):
			code = http.StatusBadRequest
			err = fault.ErrUserNameEmpty
		case errors.Is(err, fault.ErrUserNotFound):
			// X-Api-Key is invalid
			code = http.StatusBadRequest
		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to get target user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	admin := false

	switch c.Query("admin") {
	case "true":
		admin = true
	case "false":
		admin = false
	case "":
		err := fmt.Errorf("admin parameter is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	default:
		err := fmt.Errorf("invalid value for admin: %s", c.Query("admin"))
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if target.Admin == admin {
		c.Error(fault.ErrNothingToDo)
		c.JSON(http.StatusBadRequest, fault.ErrNothingToDo)
		return
	}

	err = db.UserChangeAdmin(&target, admin)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}
