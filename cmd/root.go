package cmd

import (
	"os"

	"github.com/Unheilbar/time_tracker/internal/tracker"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "time_tracker",
	Short: "Time tracker allows you to track time you spend on your activities",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: tracker.Root,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start starts timer for task",
	Run:   tracker.Start,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop stops timer for active task",
	Run:   tracker.Stop,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List shows all active projects.",
	Run:   tracker.List,
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Runs last idled task",
	Run:   tracker.Resume,
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Stops last activated task",
	Run:   tracker.Remove,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(resumeCmd)

	// create flags
	// TODO move them to internal
	startCmd.Flags().StringP(
		"task",
		"t",
		"",
		`
		name of the task to start
		`)

	startCmd.MarkFlagsOneRequired("task")

	removeCmd.Flags().StringP(
		"task",
		"t",
		"",
		`
		name of the task to remove
		`)
	removeCmd.MarkFlagsOneRequired("task")
}
