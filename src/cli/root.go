package cli

import (
	"fmt"
	"os"

	"github.com/hhiroshell/git-client-tui/src/tui"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new root command for the Git TUI application
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitui",
		Short: "Git Client TUI - A terminal user interface for Git",
		Long: `A visually friendly TUI git client built with Go and Cobra CLI library 
that provides granular staging control (file/hunk/line level), intuitive branch
management with auto-fetch, commit operations including squashing and message editing, 
and optimized UX for developer workflows.`,
		Run: runTUI,
	}

	cmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	cmd.AddCommand(newVersionCommand())

	return cmd
}

// Execute executes the root command
func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runTUI starts the TUI application
func runTUI(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if err := tui.Start(debug); err != nil {
		fmt.Printf("Error starting TUI: %v\n", err)
		os.Exit(1)
	}
}

// newVersionCommand creates a version command
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Git Client TUI v0.1.0")
		},
	}
}