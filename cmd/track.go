package cmd

import (
	"fmt"
	"time"

	"os/user"

	"hours-worked/pkg/tracking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showIcon bool

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Track minutes that the user is active",
	Long: `Track minutes that the user is active
Usage of track:
hours-worked track [OPTIONS]

Hoours worked increments the minutes which the user is active in the UI.
Configuration options can be set in the config file, and are created by default on the first run.`,
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

func track() {
	t := time.Now()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	filename := viper.GetString("tracking_file")
	// check active state
	active, err := tracking.IsUserActive(user.Name, user.Uid)
	if err != nil {
		panic(err)
	}
	// check screensaver state
	screensaver, err := tracking.IsScreenSaverOn()
	if err != nil {
		panic(err)
	}
	// increment
	if active && !screensaver {
		minutesToday, err := tracking.RecordActive(filename, user.Username, t)
		if err != nil {
			fmt.Println("error incrementing count:", err)
		}
		hoursActive := fmt.Sprintf("%.1f", float64(minutesToday)/60)
		fmt.Println("current hours active today:", hoursActive)
	}
}

func startTrack() {
	nextTime := time.Now().Truncate(time.Minute)
	for {
		nextTime = nextTime.Add(time.Minute)
		time.Sleep(time.Until(nextTime))
		track()
	}
}
