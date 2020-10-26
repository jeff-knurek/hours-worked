package reporting

import (
	"fmt"
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

// HoursWorkedThisMonth returns the total hours documented per day for the current month
func HoursWorkedThisMonth(filename, user string, daysOff []string) float64 {
	data := getTrackedData(filename)[user]
	now := time.Now()

	return sumThisMonth(data, now, daysOff)
}

func sumThisMonth(data Years, t time.Time, daysOff []string) float64 {
	thisDay := t
	y := strconv.Itoa(thisDay.Year())
	m := thisDay.Month().String() // uses month name - might "break" if user switches locales
	days := data[y][m]

	sum := 0
	for key, d := range days {
		day, _ := strconv.Atoi(key)
		if !isDayOff(daysOff, fmt.Sprintf("%s-%02d-%02d", y, int(thisDay.Month()), day)) {
			sum += d
		}
	}
	return float64(sum) / 60
}

// AvailableDaysThisMonth returns the total days that could be working days
func AvailableDaysThisMonth(t time.Time, daysOff []string) int {
	totalDays := 0
	thisDay := t
	lastDay := thisDay.Day()
	// lastDayMonth := time.Date(thisDay.Year(), time.Month(int(thisDay.Month())+1), 0, 0, 0, 0, 0, time.UTC)
	// lastDay := lastDayMonth.Day()

	for i := 1; i <= lastDay; i++ {
		currDay := time.Date(thisDay.Year(), time.Month(int(thisDay.Month())), i, 0, 0, 0, 0, time.UTC)
		if isDayOff(daysOff, currDay.Format("2006-01-02")) {
			continue
		}
		if int(currDay.Weekday()) == 0 || int(currDay.Weekday()) == 6 {
			continue
		}
		totalDays++
	}

	return totalDays
}

func isDayOff(daysOff []string, day string) bool {
	for _, d := range daysOff {
		if d == day {
			return true
		}
	}
	return false
}
