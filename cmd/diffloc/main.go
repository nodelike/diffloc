package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/nodelike/diffloc/internal/analyzer"
	"github.com/nodelike/diffloc/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	noGitignore    bool
	excludeTests   bool
	customExcludes []string
	allowedExts    []string
	path           string
	cpuProfile     string
	memProfile     string
	staticOutput   bool
	jsonOutput     bool
	maxDepth       int
)

var rootCmd = &cobra.Command{
	Use:   "diffloc",
	Short: "Diff Line Counter - analyze lines of code changes",
	Long: `diffloc is a tool for analyzing lines of code in your projects,
with special support for Git repositories to show changed vs unchanged files.`,
	Version: "1.0.2",
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a directory or Git repository",
	Long: `Analyze files in a directory or Git repository to count lines of code,
showing changes and statistics.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runAnalyze,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .diffloc.yaml)")

	// Add all flags to both root and analyze commands so they work either way
	addAnalyzeFlags := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&noGitignore, "no-gitignore", false, "Ignore .gitignore patterns (always-excluded patterns still apply)")
		cmd.Flags().BoolVar(&excludeTests, "exclude-tests", false, "Exclude test files (_test.go, test/, tests/, *.test.*, *.spec.*)")
		cmd.Flags().StringArrayVar(&customExcludes, "exclude", []string{}, "Additional exclusion pattern (can be repeated)")
		cmd.Flags().StringArrayVar(&allowedExts, "ext", []string{}, "Override allowed file extensions (can be repeated)")
		cmd.Flags().StringVar(&cpuProfile, "profile-cpu", "", "Write CPU profile to file")
		cmd.Flags().StringVar(&memProfile, "profile-mem", "", "Write memory profile to file")
		cmd.Flags().BoolVar(&staticOutput, "static", false, "Print static output without interactive UI")
		cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
		cmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Maximum directory traversal depth (0 = unlimited)")

		// Bind flags to viper
		viper.BindPFlag("no-gitignore", cmd.Flags().Lookup("no-gitignore"))
		viper.BindPFlag("exclude-tests", cmd.Flags().Lookup("exclude-tests"))
		viper.BindPFlag("exclude", cmd.Flags().Lookup("exclude"))
		viper.BindPFlag("ext", cmd.Flags().Lookup("ext"))
		viper.BindPFlag("max-depth", cmd.Flags().Lookup("max-depth"))
	}

	// Add flags to analyze command
	addAnalyzeFlags(analyzeCmd)

	// Add flags to root command for default behavior
	addAnalyzeFlags(rootCmd)

	rootCmd.AddCommand(analyzeCmd)

	// Make analyze the default command
	rootCmd.Run = analyzeCmd.Run
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".diffloc")
	}

	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func runAnalyze(cmd *cobra.Command, args []string) {
	// Get path from args or use current directory
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "."
	}

	// Get working directory if path is relative
	if path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Safety check: validate the path is safe to analyze
	if err := analyzer.ValidatePath(path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Warn about potentially large directories
	if shouldWarn, warning := analyzer.ShouldWarnLargeDirectory(path); shouldWarn {
		fmt.Fprintln(os.Stderr, warning)
	}

	// Start CPU profiling if requested
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	// Merge flags with config file values
	if !cmd.Flags().Changed("no-gitignore") {
		noGitignore = viper.GetBool("no-gitignore")
	}
	if !cmd.Flags().Changed("exclude-tests") {
		excludeTests = viper.GetBool("exclude-tests")
	}
	if len(customExcludes) == 0 {
		customExcludes = viper.GetStringSlice("exclude")
	}
	if len(allowedExts) == 0 {
		allowedExts = viper.GetStringSlice("ext")
	}
	if maxDepth == 0 {
		maxDepth = viper.GetInt("max-depth")
	}

	// Create filter
	filter := analyzer.NewFilter(allowedExts, customExcludes, !noGitignore, excludeTests)

	// Load gitignore if in a git repo and respecting gitignore
	if !noGitignore && analyzer.IsGitRepo(path) {
		repoRoot, err := analyzer.GetRepoRoot(path)
		if err == nil {
			filter.LoadGitignore(repoRoot)
		}
	}

	// Create context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, canceling...")
		cancel()
	}()

	// Analyze the directory
	stats, err := analyzer.Analyze(ctx, path, filter)
	if err != nil {
		if err == context.Canceled {
			fmt.Fprintln(os.Stderr, "\nAnalysis canceled by user")
			os.Exit(130) // Standard exit code for SIGINT
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output results
	if jsonOutput {
		// JSON output
		output, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else if staticOutput {
		// Static output (same format as TUI, but no interactivity)
		if err := ui.PrintStatic(stats); err != nil {
			fmt.Fprintf(os.Stderr, "Error printing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Interactive TUI
		if err := ui.Run(stats); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	}

	// Write memory profile if requested
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create memory profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write memory profile: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
