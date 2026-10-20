package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"strings"
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
	jsonOutput     bool
	maxDepth       int
)

// version is set by goreleaser via -ldflags; go install builds fall back to module info.
var version = ""

var rootCmd = &cobra.Command{
	Use:   "diffloc [path]",
	Short: "Diff Line Counter - analyze lines of code changes",
	Long: `diffloc is a tool for analyzing lines of code in your projects,
with special support for Git repositories to show changed vs unchanged files.`,
	Args:    cobra.MaximumNArgs(1),
	Version: resolveVersion(),
	Run:     runAnalyze,
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a directory or Git repository",
	Long: `Analyze files in a directory or Git repository to count lines of code,
showing changes and statistics.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runAnalyze,
}

func resolveVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return "dev"
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .diffloc.yaml in the target path or current directory)")

	// Add all flags to both root and analyze commands so they work either way
	addAnalyzeFlags := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&noGitignore, "no-gitignore", false, "Ignore .gitignore patterns (always-excluded patterns still apply)")
		cmd.Flags().BoolVar(&excludeTests, "exclude-tests", false, "Exclude test files (_test.go, test/, tests/, *.test.*, *.spec.*)")
		cmd.Flags().StringArrayVar(&customExcludes, "exclude", []string{}, "Additional exclusion regex (can be repeated)")
		cmd.Flags().StringArrayVar(&allowedExts, "ext", []string{}, "Override allowed file extensions (can be repeated)")
		cmd.Flags().StringVar(&cpuProfile, "profile-cpu", "", "Write CPU profile to file")
		cmd.Flags().StringVar(&memProfile, "profile-mem", "", "Write memory profile to file")
		cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
		cmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Maximum directory depth (0 = unlimited, 1 = top-level files only)")
	}

	addAnalyzeFlags(analyzeCmd)
	addAnalyzeFlags(rootCmd)

	rootCmd.AddCommand(analyzeCmd)
}

// loadConfig reads .diffloc.yaml from --config, the target path, or the current directory.
func loadConfig(target string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(target)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".diffloc")
	}

	viper.SetEnvPrefix("DIFFLOC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil
		}
		return err
	}
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	return nil
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
	os.Exit(1)
}

func runAnalyze(cmd *cobra.Command, args []string) {
	path = "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fail("failed to resolve path: %v", err)
	}
	path = absPath

	if err := analyzer.ValidatePath(path); err != nil {
		fail("%v", err)
	}

	if shouldWarn, warning := analyzer.ShouldWarnLargeDirectory(path); shouldWarn {
		fmt.Fprintln(os.Stderr, warning)
	}

	if err := loadConfig(path); err != nil {
		fail("failed to read config: %v", err)
	}

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			fail("failed to create CPU profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fail("failed to start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	flags := cmd.Flags()
	if !flags.Changed("no-gitignore") {
		noGitignore = viper.GetBool("no-gitignore")
	}
	if !flags.Changed("exclude-tests") {
		excludeTests = viper.GetBool("exclude-tests")
	}
	if !flags.Changed("exclude") {
		customExcludes = viper.GetStringSlice("exclude")
	}
	if !flags.Changed("ext") {
		allowedExts = viper.GetStringSlice("ext")
	}
	if !flags.Changed("max-depth") {
		maxDepth = viper.GetInt("max-depth")
	}
	if maxDepth < 0 {
		fail("--max-depth must be 0 or greater")
	}

	filter, err := analyzer.NewFilter(analyzer.FilterOptions{
		Exts:             allowedExts,
		Excludes:         customExcludes,
		RespectGitignore: !noGitignore,
		ExcludeTests:     excludeTests,
		MaxDepth:         maxDepth,
	})
	if err != nil {
		fail("%v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stats, err := analyzer.Analyze(ctx, path, filter)
	if err != nil {
		if ctx.Err() != nil {
			fmt.Fprintln(os.Stderr, "\nAnalysis canceled by user")
			os.Exit(130) // Standard exit code for SIGINT
		}
		fail("%v", err)
	}

	for _, p := range stats.Skipped {
		fmt.Fprintf(os.Stderr, "Warning: could not read %s\n", p)
	}

	if jsonOutput {
		output, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fail("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(output))
	} else {
		ui.PrintStatic(stats)
	}

	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			fail("failed to create memory profile: %v", err)
		}
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			fail("failed to write memory profile: %v", err)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
