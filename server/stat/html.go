package stat

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/elmasy-com/columbus-server/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ctLog struct {
	Name            string
	Index           string
	Size            string
	Remaining       string
	Complete        string
	CompletePercent float64
}

type historyStat struct {
	Num           string
	Date          string
	Total         string
	Updated       string
	Valid         string
	CTLogTotalInt int64
	CTLogTotal    string
}

type statistics struct {
	Date           string
	Total          string
	Updated        string
	UpdatedPercent string
	Valid          string
	ValidPercent   string
	CTTotalInt     int64
	CTTotal        string
	CTLogs         []ctLog
	History        []historyStat
}

//go:embed stat.html
var statHtml string

func parseStatistic() (statistics, error) {

	s, err := db.StatisticsGets()
	if err != nil {
		return statistics{}, fmt.Errorf("failed to get newset statistic: %w", err)
	}

	printer := message.NewPrinter(language.English)

	var stat statistics

	// The first element in the slice is the newest entry
	stat.Date = time.Unix(s[0].Date, 0).String()
	stat.Total = printer.Sprint(s[0].Total)
	stat.Updated = printer.Sprint(s[0].Updated)
	stat.UpdatedPercent = fmt.Sprintf("%.2f%%", float64(s[0].Updated)/float64(s[0].Total)*100)
	stat.Valid = printer.Sprint(s[0].Valid)
	stat.ValidPercent = fmt.Sprintf("%.2f%%", float64(s[0].Valid)/float64(s[0].Total)*100)

	stat.CTLogs = make([]ctLog, len(s[0].CTLogs))

	for i := range s[0].CTLogs {
		stat.CTLogs[i].Name = s[0].CTLogs[i].Name
		stat.CTLogs[i].Index = printer.Sprint(s[0].CTLogs[i].Index)
		stat.CTLogs[i].Size = printer.Sprint(s[0].CTLogs[i].Size)
		stat.CTLogs[i].Remaining = printer.Sprint(s[0].CTLogs[i].Size - s[0].CTLogs[i].Index)
		stat.CTLogs[i].CompletePercent = float64(s[0].CTLogs[i].Index) / float64(s[0].CTLogs[i].Size) * 100
		stat.CTLogs[i].Complete = fmt.Sprintf("%7.2f%%", stat.CTLogs[i].CompletePercent)

		stat.CTTotalInt += s[0].CTLogs[i].Size
	}

	stat.CTTotal = printer.Sprint(stat.CTTotalInt)

	sort.Slice(stat.CTLogs, func(i, j int) bool { return stat.CTLogs[i].CompletePercent > stat.CTLogs[j].CompletePercent })

	// The remaining elements are the history
	hs := s[1:]

	hNum := 1

	stat.History = make([]historyStat, len(hs))

	for i := range hs {

		stat.History[i].Num = printer.Sprint(hNum)
		stat.History[i].Date = time.Unix(hs[i].Date, 0).String()

		stat.History[i].Total = printer.Sprint(hs[i].Total)
		stat.History[i].Updated = printer.Sprint(hs[i].Updated)
		stat.History[i].Valid = printer.Sprint(hs[i].Valid)

		for ii := range hs[i].CTLogs {
			stat.History[i].CTLogTotalInt += hs[i].CTLogs[ii].Size
		}

		stat.History[i].CTLogTotal = printer.Sprint(stat.History[i].CTLogTotalInt)

		hNum++
	}

	return stat, nil
}

func GetStat(c *gin.Context) {

	stat, err := parseStatistic()
	if err != nil {
		c.Error(fmt.Errorf("failed to parse statistic: %w", err))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	t := template.New("stat")

	t, err = t.Parse(statHtml)
	if err != nil {
		c.Error(fmt.Errorf("failed to parse template html: %w", err))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	err = t.Execute(c.Writer, stat)
	if err != nil {
		c.Error(fmt.Errorf("failed to execute template: %w", err))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

}
