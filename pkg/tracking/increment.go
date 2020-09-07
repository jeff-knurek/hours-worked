package tracking

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

type userData map[string]years
type years map[string]months
type months map[string]days
type days map[string]int

// RecordActive increments the count of activity for the user provided on the datetime provided
// return the current minutes active for today
func RecordActive(filename, user string, currentTime time.Time) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	// unmarshall data from the file
	var obj userData
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return 0, err
	}
	// increment counter
	u, curCount := incrementTime(obj[user], currentTime)
	// update the file
	obj[user] = u
	file, _ := json.MarshalIndent(obj, "", "    ")
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return 0, err
	}
	return curCount, nil
}

func incrementTime(data years, t time.Time) (years, int) {
	y := strconv.Itoa(t.Year())
	m := t.Month().String() // uses month name - might "break" if user switches locales
	d := strconv.Itoa(t.Day())
	cur := 1
	if data[y] != nil {
		if data[y][m] != nil {
			cur = data[y][m][d] + 1
			data[y][m][d] = cur
		} else {
			data[y][m] = days{d: 1}
		}
	} else {
		data[y] = months{m: {d: 1}}
	}
	return data, cur
}
