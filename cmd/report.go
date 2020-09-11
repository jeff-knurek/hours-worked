package cmd

import (
	"fmt"
	"os/user"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"hours-worked/pkg/reporting"
)

var output string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runReport(output)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&output, "output", "o", "", "what format to output the report")
}

func runReport(format string) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	filename := viper.GetString("tracking_file")

	fmt.Printf("hours worked this week: %.1f \n", reporting.HoursWorkedThisWeek(filename, user.Username))
	fmt.Println("-------------")

	fmt.Println(reporting.TextCalendar(time.Now(), filename, user.Username))
}
