package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/rueyaa332266/ezcron/translator"
	"github.com/spf13/cobra"
)

var scheduleTypeSuggest = []prompt.Suggest{
	{Text: "Time_schedule:", Description: "Create a schedule at specific time or time interval"},
	{Text: "Daily", Description: "Create a daily schedule at specific timer"},
	{Text: "Weekly", Description: "Create a weekly schedule on specific weekday at specific time"},
	{Text: "Monthly", Description: "Create a monthly schedule on specific monthday at specific time"},
	{Text: "Yearly", Description: "create a yearly schedule in specific month on specific monthday at specific time"},
}

func makeTimeSuggest() []prompt.Suggest {
	var timeSuggest []prompt.Suggest
	for i := 0; i < 24; i++ {
		for j := 0; j < 60; j++ {
			hour := translator.AddZeorforTenDigit(strconv.Itoa(i))
			minute := translator.AddZeorforTenDigit(strconv.Itoa(j))
			suggest := prompt.Suggest{Text: hour + ":" + minute}
			timeSuggest = append(timeSuggest, suggest)
		}

	}
	return timeSuggest
}

func makeMinuteSuggest() []prompt.Suggest {
	var minuteSuggest []prompt.Suggest
	for i := 0; i < 60; i++ {
		minute := strconv.Itoa(i)
		suggest := prompt.Suggest{Text: minute + "_minute"}
		minuteSuggest = append(minuteSuggest, suggest)
	}
	return minuteSuggest
}

func executor(in string) {
	inputs := strings.Split(in, " ")
	// fmt.Println(inputs)
	switch inputs[0] {
	case "Time_schedule:":
		last := inputs[len(inputs)-1]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if strings.Contains(last, "minute") {
			fmt.Println("*/" + strings.Split(last, "_")[0] + " * * * *")
		} else if re.MatchString(last) {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			fmt.Println(minute + " " + hour + " * * *")
		} else {
			fmt.Println("Time schedule is not completed")
		}
	default:
		fmt.Println("not implement")
	}

	os.Exit(0)
}

func completer(in prompt.Document) []prompt.Suggest {
	args := strings.Split(in.TextBeforeCursor(), " ")
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(scheduleTypeSuggest, args[0], true)
	}
	first := args[0]
	switch first {
	case "Time_schedule:":
		second := args[1]
		if len(args) == 2 {
			timeAdposition := []prompt.Suggest{{Text: "at", Description: "__:__"}, {Text: "every", Description: "per minute"}}
			return prompt.FilterHasPrefix(timeAdposition, second, true)
		}
		third := args[2]
		switch second {
		case "at":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeTimeSuggest(), third, true)
			}
		case "every":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeMinuteSuggest(), third, true)
			}
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
		fmt.Println("Select the schedule of type and press space")
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("Press tab and select >> "),
			prompt.OptionTitle("excron create"),
		)
		p.Run()
	},
}
