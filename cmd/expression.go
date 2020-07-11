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
	{Text: "Monthly_schedule:", Description: "Create a monthly schedule on specific monthday at specific time"},
	{Text: "Yearly_schedule:", Description: "Create a yearly schedule in specific month on specific monthday at specific time"},
}

var dayWList = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var monthList = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var dayList = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

func makeTimeSuggest(t string) []prompt.Suggest {
	var timeSuggest []prompt.Suggest
	var suggest prompt.Suggest
	for i := 0; i < 24; i++ {
		if t == "minute" && i > 0 {
			break
		} else if t == "hour" {
			hour := strconv.Itoa(i + 1)
			suggest = prompt.Suggest{Text: hour + "_hour"}
			timeSuggest = append(timeSuggest, suggest)
		}
		for j := 0; j < 60; j++ {
			if t == "hour" {
				break
			} else if t == "minute" {
				minute := strconv.Itoa(j + 1)
				suggest = prompt.Suggest{Text: minute + "_minute"}
			} else if t == "time" {
				hour := translator.AddZeorforTenDigit(strconv.Itoa(i))
				minute := translator.AddZeorforTenDigit(strconv.Itoa(j))
				suggest = prompt.Suggest{Text: hour + ":" + minute}
			}
			timeSuggest = append(timeSuggest, suggest)
		}
	}
	return timeSuggest
}

func makeWeekdaySuggest() []prompt.Suggest {
	var weekDaySuggest []prompt.Suggest
	for _, v := range dayWList {
		suggest := prompt.Suggest{Text: v, Description: "default at 00:00"}
		weekDaySuggest = append(weekDaySuggest, suggest)
	}
	return weekDaySuggest
}

func makeMonthdaySuggest() []prompt.Suggest {
	var monthDaySuggest []prompt.Suggest
	for i := range dayList {
		day := translator.OrdinalDay(dayList[i])
		suggest := prompt.Suggest{Text: day + "_day", Description: "of month"}
		monthDaySuggest = append(monthDaySuggest, suggest)
	}
	return monthDaySuggest
}

func makeMonthdayNumberSuggest(src []string) []prompt.Suggest {
	var monthDayNumbersSuggest []prompt.Suggest
	for _, v := range src {
		suggest := prompt.Suggest{Text: translator.OrdinalDay(v), Description: "default at 00:00"}
		monthDayNumbersSuggest = append(monthDayNumbersSuggest, suggest)
	}
	return monthDayNumbersSuggest
}

func makeMonthNumSuggest() []prompt.Suggest {
	var monthNumSuggest []prompt.Suggest
	for i := 1; i < 13; i++ {
		suggest := prompt.Suggest{Text: strconv.Itoa(i) + "_month", Description: "default at 00:00"}
		monthNumSuggest = append(monthNumSuggest, suggest)
	}
	return monthNumSuggest
}

func makeMonthSuggest() []prompt.Suggest {
	var monthSuggest []prompt.Suggest
	for _, v := range monthList {
		suggest := prompt.Suggest{Text: v}
		monthSuggest = append(monthSuggest, suggest)
	}
	return monthSuggest
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
	switch inputs[0] {
	case "Time_schedule:":
		executeTimeSchedule(inputs)
	case "Daily_schedule:":
		executeDailySchedule(inputs)
	case "Weekly_schedule:":
		executeWeeklySchedule(inputs)
	case "Monthly_schedule:":
		executeMonthlySchedule(inputs)
	case "Yearly_schedule:":
		executeYearlySchedule(inputs)
	default:
		fmt.Println("invalid input")
	}
	os.Exit(0)
}

func executeTimeSchedule(inputs []string) {
	if len(inputs) == 3 {
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
			fmt.Println("Invalid time schedule")
		}
	} else {
		fmt.Println("Invalid time schedule")
	}
}

func executeDailySchedule(inputs []string) {
	if len(inputs) == 2 || len(inputs) == 4 {
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
			fmt.Println("Invalid daily schedule")
		}
	} else {
		fmt.Println("Invalid daily schedule")
	}
}

func executeWeeklySchedule(inputs []string) {
	last := inputs[len(inputs)-1]
	if len(inputs) == 5 {
		weekDay := inputs[len(inputs)-3]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if re.MatchString(last) && contains(dayWList, weekDay) {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			fmt.Println(minute + " " + hour + " * * " + translator.WeekDayToNum(weekDay))
		} else {
			fmt.Println("Invalid weekly schedule")
		}
	} else if len(inputs) == 3 {
		if contains(dayWList, last) {
			fmt.Println("0 0 * * " + translator.WeekDayToNum(last))
		} else {
			fmt.Println("Invalid weekly schedule")
		}
	} else {
		fmt.Println("Invalid weekly schedule")
	}
}

func executeMonthlySchedule(inputs []string) {
	last := inputs[len(inputs)-1]
	if len(inputs) == 7 {
		monthDay := inputs[2]
		perMonth := inputs[4]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if re.MatchString(last) && strings.Contains(monthDay, "_day") && strings.Contains(perMonth, "_month") {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			perMonth := strings.TrimRight(perMonth, "_month")
			fmt.Println(minute + " " + hour + " " + monthDay + " */" + perMonth + " *")
		} else {
			fmt.Println("Invalid monthly schedule")
		}
	} else if len(inputs) == 6 {
		monthDay := inputs[2]
		re := regexp.MustCompile(`\d\d:\d\d`)
		if re.MatchString(last) && strings.Contains(monthDay, "_day") {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			fmt.Println(minute + " " + hour + " " + monthDay + " */1 *")
		} else {
			fmt.Println("Invalid monthly schedule")
		}
	} else if len(inputs) == 5 {
		monthDay := inputs[2]
		perMonth := inputs[4]
		if strings.Contains(monthDay, "_day") && strings.Contains(perMonth, "_month") {
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			perMonth := strings.TrimRight(perMonth, "_month")
			fmt.Println("0 0 " + monthDay + " */" + perMonth + " *")
		} else {
			fmt.Println("Invalid monthly schedule")
		}
	} else if len(inputs) == 4 {
		monthDay := inputs[2]
		if strings.Contains(monthDay, "_day") {
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			fmt.Println("0 0 " + monthDay + " */1 *")
		} else {
			fmt.Println("Invalid monthly schedule")
		}
	} else {
		fmt.Println("Invalid monthly schedule")
	}
}

func executeYearlySchedule(inputs []string) {
	last := inputs[len(inputs)-1]
	if len(inputs) == 6 {
		month := inputs[2]
		monthDay := inputs[3]
		reTime := regexp.MustCompile(`\d\d:\d\d`)
		reDay := regexp.MustCompile(`^\d{1,2}[a-z]{2}$`)
		if reTime.MatchString(last) && contains(monthList, month) && reDay.MatchString(monthDay) {
			time := strings.Split(last, ":")
			minute := strings.TrimPrefix(time[1], "0")
			hour := strings.TrimPrefix(time[0], "0")
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			month := translator.MonthToNum(month)
			fmt.Println(minute + " " + hour + " " + monthDay + " " + month + " *")
		} else {
			fmt.Println("Invalid yearly schedule")
		}
	} else if len(inputs) == 4 {
		month := inputs[2]
		monthDay := inputs[3]
		reDay := regexp.MustCompile(`^\d{1,2}[a-z]{2}$`)
		if contains(monthList, month) && reDay.MatchString(monthDay) {
			re := regexp.MustCompile(`\d{1,2}`)
			monthDay := re.FindAllString(monthDay, 1)[0]
			month := translator.MonthToNum(month)
			fmt.Println("0 0 " + monthDay + " " + month + " *")
		} else {
			fmt.Println("Invalid yearly schedule")
		}
	} else {
		fmt.Println("Invalid yearly schedule")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	args := strings.Split(in.TextBeforeCursor(), " ")
	var suggest []prompt.Suggest
	var sub string
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(scheduleTypeSuggest, args[0], true)
	}
	switch args[0] {
	case "Time_schedule:":
		suggest, sub = completeTimeSchedule(args)
	case "Daily_schedule:":
		suggest, sub = completeDailySchedule(args)
	case "Weekly_schedule:":
		suggest, sub = completeWeeklySchedule(args)
	case "Monthly_schedule:":
		suggest, sub = completeMonthlySchedule(args)
	case "Yearly_schedule:":
		suggest, sub = completeYearlySchedule(args)
	}
	return prompt.FilterHasPrefix(suggest, sub, true)
}

func completeTimeSchedule(args []string) ([]prompt.Suggest, string) {
	var suggest []prompt.Suggest
	sub := args[1]
	if len(args) == 2 {
		suggest = []prompt.Suggest{{Text: "at", Description: "__:__ every day"}, {Text: "every_minute", Description: "per minute"}, {Text: "every_hour", Description: "per hour"}}
		return suggest, sub
	}
	sub = args[2]
	if len(args) == 3 {
		suggest = makeSuggestByPreWord(args[1])
	}
	return suggest, sub
}

func completeDailySchedule(args []string) ([]prompt.Suggest, string) {
	var suggest []prompt.Suggest
	sub := args[1]
	if len(args) == 2 {
		suggest = []prompt.Suggest{{Text: "every_day", Description: "default at 00:00"}}
		return suggest, sub
	}
	sub = args[2]
	if len(args) == 3 {
		suggest = []prompt.Suggest{{Text: "at", Description: "__:__"}}
		return suggest, sub
	}
	sub = args[3]
	if len(args) == 4 {
		suggest = makeSuggestByPreWord(args[2])
	}
	return suggest, sub
}

func completeWeeklySchedule(args []string) ([]prompt.Suggest, string) {
	var suggest []prompt.Suggest
	sub := args[1]
	if len(args) == 2 {
		suggest = []prompt.Suggest{{Text: "on_every", Description: "weekday"}}
		return suggest, sub
	}
	sub = args[2]
	if args[1] == "on_every" {
		if len(args) == 3 {
			suggest = makeWeekdaySuggest()
			return suggest, sub
		}
		sub = args[3]
		if contains(dayWList, args[2]) {
			if len(args) == 4 {
				suggest = []prompt.Suggest{{Text: "at", Description: "__:__"}}
				return suggest, sub
			}
			sub = args[4]
			if len(args) == 5 {
				suggest = makeSuggestByPreWord(args[3])
			}
		}
	}
	return suggest, sub
}

func completeMonthlySchedule(args []string) ([]prompt.Suggest, string) {
	var suggest []prompt.Suggest
	sub := args[1]
	if len(args) == 2 {
		suggest = []prompt.Suggest{{Text: "on", Description: "monthday"}}
		return suggest, sub
	}
	sub = args[2]
	if args[1] == "on" {
		if len(args) == 3 {
			suggest = makeMonthdaySuggest()
			return suggest, sub
		}
		sub = args[3]
		if strings.Contains(args[2], "_day") {
			if len(args) == 4 {
				suggest = []prompt.Suggest{{Text: "of_every_month", Description: "per month, default at 00:00"}, {Text: "of_every", Description: "period of month"}}
				return suggest, sub
			}
			sub = args[4]
			switch args[3] {
			case "of_every_month":
				if len(args) == 5 {
					suggest = []prompt.Suggest{{Text: "at", Description: "__:__"}}
					return suggest, sub
				}
				sub = args[5]
				if len(args) == 6 {
					suggest = makeSuggestByPreWord(args[4])
					return suggest, sub
				}
			case "of_every":
				if len(args) == 5 {
					suggest = makeMonthNumSuggest()
					return suggest, sub
				}
				sub = args[5]
				if strings.Contains(args[4], "_month") {
					if len(args) == 6 {
						suggest = []prompt.Suggest{{Text: "at", Description: "__:__"}}
						return suggest, sub
					}
					sub = args[6]
					if len(args) == 7 {
						suggest = makeSuggestByPreWord(args[5])
					}
				}
			}
		}
	}
	return suggest, sub
}

func completeYearlySchedule(args []string) ([]prompt.Suggest, string) {
	var suggest []prompt.Suggest
	sub := args[1]
	if len(args) == 2 {
		suggest = []prompt.Suggest{{Text: "in_every", Description: "month_day"}}
		return suggest, sub
	}
	sub = args[2]
	if args[1] == "in_every" {
		if len(args) == 3 {
			suggest = makeMonthSuggest()
			return suggest, sub
		}
		sub = args[3]
		if contains(monthList, args[2]) {
			if len(args) == 4 {
				switch args[2] {
				case "February":
					suggest = makeMonthdayNumberSuggest(dayList[:28])
				case "April", "June", "September", "November":
					suggest = makeMonthdayNumberSuggest(dayList[:30])
				default:
					suggest = makeMonthdayNumberSuggest(dayList)
				}
				return suggest, sub
			}
			sub = args[4]
			re := regexp.MustCompile(`^\d{1,2}[a-z]{2}$`)
			if re.MatchString(args[3]) {
				if len(args) == 5 {
					suggest = []prompt.Suggest{{Text: "at", Description: "__:__"}}
					return suggest, sub
				}
				sub = args[5]
				if len(args) == 6 {
					suggest = makeSuggestByPreWord(args[4])
				}
			}
		}
	}
	return suggest, sub
}

func makeSuggestByPreWord(pre string) []prompt.Suggest {
	var suggest []prompt.Suggest
	switch pre {
	case "at":
		suggest = makeTimeSuggest("time")
	case "every_minute":
		suggest = makeTimeSuggest("minute")
	case "every_hour":
		suggest = makeTimeSuggest("hour")
	}
	return suggest
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

var cmdExpression = &cobra.Command{
	Use:   "expression",
	Short: "Create a cron expression",
	Long:  `Create a cron expression with prompt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Follow the prompts and create the cron expression")
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("Press tab for prompts >> "),
			prompt.OptionTitle("Create cron expression"),
		)
		p.Run()
	},
}
