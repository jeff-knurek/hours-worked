package tracking

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// IsUserActive checks loginctl to determine the state of the user
func IsUserActive(username, userid string) (bool, error) {
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

// IsScreenSaverOn checks the state of gnome-screensaver-command for the current user
func IsScreenSaverOn() (bool, error) {
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
