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

// GET /user
// Return the current user based on X-Api-Key
func UserGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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
		}

		c.Error(err)

		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /user
// Create a new user
func UserPut(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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

	name := c.Query("name")
	if name == "" {
		c.Error(fault.ErrNameEmpty)
		c.JSON(http.StatusBadRequest, fault.ErrNameEmpty)
		return
	}

	admin := false

	switch adminStr := c.DefaultQuery("admin", "false"); adminStr {
	case "true":
		admin = true
	case "false":
		admin = false
	default:
		err := fmt.Errorf("admin must be boolean, got: %s", adminStr)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := db.UserCreate(name, admin)
	if err != nil {

		code := 0

		switch {
		case errors.Is(err, fault.ErrNameTaken):
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// DELETE /user
func UserDelete(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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

	switch c.Query("confirmation") {
	case "":
		c.Error(fault.ErrConfirmMissing)
		c.JSON(http.StatusBadRequest, fault.ErrConfirmMissing)
	case "true":
		err = db.UserDelete(user.Key, user.Name)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	case "false":
		c.Error(fault.ErrNotConfirmed)
		c.JSON(http.StatusBadRequest, fault.ErrNotConfirmed)
	default:
		err = fmt.Errorf("invalid value for confirmation: %s", c.Query("confirmation"))
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// PATCH /user/key
// Regenerate the user key
func UserKeyPatch(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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

	err = db.UserChangeKey(&user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PATCH /user/key
// Regenerate the user key
func UserNamePatch(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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

	name := c.Query("name")
	if name == "" {
		c.Error(fault.ErrNameEmpty)
		c.JSON(http.StatusBadRequest, fault.ErrNameEmpty)
		return
	}

	err = db.UserChangeName(&user, name)
	if err != nil {

		code := 0

		switch {
		case errors.Is(err, fault.ErrNameTaken):
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}
		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GET /users
// Returns a list of evry user
func UsersGet(c *gin.Context) {

	if blacklist.IsBlocked(c.ClientIP()) {
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

	users, err := db.UserList()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
