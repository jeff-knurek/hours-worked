package cmd

import (
	"time"

	"os/user"

	"hours-worked/pkg/tracking"

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

func track() {
	t := time.Now()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
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
		tracking.RecordActive(user.Username, t)
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
