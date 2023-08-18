package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

	// interactive mode
	Interactive bool

	// force mode
	Force bool
}

func NewRevarCommand() *cobra.Command {
	opts := &options{}
	cmd := &cobra.Command{
		Use:     "revar [flags] [regex] [replacement] [file|dir]",
		Short:   "revar is a tool to replace variables in files",
		Long:    "revar is a tool to replace variables in files",
		Version: fmt.Sprintf("%s-%s-%s-%s-%s", Version, GitCommit, BuildTime, BuildArch, BuildOs),
		Args:    cobra.MinimumNArgs(3),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.Debug {
				debug.Enable()
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			regex := args[0]
			replacement := args[1]
			p := args[2]
			return GoRevar(opts, regex, replacement, p)
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "debug mode")
	cmd.PersistentFlags().BoolVarP(&opts.Recursive, "recursive", "r", false, "recursive mode")
	cmd.PersistentFlags().BoolVarP(&opts.DryRun, "dry-run", "n", false, "dry run mode")
	cmd.PersistentFlags().BoolVarP(&opts.Interactive, "interactive", "i", false, "interactive mode")
	cmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose mode")
	cmd.PersistentFlags().BoolVarP(&opts.Force, "force", "f", false, "force mode")

	cmd.MarkFlagsMutuallyExclusive("dry-run", "interactive")
	cmd.MarkFlagsMutuallyExclusive("dry-run", "force")

	return cmd

}

func main() {
	cmd := NewRevarCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func GoRevar(opts *options, regex string, replacement string, p string) error {
	re, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	fileInfo, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("failed to get info: %w", err)
	}
	if opts.DryRun {
		fmt.Printf("Dry run mode enabled. No changes will be made.\n")
	}
	if fileInfo.IsDir() {
		var allFiles []string
		switch opts.Recursive {
		case true:

			allFiles, err = LoadFilesRecursive(p)
			if err != nil {
				return fmt.Errorf("failed to load files recursively: %w", err)
			}
		case false:
            allFiles, err = LoadFilesFromDir(p)
            if err != nil {
                return fmt.Errorf("failed to load files from dir: %w", err)
            }
		}
		fmt.Printf("Found %d files\n", len(allFiles))
		fmt.Printf("Replacing...\n")
		fmt.Printf("========================================\n")
		for _, file := range allFiles {
			if err := GoRevarFile(opts, re, replacement, file); err != nil {
				return fmt.Errorf("failed to revar file: %w", err)
			}
		}
	} else {
		if opts.Recursive {
			return fmt.Errorf("recursive mode is only available for directories")
		}
		if err := GoRevarFile(opts, re, replacement, p); err != nil {
			return fmt.Errorf("failed to revar file: %w", err)
		}

	}

	fmt.Fprintf(os.Stdout, "Done!\n")
	return nil
}

func LoadFilesRecursive(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failed to walk dir: %v\n", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func LoadFilesFromDir(dir string) ([]string, error) {
    var files []string
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    for _, file := range entries {
        if file.IsDir() {
            continue
        }
        files = append(files, filepath.Join(dir, file.Name()))
    }
    return files, nil
}

func GoRevarFile(opts *options, re *regexp.Regexp, replacement string, p string) error {
	file, err := os.OpenFile(p, os.O_RDWR, 0644)
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


	if opts.Interactive && opts.Force {
		opts.Verbose = true
	}

	for lineNumber, line := range lines {
		if re.MatchString(line) {
			if opts.Verbose || opts.DryRun {
				DisplayMatchedStringsInLine(file.Name(), lineNumber+1, line, re, replacement)
			}
			if opts.Interactive && !opts.Force {
				lines[lineNumber] = InteractiveReplacement(file.Name(), lineNumber+1, line, re, replacement)

			} else {
				lines[lineNumber] = ReplaceAllMatchedStringsInLine(line, re, replacement)
			}
		}
	}
	if opts.Interactive && opts.Force {
		fmt.Printf("Do you want to replace all this lines? [y/n]: ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" {
			opts.DryRun = false
		} else {
			opts.DryRun = true
		}
	}

	if !opts.DryRun {

		if err := OverWriteFileContent(file, strings.Join(lines, "\n")); err != nil {
			return fmt.Errorf("failed to overwrite file: %w", err)
		}

	}

	return nil
}

func ReplaceAllMatchedStringsInLine(line string, re *regexp.Regexp, replacement string) string {
	matches := re.FindAllString(line, -1)
	for _, match := range matches {
		line = strings.Replace(line, match, replacement, 1)
	}
	return line
}

func DisplayMatchedStringsInLine(fileName string, lineNumber int, line string, re *regexp.Regexp, replacement string) {
	matches := re.FindAllString(line, -1)
	boldReplacement := fmt.Sprintf("\033[1m%s\033[0m", replacement)
	strikethroughText := fmt.Sprintf("\033[9m%s\033[0m", matches[0])

	line = strings.Replace(line, matches[0], strikethroughText+"/"+boldReplacement, len(matches))

	fmt.Printf("%s +%d: %s\n", fileName, lineNumber, line)
}

func InteractiveReplacement(fileName string, fileNumber int, line string, re *regexp.Regexp, replacement string) string {
	matches := re.FindAllString(line, -1)
	for matchIndex, match := range matches {
		DisplaySingleMatchedStringInLine(fileName, fileNumber, line, match, replacement, matchIndex)
		fmt.Printf("Do you want to replace this line? [y/n]: ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" {
			line = strings.Replace(line, match, replacement, matchIndex+1)
		}
	}
	return line
}

func DisplaySingleMatchedStringInLine(fileName string, lineNumber int, line, match, replacement string, matchIndex int) {
	boldReplacement := fmt.Sprintf("\033[1m%s\033[0m", replacement)
	strikethroughText := fmt.Sprintf("\033[9m%s\033[0m", match)

	line = strings.Replace(line, match, strikethroughText+"/"+boldReplacement, matchIndex+1)

	fmt.Printf("%s +%d: %s\n", fileName, lineNumber, line)
}

func OverWriteFileContent(file *os.File, content string) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

