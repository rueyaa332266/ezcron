package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rueyaa332266/ezcron/translator"
	"github.com/spf13/cobra"
)

var cmdTranslate = &cobra.Command{
	Use:   "translate [cron expression]",
	Short: "Translate into human-friendly language",
	Long:  `Translate cron expression into human-friendly language`,
	Run:   translate,
}

func translate(cmd *cobra.Command, args []string) {
	// show help message if got no args
	if len(args) < 1 {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}
	cronExpression := strings.Join(args, " ")
	valid, checkResult := translator.MatchCronReg(cronExpression)
	if valid {
		translator.Explain(checkResult)
	} else {
		fmt.Println("invalid syntax")
		os.Exit(1)
	}
}
