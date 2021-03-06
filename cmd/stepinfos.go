// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/readerutil"
	"github.com/spf13/cobra"
)

var (
	isTimeOnlyModeFlag *bool
)

// stepinfosCmd represents the stepinfos command
var stepinfosCmd = &cobra.Command{
	Use:   "stepinfos BITRISE-LOG-FILE-PATH",
	Short: "Filter only step infos from the bitrise log",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("stepinfos:")
		if len(args) < 1 {
			return errors.New("No log file path specified")
		}
		logFilePath := args[0]

		isTimeOnlyMode := false
		if isTimeOnlyModeFlag != nil {
			isTimeOnlyMode = *isTimeOnlyModeFlag
		}
		fmt.Println(" * isTimeOnlyMode: ", isTimeOnlyMode)
		fmt.Println()

		return filterStepInfosFromLogFile(logFilePath, isTimeOnlyMode)
	},
}

func filterStepInfosFromLogFile(logFilePath string, isTimeOnlyMode bool) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read log file (%s), error: %s", logFilePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(" [!] Failed to close file, error: ", err)
		}
	}()

	err = readerutil.WalkLines(file, func(line string) error {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 2 {
			switch trimmedLine[0:2] {
			case "+-", "| ":
				if isTimeOnlyMode {
					pattern := `(?i)^\| .+\| [0-9.]+ sec[[:space:]]*\|$`
					if isMatch, err := regexp.MatchString(pattern, trimmedLine); err != nil {
						return fmt.Errorf("Failed to match line (%s) with regex (%s), error: %s", line, pattern, err)
					} else if isMatch {
						fmt.Println(line)
					}
				} else {
					fmt.Println(line)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to scan through the file (%s), error: %s", logFilePath, err)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(stepinfosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stepinfosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	isTimeOnlyModeFlag = stepinfosCmd.Flags().Bool("only-times", false, "If enabled will only print step run times")
}
