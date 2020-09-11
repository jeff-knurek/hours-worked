package reporting

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TextCalendar is based on the Markdown calendar from github.com/binzume/go-calendar
// it returns a month view of the tracked hours to be displayed as text/markdown
func TextCalendar(t time.Time, filename, user string) string {
	trackedData := getTrackedData(filename, user)
	weekLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	y := strconv.Itoa(t.Year())
	if trackedData[y] == nil || trackedData[y][t.Month().String()] == nil {
		return "--no data tracked for this month--"
	}
	trackedDays := trackedData[y][t.Month().String()]

	wc := len(weekLabels)
	wd := (int(t.Weekday()) - (t.Day() - 1) + wc*30) % wc
	lastTracked := t.Day()
	last := t.AddDate(0, 1, -t.Day()).Day()

	s := t.Month().String() + " \n\n"
	s += "| " + strings.Join(weekLabels, " | ") + " |\n"
	s += "|" + strings.Repeat("----:|", wc) + "\n"
	s += "|" + strings.Repeat("     |", wd)
	for d := 1; d <= last; d++ {
		if d > lastTracked {
			s += fmt.Sprintf("   . |")
		} else {
			h := float64(trackedDays[strconv.Itoa(d)]) / 60
			switch {
			case h > 10:
				s += fmt.Sprintf("  %.0f |", h)
			case h == 0:
				s += fmt.Sprintf("   0 |")
			default:
				s += fmt.Sprintf(" %.1f |", h)
			}
		}
		wd = (wd + 1) % wc
		if wd == 0 {
			s += "\n"
			if d != last {
				s += "|"
			}
		}
	}
	s += strings.Repeat("     |", wc-wd)

	return s
}
