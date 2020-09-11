package cmd

import (
	"fmt"
	"time"

	"os/user"

	"hours-worked/pkg/reporting"
	"hours-worked/pkg/tracking"

	"github.com/getlantern/systray"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showIcon bool
var iconChoice string

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
		if showIcon {
			go systray.Run(startMenuItem, nil)
		}
		startTrack()
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)

	trackCmd.Flags().BoolVarP(&showIcon, "icon", "i", false, "add icon in menu bar that increments")
	trackCmd.Flags().StringVar(&iconChoice, "icon-color", "light", "show a white or dark icon")
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
		if showIcon {
			hoursActive := fmt.Sprintf("%.1f", float64(minutesToday)/60)
			systray.SetTitle(hoursActive)
		}
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

func startMenuItem() {
	systray.SetTitle("starting up...")
	systray.SetTooltip("hours worked today")
	switch iconChoice {
	case "light":
		systray.SetTemplateIcon(nil, reporting.HourGlassIconLight)
	case "dark":
		systray.SetTemplateIcon(nil, reporting.HourGlassIconDark)
	default:
		systray.SetTemplateIcon(nil, reporting.HourGlassIconLight)
	}
	systray.AddMenuItem("hours active today", "hours active today")
}
