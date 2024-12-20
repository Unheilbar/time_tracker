package cmd

import (
	"log"
	"os"

	"github.com/Unheilbar/time_tracker/internal/flags"
	"github.com/Unheilbar/time_tracker/internal/repository"
	"github.com/Unheilbar/time_tracker/internal/tracker"
	"github.com/spf13/cobra"
)

var app *tracker.App = newApp()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "time_tracker",
	Short: "Time tracker allows you to track time you spend on your activities",
	Run:   app.Root,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start starts timer for task",
	Run:   app.Start,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop stops timer for active task",
	Run:   app.Stop,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List shows all active projects.",
	Run:   app.List,
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Runs last idled task",
	Run:   app.Resume,
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Stops last activated task",
	Run:   app.Remove,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var dataDir = ".time_tracker/badger"
var logDir = ".time_tracker/logs"

const (
	envDataDir = "GO_TIME_TRACKER_DATA_PATH"
	envLogDir  = "GO_TIME_TRACKER_LOG_PATH"
)

func newApp() *tracker.App {
	if val, ok := os.LookupEnv(envDataDir); ok {
		dataDir = val
	}

	if val, ok := os.LookupEnv(envLogDir); ok {
		logDir = val
	}

	db, err := repository.NewBadgerDB(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepo(db)
	//repoJson := repository.NewFileBackend(path.Join(dataDir, "bd.json"))

	return tracker.NewApp(repo)
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(resumeCmd)

	removeCmd.Flags().BoolP(
		flags.All.Name,
		flags.All.Shorthand,
		false,
		"-all to remove all the tasks")
}
