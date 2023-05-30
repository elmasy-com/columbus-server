package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Date     int64
	Totalnum int64
	m        sync.Mutex
}

var (
	Current Stat
)

func (s *Stat) Update(TotalNum int64) {

	s.m.Lock()
	defer s.m.Unlock()

	s.Date = time.Now().Unix()
	s.Totalnum = TotalNum
}

func (s *Stat) GetDate() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Date
}

func (s *Stat) GetTotalNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Totalnum
}

func (s *Stat) IsEmpty() bool {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Date == 0 && s.Totalnum == 0
}

// UpdateStat is created to run as a goroutine.
// Started in the main.
// Updates the Current variable every 60 minutes and updates the unique collection via db.UpdateUniques() every config.StatAPIWait minutes.
func UpdateStat() {

	ticker := time.NewTicker(6 * time.Hour)

	// Update stats at the beginning
	if total, err := db.GetStat(); err == nil {
		Current.Update(total)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
	}

	for range ticker.C {

		if total, err := db.GetStat(); err == nil {
			Current.Update(total)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
		}
	}
}

func StatGet(c *gin.Context) {

	if Current.IsEmpty() {
		c.Status(http.StatusNoContent)
		return
	}

	// Return a copy only.
	// This was the easiest way to control the write (update) / read process with the mutex.
	c.JSON(http.StatusOK, gin.H{"date": Current.GetDate(), "total": Current.GetTotalNum()})
}
