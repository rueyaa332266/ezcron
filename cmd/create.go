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
	{Text: "Daily_schedule:", Description: "Create a daily schedule at specific time"},
	{Text: "Weekly_schedule:", Description: "Create a weekly schedule on specific weekday at specific time"},
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
	for i := 1; i <= 60; i++ {
		minute := strconv.Itoa(i)
		suggest := prompt.Suggest{Text: minute + "_minute"}
		minuteSuggest = append(minuteSuggest, suggest)
	}
	return minuteSuggest
}

func makeHourSuggest() []prompt.Suggest {
	var hourSuggest []prompt.Suggest
	for i := 1; i <= 24; i++ {
		minute := strconv.Itoa(i)
		suggest := prompt.Suggest{Text: minute + "_hour"}
		hourSuggest = append(hourSuggest, suggest)
	}
	return hourSuggest
}

func makeWeekdaySuggest() []prompt.Suggest {
	var weekDaysuggest []prompt.Suggest
	dayList := [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for _, v := range dayList {
		suggest := prompt.Suggest{Text: v, Description: "at 00:00"}
		weekDaysuggest = append(weekDaysuggest, suggest)
	}
	return weekDaysuggest
}

func executor(in string) {
	if in == "" {
		fmt.Println("Empty input")
		os.Exit(1)
	}

	// split and ignore space
	f := func(c rune) bool {
		return c == ' '
	}
	inputs := strings.FieldsFunc(in, f)
	// fmt.Println(inputs)
	switch inputs[0] {
	case "Time_schedule:":
		last := inputs[len(inputs)-1]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if strings.Contains(last, "minute") {
			fmt.Println("*/" + strings.Split(last, "_")[0] + " * * * *")
		} else if strings.Contains(last, "hour") {
			fmt.Println("* */" + strings.Split(last, "_")[0] + " * * *")
		} else if re.MatchString(last) {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			fmt.Println(minute + " " + hour + " * * *")
		} else {
			fmt.Println("Time schedule is not completed")
		}
	case "Daily_schedule:":
		last := inputs[len(inputs)-1]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if re.MatchString(last) {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			fmt.Println(minute + " " + hour + " */1 * *")
		} else if last == "every_day" {
			fmt.Println("0 0 */1 * *")
		} else {
			fmt.Println("Daily schedule is not completed")
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
			timeAdposition := []prompt.Suggest{{Text: "at", Description: "__:__"}, {Text: "every_minute", Description: "per minute"}, {Text: "every_hour", Description: "per hour"}}
			prompt.OptionPreviewSuggestionTextColor(prompt.Red)
			return prompt.FilterHasPrefix(timeAdposition, second, true)
		}
		third := args[2]
		switch second {
		case "at":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeTimeSuggest(), third, true)
			}
		case "every_minute":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeMinuteSuggest(), third, true)
			}
		case "every_hour":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeHourSuggest(), third, true)
			}
		}
	case "Daily_schedule:":
		second := args[1]
		if len(args) == 2 {
			dayAdposition := []prompt.Suggest{{Text: "every_day", Description: "every day at 00:00"}, {Text: "every_day_at", Description: "every day at __:__"}}
			return prompt.FilterHasPrefix(dayAdposition, second, true)
		}
		third := args[2]
		switch second {
		case "every_day_at":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeTimeSuggest(), third, true)
			}
		}
	case "Weekly_schedule:":
		second := args[1]
		if len(args) == 2 {
			dayAdposition := []prompt.Suggest{{Text: "on_every", Description: "weekday"}}
			return prompt.FilterHasPrefix(dayAdposition, second, true)
		}
		third := args[2]
		switch second {
		case "on_every":
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeWeekdaySuggest(), third, true)
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
