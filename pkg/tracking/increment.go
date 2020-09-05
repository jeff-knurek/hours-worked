package tracking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type userData map[string]years
type years map[string]months
type months map[string]days
type days map[string]int

// RecordActive increments the count of activity for the user provided on the datetime provided
func RecordActive(user string, currentTime time.Time) {
	//TODO: Set the info about file elsewhere
	filename := "/home/jeff/.hours-worked/tracked.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// unmarshall data from the file
	var obj userData
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	// increment counter
	obj[user] = incrementTime(obj[user], currentTime)
	// update the file
	file, _ := json.MarshalIndent(obj, "", "    ")
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		fmt.Println("failed to update file:", err)
	}
}

func incrementTime(data years, t time.Time) years {
	y := strconv.Itoa(t.Year())
	m := t.Month().String() // uses month name - might "break" if user switches locales
	d := strconv.Itoa(t.Day())
	if data[y] != nil {
		if data[y][m] != nil {
			data[y][m][d]++
		} else {
			data[y][m] = days{d: 1}
		}
	} else {
		data[y] = months{m: {d: 1}}
	}
	return data
}
