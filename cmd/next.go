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

func nextExecTime(cronExpression string) (time.Time, error) {
	valid, _ := translator.MatchCronReg(cronExpression)
	if valid {
		cronExpression = translator.MonthToNum(cronExpression)
		cronExpression = translator.WeekDayToNum(cronExpression)
		sched, err := cron.ParseStandard(cronExpression)
		if err != nil {
			fmt.Println(cronExpression, "Invalid syntax")
			return time.Now(), err
		}
		return sched.Next(time.Now()), nil
	}
	fmt.Println(cronExpression, "Invalid syntax")
	err := errors.New("Invalid syntax")
	return time.Now(), err
}

var cmdNext = &cobra.Command{
	Use:   "next [cron expression]",
	Short: "Return next execute time",
	Long:  `Show the next execute time when inputing cron expression`,
	Run:   getNextTime,
}

func getNextTime(cmd *cobra.Command, args []string) {
	// show help message if got no args
	if len(args) < 1 {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}
	cronExpression := strings.Join(args, " ")
	valid, _ := translator.MatchCronReg(cronExpression)
	if valid {
		next, _ := nextExecTime(cronExpression)
		fmt.Println("Next execute time:", next)
	} else {
		fmt.Println("Invalid syntax")
	}
}
