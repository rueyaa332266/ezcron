package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	cron "github.com/robfig/cron/v3"
	"github.com/rueyaa332266/ezcron/translator"
	"github.com/spf13/cobra"
)

func nextExecTime(cronExpression string) error {
	valid, _ := translator.MatchCronReg(cronExpression)
	if valid {
		cronExpression = translator.MonthToNum(cronExpression)
		cronExpression = translator.WeekDayToNum(cronExpression)
		sched, err := cron.ParseStandard(cronExpression)
		if err != nil {
			fmt.Println(cronExpression, "invalid syntax")
			return err
		} else {
			fmt.Println("Next execute time:", sched.Next(time.Now()))
			return nil
		}
	} else {
		fmt.Println(cronExpression, "invalid syntax")
		err := errors.New("invalid syntax")
		return err
	}
}

var cmdNext = &cobra.Command{
	Use:   "next [cron expression]",
	Short: "return next execute time",
	Long:  `Show the next execute time when inputing cron expression`,
	Run: func(cmd *cobra.Command, args []string) {
		// when input from pipe
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("not avalibe from pipe")
			os.Exit(1)
		} else {
			// show help message if got no args
			if len(args) < 1 {
				cmd.Help()
				os.Exit(0)
			}
			cronExpression := strings.Join(args, " ")
			valid, _ := translator.MatchCronReg(cronExpression)
			if valid {
				nextExecTime(cronExpression)
			} else {
				fmt.Println("invalid syntax")
				os.Exit(1)
			}

		}
	},
}
