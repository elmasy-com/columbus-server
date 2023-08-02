package server

import (
	"net/http"

	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
)

func StatGet(c *gin.Context) {

	s, err := db.StatisticsGetNewest()
	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, s)
}
