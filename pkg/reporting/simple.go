package reporting

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

// TODO - this structure is duplicated - shouled be shared
type userData map[string]years
type years map[string]months
type months map[string]days
type days map[string]int

func getTrackedData(filename, user string) years {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	// unmarshall data from the file
	var obj userData
	err = json.Unmarshal(data, &obj)
	if err != nil {
		// possible that the file doesn't exist, or the user hasn't had any data collected yet
		panic(err)
	}
	return obj[user]
}

// HoursWorkedThisWeek returns the total hours documented per day since Sunday
// (ie work week starts Monday at 00:01)
func HoursWorkedThisWeek(filename, user string) float64 {
	data := getTrackedData(filename, user)
	now := time.Now()

	return sumThisWeek(data, now)
}

func sumThisWeek(data years, t time.Time) float64 {
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
