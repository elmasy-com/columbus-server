package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	sdk "github.com/elmasy-com/columbus-sdk"
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

	user, err := db.UserGet(c.GetHeader("X-Api-Key"))
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

	user, err := db.UserGet(c.GetHeader("X-Api-Key"))
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

	user, err := db.UserGet(c.GetHeader("X-Api-Key"))
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

	user, err := db.UserGet(c.GetHeader("X-Api-Key"))
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

	user, err := db.UserGet(c.GetHeader("X-Api-Key"))
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
	if !user.Admin {

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("user must be admin")
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	out, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body: %w", err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(out) == 0 {
		err = fmt.Errorf("missing request body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userTarget := sdk.User{}

	err = json.Unmarshal(out, &userTarget)
	if err != nil {

		var code int
		var err error

		switch err.(type) {
		case *json.SyntaxError:
			code = http.StatusBadRequest
			err = fmt.Errorf("syntax error: %w", err)
		case *json.UnmarshalTypeError:
			code = http.StatusBadRequest
			err = fmt.Errorf("type error: %w", err)
		default:
			code = http.StatusInternalServerError
			err = fmt.Errorf("unmarshal error: %w", err)
		}

		c.Error(err)
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	if !user.Admin && user.Key != userTarget.Key {
		// Prevent to non admi user modify other user

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("only admins can modify other users")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userTarget.Key == "" {
		err := fmt.Errorf("body key is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userTarget.Name == "" {
		err := fmt.Errorf("body name is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := db.UserGet(userTarget.Key)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {

			err = fmt.Errorf("target user not exist")
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userTarget.Name != target.Name {
		err := fmt.Errorf("target name is not match with the key")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Need key update?
	upKey := c.DefaultQuery("key", "false")
	if upKey != "true" && upKey != "false" {
		err := fmt.Errorf("invalid value for parameter key: %s", upKey)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Need name update?
	upName := c.DefaultQuery("name", "")

	// Need admin update?
	upAdmin := c.Query("admin")
	if upAdmin != "" && !user.Admin {
		// Prevent user to set admin field to self and to other user

		blacklist.Block(c.ClientIP())

		err := fmt.Errorf("only admins can modify admin field")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if upAdmin != "true" && upAdmin != "false" {
		err := fmt.Errorf("invalid value for parameter admin: %s", upAdmin)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if upKey == "" && upName == "" && upAdmin == "" {

	}

	// if upKey == "true" {
	// 	err := db.UserChangeKey(&target)
	// 	if err != nil {
	// 		err := fmt.Errorf()
	// 	}
	// }
}
