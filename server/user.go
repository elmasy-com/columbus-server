package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
)

// GET /user
// Return the current user based on X-Api-Key
func UserGet(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int
		var err2 error

		switch {
		case errors.Is(err, db.ErrUserKeyEmpty):
			// X-Api-Key header is missing
			code = http.StatusUnauthorized
			err2 = fmt.Errorf("missing X-Api-Key")

		case errors.Is(err, db.ErrUserNotFound):
			// X-Api-Key is invalid
			code = http.StatusUnauthorized
			err2 = fmt.Errorf("invalid X-Api-Key")

			blacklist.Block(c.ClientIP())

		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err2 = fmt.Errorf("failed to get user: %w", err)
		}

		c.Error(err2)

		c.JSON(code, gin.H{"error": err2.Error()})

		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /user
// Create a new user
func UserPut(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int
		var err2 error

		switch {
		case errors.Is(err, db.ErrUserKeyEmpty):
			// X-Api-Key header is missing
			code = http.StatusUnauthorized
			err2 = fmt.Errorf("missing X-Api-Key")

		case errors.Is(err, db.ErrUserNotFound):
			// X-Api-Key is invalid
			code = http.StatusUnauthorized
			err2 = fmt.Errorf("invalid X-Api-Key")

			blacklist.Block(c.ClientIP())

		default:
			// Server error while trying to get user
			code = http.StatusInternalServerError
			err2 = fmt.Errorf("failed to get user: %w", err)
		}

		c.Error(err2)

		c.JSON(code, gin.H{"error": err2.Error()})

		return
	}

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("not admin")
		c.Error(err)

		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")
	if name == "" {
		err := fmt.Errorf("name is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		case errors.Is(err, db.ErrUserNameTaken):
			code = http.StatusConflict
			err = db.ErrUserNameTaken
		default:
			code = http.StatusInternalServerError
			err = fmt.Errorf("failed to create user: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// DELETE /user
func UserDelete(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int

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
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	switch c.Query("confirmation") {
	case "":
		err := fmt.Errorf("confirmation is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case "true":
		err = db.UserDelete(user.Key, user.Name)
		if err != nil {
			err = fmt.Errorf("failed to delete user: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	case "false":
		err = fmt.Errorf("delete must be confirmed")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		err = fmt.Errorf("invalid value for confirmation: %s", c.Query("confirmation"))
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// PATCH /user
func UserPatch(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int

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
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	upKey := c.Query("key")
	upName := c.Query("name")

	switch {
	case upKey != "true" && upKey != "false" && upKey != "":
		// key paramaters must be a bool or empty
		err := fmt.Errorf("invalid value for key: %s", upKey)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case upKey == "" && upName == "", upName == "" && upKey == "false":
		// Query params are empty or name is empty and key is false, nothing to do
		err := fmt.Errorf("nothing to do")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case upKey != "" && upName != "":
		// Cant change both parameter at once.
		err := fmt.Errorf("two update at a time")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case user.Name == upName:
		// Old and new name are same
		err := fmt.Errorf("old and new name are same")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case upName != "":

		err = db.UserChangeName(&user, upName)
		if err != nil {

			code := 0

			switch {
			case errors.Is(err, db.ErrUserNameTaken):
				code = http.StatusConflict
			default:
				code = http.StatusInternalServerError
				err = fmt.Errorf("failed to update name: %w", err)
			}

			c.Error(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
		return

	case upKey == "true":
		err := db.UserChangeKey(&user)
		if err != nil {
			err = fmt.Errorf("failed to update key: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// PATCH /user/other
func UserOtherPatch(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int

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
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("not admin")
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	target, err := db.UserGetName(c.Query("username"))
	if err != nil {

		code := 0

		switch {
		case errors.Is(err, db.ErrUserNotFound):
			code = http.StatusNotFound
		case errors.Is(err, db.ErrUserNameEmpty):
			code = http.StatusBadRequest
			err = fmt.Errorf("username is empty")
		default:
			code = http.StatusInternalServerError
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	key := c.Query("key")
	if key != "" && key != "true" && key != "false" {
		err := fmt.Errorf("invalid value for key: %s", key)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")
	if name != "" {

		if key != "" {
			err := fmt.Errorf("one update at a time")
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if target.Name == name {
			err := fmt.Errorf("same name")
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		taken, err := db.IsNameTaken(name)
		if err != nil {
			err = fmt.Errorf("failed to check if name is taken: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if taken {
			err = fmt.Errorf("name is taken")
			c.Error(err)
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
	}

	adminBool := false

	admin := c.Query("admin")
	if admin != "" {

		switch admin {
		case "true":
			adminBool = true
		case "false":
			adminBool = false
		default:
			err := fmt.Errorf("invalid value for admin: %s", admin)
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if key != "" || name != "" {
			err := fmt.Errorf("one update at a time")
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	switch {
	case key == "true":
		err := db.UserChangeKey(&target)
		if err != nil {
			err = fmt.Errorf("failed to change key: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case name != "":
		err := db.UserChangeName(&target, name)
		if err != nil {
			err = fmt.Errorf("failed to change name: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case admin != "":
		err := db.UserChangeAdmin(&target, adminBool)
		if err != nil {
			err = fmt.Errorf("failed to change admin: %w", err)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		err := fmt.Errorf("nothing to do")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}

func UsersGet(c *gin.Context) {

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	user, err := db.UserGetKey(c.GetHeader("X-Api-Key"))
	if err != nil {

		var code int

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
		c.JSON(code, gin.H{"error": err.Error()})

		return
	}

	if !user.Admin {

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("not admin")
		c.Error(err)

		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	users, err := db.UserList()
	if err != nil {
		err = fmt.Errorf("failed to get users: %w", err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
