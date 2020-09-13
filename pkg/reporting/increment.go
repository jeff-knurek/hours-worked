package reporting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// UserData defines the full object of tracking data with key:username
type UserData map[string]Years

// Years defines the object of tracking data for the user, with key:year
type Years map[string]Months

// Months defines the object of tracking data for the user, with key:month
type Months map[string]Days

// Days defines the object of tracking data for the user, with key:day, and value[int] of minutes active
type Days map[string]int

// RecordActive increments the count of activity for the user provided on the datetime provided
// return the current minutes active for today
func RecordActive(filename, user string, currentTime time.Time) (int, error) {
	checkFileExists(filename, user)
	data := getTrackedData(filename)
	u := data[user]

	// increment counter
	u, curCount := incrementTime(u, currentTime)
	// update the file
	data[user] = u
	updateFile(filename, data)
	return curCount, nil
}

func incrementTime(data Years, t time.Time) (Years, int) {
	y := strconv.Itoa(t.Year())
	m := t.Month().String() // uses month name - might "break" if user switches locales
	d := strconv.Itoa(t.Day())
	cur := 1
	if data[y] != nil {
		if data[y][m] != nil {
			cur = data[y][m][d] + 1
			data[y][m][d] = cur
		} else {
			data[y][m] = Days{d: 1}
		}
	} else {
		data[y] = Months{m: {d: 1}}
	}
	return data, cur
}

func checkFileExists(filename, user string) {
	if !fileExists(filename) {
		// write empty obj to it
		obj := UserData{user: {}}
		file, _ := json.MarshalIndent(obj, "", "    ")
		err := ioutil.WriteFile(filename, file, 0644)
		if err != nil {
			panic(fmt.Errorf("not able to write to the new file: %s", err))
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getTrackedData(filename string) UserData {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	// unmarshall data from the file
	var obj UserData
	err = json.Unmarshal(data, &obj)
	if err != nil {
		// possible that the file doesn't exist, or the user hasn't had any data collected yet
		// todo, return error
		panic(err)
	}
	return obj
}

func updateFile(filename string, data UserData) error {
	file, _ := json.MarshalIndent(data, "", "    ")
	err := ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
