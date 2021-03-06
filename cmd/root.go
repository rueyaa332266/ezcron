package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/rueyaa332266/ezcron/translator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ezcron",
	Run: translateFromPipe,
}

func translateFromPipe(cmd *cobra.Command, args []string) {
	// input from pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(os.Stdin)
		cronExpression := strings.TrimSuffix(buf.String(), "\n")
		valid, result := translator.MatchCronReg(cronExpression)
		if valid {
			translator.Explain(result)
		} else {
			fmt.Println("invalid syntax")
			os.Exit(1)
		}
	} else {
		cmd.Help()
	}
}

func init() {
	rootCmd.AddCommand(cmdNext)
	rootCmd.AddCommand(cmdTranslate)
	rootCmd.AddCommand(cmdExpression)
}

// Execute for cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
