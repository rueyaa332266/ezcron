package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/rueyaa332266/ezcron/translator"
	"github.com/spf13/cobra"
)

var scheduleTypeSuggest = []prompt.Suggest{
	{Text: "Time", Description: "Create a schedule at specific time or time interval"},
	{Text: "Daily", Description: "Create a daily schedule at specific timer"},
	{Text: "Weekly", Description: "Create a weekly schedule on specific weekday at specific time"},
	{Text: "Monthly", Description: "Create a monthly schedule on specific monthday at specific time"},
	{Text: "Yearly", Description: "create a yearly schedule in specific month on specific monthday at specific time"},
}

func makeTimeSuggest() []prompt.Suggest {
	var timeSuggest []prompt.Suggest
	for i := 0; i < 23; i++ {
		for j := 0; j < 59; j++ {
			hour := translator.AddZeorforTenDigit(strconv.Itoa(i))
			minute := translator.AddZeorforTenDigit(strconv.Itoa(j))
			suggest := prompt.Suggest{Text: hour + ":" + minute}
			timeSuggest = append(timeSuggest, suggest)
		}

	}
	return timeSuggest
}

func executor(in string) {
	// for the create func
}

func completer(in prompt.Document) []prompt.Suggest {
	args := strings.Split(in.TextBeforeCursor(), " ")
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(scheduleTypeSuggest, args[0], true)
	}
	first := args[0]
	switch first {
	case "Time":
		second := args[1]
		if len(args) == 2 {
			timeAdposition := []prompt.Suggest{{Text: "at"}}
			return prompt.FilterHasPrefix(timeAdposition, second, true)
		}
		third := args[2]
		if len(args) == 3 {
			return prompt.FilterHasPrefix(makeTimeSuggest(), third, true)
		}

	default:
		return prompt.FilterHasPrefix(scheduleTypeSuggest, in.GetWordBeforeCursor(), true)
	}
	return []prompt.Suggest{}
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create cron expression",
	Long:  `Create cron expression with prompt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Select the schedule of type and press entre")
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("Press tab and select >> "),
			prompt.OptionTitle("excron create"),
		)
		p.Run()
	},
}
