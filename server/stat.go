package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Date     int64              `json:"date"`
	Total    int64              `json:"total"`
	Updated  int64              `json:"updated"`
	Valid    int64              `json:"valid"`
	Scanners []db.ScannerSchema `json:"scanners"`
	m        sync.Mutex         `json:"-"`
}

var (
	Current Stat
)

func (s *Stat) Update(total, updated, valid int64, scanner []db.ScannerSchema) {

	s.m.Lock()
	defer s.m.Unlock()

	s.Date = time.Now().Unix()
	s.Total = total
	s.Updated = updated
	s.Valid = valid
	s.Scanners = scanner
}

func (s *Stat) Get() Stat {
	s.m.Lock()
	defer s.m.Unlock()

	return Stat{Date: s.Date, Total: s.Total, Updated: s.Updated, Valid: s.Valid, Scanners: s.Scanners}
}

func (s *Stat) IsEmpty() bool {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Date == 0 || s.Total == 0 || s.Valid == 0 || len(s.Scanners) == 0
}

// UpdateStat is created to run as a goroutine.
// Started in the main.
// Updates the Current variable every 60 minutes and updates the unique collection via db.UpdateUniques() every config.StatAPIWait minutes.
func UpdateStat() {

	// Update stats at the beginning
	total, updated, valid, scanners, err := db.GetStat()
	if err == nil {
		Current.Update(total, updated, valid, scanners)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
	}

	for {

		time.Sleep(time.Duration(rand.Int63n(7200)+7200) * time.Second)

		if total, updated, valid, scanners, err := db.GetStat(); err == nil {
			Current.Update(total, updated, valid, scanners)
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
	c.JSON(http.StatusOK, Current.Get())
}
