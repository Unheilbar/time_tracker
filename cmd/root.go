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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Start starts timer for task",
	Run:   tracker.Create,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop stops timer for task",
	Run:   tracker.Stop,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List shows all active projects.",
	Run:   tracker.List,
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Runs last stopped task",
	Run:   tracker.Resume,
}

var afkCmd = &cobra.Command{
	Use:   "afk",
	Short: "Stops last activated task",
	Run:   tracker.AFK,
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
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(listCmd)

	// create flags
	// TODO move them to internal
	createCmd.Flags().StringP(
		"task",
		"t",
		"",
		`
		name of the task to create
		`)
	createCmd.Flags().StringSliceP(
		"description",
		"d",
		nil,
		`
		description of current task session
		if not provided is empty
		actual shortland is -desc because of bug in cobra library
		`)

	createCmd.MarkFlagsOneRequired("task")

	stopCmd.Flags().StringP(
		"task",
		"t",
		"",
		`
		name of the task to stop
		`)
}
