package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Sannrox/revar/internal/debug"
	"github.com/spf13/cobra"
)

var (
	PlatformName = ""
	Version      = "unknown-version"
	GitCommit    = "unknown-commit"
	BuildTime    = "unknown-buildtime"
	BuildArch    = "unknown-buildarch"
	BuildOs      = "unknown-buildos"
)

type options struct {
	// debug mode
	Debug bool

	// recursive mode
	Recursive bool

	// dry run mode
	DryRun bool

	// verbose mode
	Verbose bool

	// file to read
	File string

	// directory to read
	Dir string
}

func NewRevarCommand() *cobra.Command {
	opts := &options{}
	cmd := &cobra.Command{
		Use:     "revar [flags] [regex] [replacement]",
		Short:   "revar is a tool to replace variables in files",
		Long:    "revar is a tool to replace variables in files",
		Version: fmt.Sprintf("%s-%s-%s-%s-%s", Version, GitCommit, BuildTime, BuildArch, BuildOs),
		Args:    cobra.MinimumNArgs(2),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.Debug {
				debug.Enable()
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			regex := args[0]
			replacement := args[1]
			return GoRevar(opts, regex, replacement)
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "debug mode")
	cmd.PersistentFlags().BoolVarP(&opts.Recursive, "recursive", "r", false, "recursive mode")
	cmd.PersistentFlags().BoolVarP(&opts.DryRun, "dry-run", "n", false, "dry run mode")
	cmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose mode")
	cmd.PersistentFlags().StringVarP(&opts.File, "file", "f", "", "file to read")
	cmd.PersistentFlags().StringVarP(&opts.Dir, "dir", "d", "", "directory to read")

	return cmd

}

func main() {
	cmd := NewRevarCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func GoRevar(opts *options, regex string, replacement string) error {
	fmt.Fprintf(os.Stdout, "Let's replace %s with %s\n", regex, replacement)

	re, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	file, err := os.OpenFile(opts.File, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	for lineNumber, line := range lines {
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			boldMatch := fmt.Sprintf("\033[1m%s\033[0m", match)
			text := strings.Replace(line, match, boldMatch, 1)
			fmt.Printf("%s +%d: %s\n", opts.File, lineNumber, text)

			if !opts.DryRun {
				lines[lineNumber] = re.ReplaceAllString(line, replacement)

			}
		}
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Done!\n")
	return nil
}
