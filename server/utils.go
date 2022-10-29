package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sdk "github.com/elmasy-com/columbus-sdk"
	"github.com/gin-gonic/gin"
)

// bodyToUser unmarhsal JSON Post body to sdk.User.
//
// The fields are not checked in this function.
//
// If any error occured, returns nil and set the status code and the error in context (so dont need to handle outside of this function).
func bodyToUser(c *gin.Context) *sdk.User {

	out, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body: %w", err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil
	}

	if len(out) == 0 {
		err = fmt.Errorf("missing request body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	user := sdk.User{}

	err = json.Unmarshal(out, &user)
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
		return nil
	}

	return &user
}
