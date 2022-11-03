package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Date      int64      `json:"date"`
	DomainNum int64      `json:"domain"`
	SubNum    int64      `json:"sub"`
	m         sync.Mutex `json:"-"`
}

var (
	Current         Stat
	IsUpdateRunning bool // Indicate that the UpdateStat() goroutine is running.
)

func (s *Stat) Update(DomainNum int64, SubNum int64) {

	s.m.Lock()
	defer s.m.Unlock()

	s.Date = time.Now().Unix()
	s.DomainNum = DomainNum
	s.SubNum = SubNum
}

func (s *Stat) Copy() Stat {

	s.m.Lock()
	defer s.m.Unlock()

	return Stat{Date: s.Date, DomainNum: s.DomainNum, SubNum: s.SubNum}
}

// UpdateStat is created to run a goroutine.
// Started at the first call to GET /stat.
// Normally, UpdateStat updates the Current variable every X + 60 minute, but in case of error, UpdateStat retries in X + 30 minute.
func UpdateStat() {

	for {

		d, s, err := db.GetStat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
			time.Sleep(30 * time.Minute)
			continue
		}

		Current.Update(d, s)
		time.Sleep(60 * time.Minute)
	}
}

func StatGet(c *gin.Context) {

	if !IsUpdateRunning {
		go UpdateStat()
		IsUpdateRunning = true
		fmt.Printf("UpdateStat() goroutine started!\n")
	}

	// Allow any origin
	c.Header("Access-Control-Allow-Origin", "*")

	if blacklist.IsBlocked(c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "blocked"})
		return
	}

	c.JSON(http.StatusOK, Current.Copy())
}
