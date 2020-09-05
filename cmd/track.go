package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"os/exec"
	"os/user"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use: "track",
	// TODO: update the descriptions
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTrack()
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startTrack() error {
	// get interval
	interval := 60
	fmt.Printf("track called with interval: %d seconds \n", interval)

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// fmt.Println("Hi " + user.Name + " (id: " + user.Uid + ")")

	fmt.Println("---------")
	active, err := isUserActive(user.Name, user.Uid)
	if err != nil {
		panic(err)
	}

	screensaver, err := isScreenSaverOn()
	if err != nil {
		panic(err)
	}

	// fmt.Println("---------")
	if active && !screensaver {
		fmt.Println("user is active")
		recordActive(user.Username, time.Now())
	}
	return nil
}

// -------------------
// -------------------
// -------------------

func isUserActive(username, userid string) (bool, error) {
	sOut, err := exec.Command("loginctl", "show-user", userid, "-p", "Sessions").Output()
	if err != nil {
		return false, err
	}
	sIDs, err := parseSessionID(string(sOut))
	if err != nil {
		return false, err
	}

	active := false
	for _, sid := range sIDs {
		aOut, err := exec.Command("loginctl", "show-session", "-p", "State", sid).Output()
		if err != nil {
			return false, err
		}
		active, err = parseActive(string(aOut))
		if err != nil {
			return false, err
		}
		if active {
			return true, nil
		}

	}
	fmt.Println(active)
	return active, nil
}

func parseSessionID(sessionOutput string) ([]string, error) {
	match := regexp.MustCompile("Sessions=(.*)").FindStringSubmatch(sessionOutput)
	if len(match) > 1 && match[1] != "" {
		return strings.Split(match[1], " "), nil
	}
	return nil, fmt.Errorf("session id not found in: %s", sessionOutput)
}

func parseActive(activeOutput string) (bool, error) {
	match := regexp.MustCompile("State=(.*)").FindStringSubmatch(activeOutput)
	if len(match) != 2 {
		return false, fmt.Errorf("state id not found: %s", activeOutput)
	}
	if match[1] == "active" {
		return true, nil
	}
	return false, nil
}

func isScreenSaverOn() (bool, error) {
	out, err := exec.Command("gnome-screensaver-command", "--query").Output()
	if err != nil {
		return false, err
	}
	return parseScreenSaver(string(out))
}

func parseScreenSaver(cmdOut string) (bool, error) {
	match := regexp.MustCompile("The screensaver is (.*)").FindStringSubmatch(cmdOut)
	if len(match) != 2 {
		return false, fmt.Errorf("screensaver info not found: %s", cmdOut)
	}
	if match[1] == "active" {
		return true, nil
	}
	return false, nil
}

// -------------------
// -------------------
// -------------------

type userData map[string]years
type years map[string]months
type months map[string]days
type days map[string]int

func recordActive(user string, currentTime time.Time) {
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
