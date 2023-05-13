package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elmasy-com/columbus-sdk/db"
	"github.com/elmasy-com/columbus-server/config"
	"github.com/gin-gonic/gin"
)

type Stat struct {
	Date          int64
	Totalnum      int64
	TldNum        int64
	DomainNum     int64
	FullDomainNum int64
	SubNum        int64
	m             sync.Mutex
}

var (
	Current Stat
)

func (s *Stat) Update(TotalNum int64, TldNum int64, DomainNum int64, FullDomainNum int64, SubNum int64) {

	s.m.Lock()
	defer s.m.Unlock()

	s.Date = time.Now().Unix()
	s.Totalnum = TotalNum
	s.TldNum = TldNum
	s.DomainNum = DomainNum
	s.FullDomainNum = FullDomainNum
	s.SubNum = SubNum
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

func (s *Stat) GetTldNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.TldNum
}

func (s *Stat) GetDomainNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.DomainNum
}

func (s *Stat) GetFullDomainNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.FullDomainNum
}
func (s *Stat) GetSubNum() int64 {

	s.m.Lock()
	defer s.m.Unlock()

	return s.SubNum
}

func (s *Stat) IsEmpty() bool {

	s.m.Lock()
	defer s.m.Unlock()

	return s.Date == 0 && s.Totalnum == 0 && s.TldNum == 0 &&
		s.DomainNum == 0 && s.FullDomainNum == 0 && s.SubNum == 0
}

// UpdateStat is created to run as a goroutine.
// Started in the main.
// Updates the Current variable every 60 minutes and updates the unique collection via db.UpdateUniques() every config.StatAPIWait minutes.
func UpdateStat() {

	updateUniques := time.NewTicker(time.Duration(config.StatAPIWait) * time.Minute)
	getStat := time.NewTicker(60 * time.Minute)

	// Update stats first
	if total, tlds, domains, fullDomains, subs, err := db.GetStat(); err == nil {
		Current.Update(total, tlds, domains, fullDomains, subs)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
	}

	for {

		select {
		case <-updateUniques.C:

			err := db.UpdateUniques()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update uniques: %s\n", err)
			}

		case <-getStat.C:

			total, tlds, domains, fullDomains, subs, err := db.GetStat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get DB stat: %s\n", err)
				continue
			}

			Current.Update(total, tlds, domains, fullDomains, subs)
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
	c.JSON(http.StatusOK,
		gin.H{"date": Current.GetDate(), "total": Current.GetTotalNum(), "tld": Current.GetTldNum(),
			"domain": Current.GetDomainNum(), "fulldomain": Current.GetFullDomainNum(), "sub": Current.GetSubNum()})
}
