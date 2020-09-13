package reporting

import (
	"strconv"
	"time"
)

// HoursWorkedThisWeek returns the total hours documented per day since Sunday
// (ie work week starts Monday at 00:01)
func HoursWorkedThisWeek(filename, user string) float64 {
	data := getTrackedData(filename)[user]
	now := time.Now()

	return sumThisWeek(data, now)
}

func sumThisWeek(data Years, t time.Time) float64 {
	thisDay := t
	sum := 0
	for i := int(t.Weekday()); i > 0; i-- {
		y := strconv.Itoa(thisDay.Year())
		m := thisDay.Month().String() // uses month name - might "break" if user switches locales
		d := strconv.Itoa(thisDay.Day())
		if data[y] != nil {
			if data[y][m] != nil {
				sum += data[y][m][d]
			}
		}
		thisDay = thisDay.AddDate(0, 0, -1)
	}
	return float64(sum) / 60
}
