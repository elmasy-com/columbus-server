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
}

//go:embed stat.html
var statHtml string

func GetStat(c *gin.Context) {

	s, err := db.StatisticsGetNewest()
	if err != nil {
		c.Error(fmt.Errorf("failed to get newset statistic: %w", err))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	printer := message.NewPrinter(language.English)

	var stat statistics

	stat.Date = time.Unix(s.Date, 0).String()
	stat.Total = printer.Sprint(s.Total)
	stat.Updated = printer.Sprint(s.Updated)
	stat.UpdatedPercent = fmt.Sprintf("%.2f%%", float64(s.Updated)/float64(s.Total)*100)
	stat.Valid = printer.Sprint(s.Valid)
	stat.ValidPercent = fmt.Sprintf("%.2f%%", float64(s.Valid)/float64(s.Total)*100)

	stat.CTLogs = make([]ctLog, len(s.CTLogs))

	for i := range s.CTLogs {
		stat.CTLogs[i].Name = s.CTLogs[i].Name
		stat.CTLogs[i].Index = printer.Sprint(s.CTLogs[i].Index)
		stat.CTLogs[i].Size = printer.Sprint(s.CTLogs[i].Size)
		stat.CTLogs[i].Remaining = printer.Sprint(s.CTLogs[i].Size - s.CTLogs[i].Index)
		stat.CTLogs[i].CompletePercent = float64(s.CTLogs[i].Index) / float64(s.CTLogs[i].Size) * 100
		stat.CTLogs[i].Complete = fmt.Sprintf("%7.2f%%", stat.CTLogs[i].CompletePercent)

		stat.CTTotalInt += s.CTLogs[i].Size
	}

	stat.CTTotal = printer.Sprint(stat.CTTotalInt)

	sort.Slice(stat.CTLogs, func(i, j int) bool { return stat.CTLogs[i].CompletePercent > stat.CTLogs[j].CompletePercent })

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
