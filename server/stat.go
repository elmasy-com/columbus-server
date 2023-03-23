package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-sdk/db"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Date      int64
	DomainNum int64
	SubNum    int64
	m         sync.Mutex
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

func (s *Stat) GetDate() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Date
}

func (s *Stat) GetDomainNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.DomainNum
}

func (s *Stat) GetSubNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.SubNum
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

	if Current.GetDate() == 0 && Current.GetDomainNum() == 0 && Current.GetSubNum() == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Return a copy only.
	// This was the easiest way to control the write (update) / read process with the mutex.
	c.JSON(http.StatusOK, gin.H{"date": Current.GetDate(), "domain": Current.GetDomainNum(), "sub": Current.GetSubNum()})
}
