package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const dFormat = "2006-01-02"

// dayOffCmd represents the dayOff command
var dayOffCmd = &cobra.Command{
	Use:   "dayOff [day]",
	Short: "Add a day that isn't calcualted for total. (format: 2016-04-28)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a day (format: 2016-04-28)")
		}
		if validDay(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid day format (yyyy-mm-dd): %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		addDayOff(args[0])
	},
}

func init() {
	rootCmd.AddCommand(dayOffCmd)
}

func validDay(day string) bool {
	_, err := time.Parse(dFormat, day)
	if err != nil {
		return false
	}
	return true
}

func addDayOff(day string) {
	t, _ := time.Parse(dFormat, day)
	newDay := t.Format(dFormat)
	days := viper.GetStringSlice("days_off")
	for _, d := range days {
		if d == newDay {
			fmt.Println(newDay, "already exists")
			return
		}
	}
	days = append(days, newDay)

	viper.Set("days_off", days)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("problem with updating config file", err)
	}
	return
}
