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
	{Text: "Yearly_schedule:", Description: "create a yearly schedule in specific month on specific monthday at specific time"},
}

var dayWList = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var monthList = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var dayMList []string

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
	for _, v := range dayWList {
		suggest := prompt.Suggest{Text: v, Description: "default at 00:00"}
		weekDaysuggest = append(weekDaysuggest, suggest)
	}
	return weekDaysuggest
}

func makeMonthdaySuggest() []prompt.Suggest {
	var monthDaysuggest []prompt.Suggest
	for i := 1; i < 32; i++ {
		day := translator.OrdinalFromStr(strconv.Itoa(i))
		suggest := prompt.Suggest{Text: day + "_day", Description: "of month"}
		monthDaysuggest = append(monthDaysuggest, suggest)
	}
	return monthDaysuggest
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
	// fmt.Println(inputs)
	switch inputs[0] {
	case "Time_schedule:":
		if len(inputs) != 3 {
			fmt.Println("input not valid")
			os.Exit(1)
		}
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
			fmt.Println("Time schedule is not valid")
		}
	case "Daily_schedule:":
		if len(inputs) != 2 && len(inputs) != 3 {
			fmt.Println("input not valid")
			os.Exit(1)
		}
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
			fmt.Println("Daily schedule is not valid")
		}
	case "Weekly_schedule:":
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
				fmt.Println("Weekly schedule is not valid")
			}
		} else if len(inputs) == 3 {
			if contains(dayWList, last) {
				fmt.Println("0 0 * * " + translator.WeekDayToNum(last))
			} else {
				fmt.Println("Weekly schedule is not valid")
			}
		} else {
			fmt.Println("input not valid")
			os.Exit(1)
		}
	case "Monthly_schedule:":
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
				fmt.Println("Monthly schedule is not valid")
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
				fmt.Println("Monthly schedule is not valid")
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
				fmt.Println("Monthly schedule is not valid")
			}
		} else if len(inputs) == 4 {
			monthDay := inputs[2]
			if strings.Contains(monthDay, "_day") {
				re := regexp.MustCompile(`\d{1,2}`)
				monthDay := re.FindAllString(monthDay, 1)[0]
				fmt.Println("0 0 " + monthDay + " */1 *")
			} else {
				fmt.Println("Monthly schedule is not valid")
			}
		} else {
			fmt.Println("input not valid")
			os.Exit(1)
		}
	case "Yearly_schedule:":
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
				fmt.Println("Weekly schedule is not valid")
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
			}
		} else {
			fmt.Println("input not valid")
			os.Exit(1)
		}
	default:
		fmt.Println("input not valid")
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
		if second == "every_day_at" && len(args) == 3 {
			return prompt.FilterHasPrefix(makeTimeSuggest(), third, true)
		}
	case "Weekly_schedule:":
		second := args[1]
		if len(args) == 2 {
			dayAdposition := []prompt.Suggest{{Text: "on_every", Description: "weekday"}}
			return prompt.FilterHasPrefix(dayAdposition, second, true)
		}
		third := args[2]
		if second == "on_every" {
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeWeekdaySuggest(), third, true)
			}
			fourth := args[3]
			if contains(dayWList, third) {
				if len(args) == 4 {
					return prompt.FilterHasPrefix([]prompt.Suggest{{Text: "at", Description: "__:__"}}, fourth, true)
				}
				fifth := args[4]
				if fourth == "at" && len(args) == 5 {
					return prompt.FilterHasPrefix(makeTimeSuggest(), fifth, true)
				}
			}
		}
	case "Monthly_schedule:":
		second := args[1]
		if len(args) == 2 {
			dayAdposition := []prompt.Suggest{{Text: "on", Description: "monthday"}}
			return prompt.FilterHasPrefix(dayAdposition, second, true)
		}
		third := args[2]
		if second == "on" {
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeMonthdaySuggest(), third, true)
			}
			fourth := args[3]
			if strings.Contains(third, "_day") {
				if len(args) == 4 {
					return prompt.FilterHasPrefix([]prompt.Suggest{{Text: "of_every_month", Description: "per month, default at 00:00"}, {Text: "of_every", Description: "period of month"}}, fourth, true)
				}
				fifth := args[4]
				switch fourth {
				case "of_every_month":
					if len(args) == 5 {
						return prompt.FilterHasPrefix([]prompt.Suggest{{Text: "at", Description: "__:__"}}, fifth, true)
					}
					sixth := args[5]
					if fifth == "at" && len(args) == 6 {
						return prompt.FilterHasPrefix(makeTimeSuggest(), sixth, true)
					}
				case "of_every":
					if len(args) == 5 {
						return prompt.FilterHasPrefix(makeMonthNumSuggest(), fifth, true)
					}
					sixth := args[5]
					if strings.Contains(fifth, "_month") {
						if len(args) == 6 {
							return prompt.FilterHasPrefix([]prompt.Suggest{{Text: "at", Description: "__:__"}}, sixth, true)
						}
						seventh := args[6]
						if sixth == "at" && len(args) == 7 {
							return prompt.FilterHasPrefix(makeTimeSuggest(), seventh, true)
						}
					}
				}
			}
		}
	case "Yearly_schedule:":
		second := args[1]
		if len(args) == 2 {
			dayAdposition := []prompt.Suggest{{Text: "in_every", Description: "month_day"}}
			return prompt.FilterHasPrefix(dayAdposition, second, true)
		}
		third := args[2]
		if second == "in_every" {
			if len(args) == 3 {
				return prompt.FilterHasPrefix(makeMonthSuggest(), third, true)
			}
			fourth := args[3]
			if contains(monthList, third) {
				if len(args) == 4 {
					// make date 28 30 31
					var day []string
					for i := 1; i < 32; i++ {
						day = append(day, strconv.Itoa(i))
					}
					day28 := day[:28]
					day30 := day[:30]
					f := func(src []string) []prompt.Suggest {
						var suggests []prompt.Suggest
						for _, v := range src {
							suggest := prompt.Suggest{Text: translator.OrdinalFromStr(v), Description: "default at 00:00"}
							suggests = append(suggests, suggest)
						}
						return suggests
					}
					switch third {
					case "February":
						return prompt.FilterHasPrefix(f(day28), fourth, true)
					case "April", "June", "September", "November":
						return prompt.FilterHasPrefix(f(day30), fourth, true)
					default:
						return prompt.FilterHasPrefix(f(day), fourth, true)
					}

				}
				fifth := args[4]
				re := regexp.MustCompile(`^\d{1,2}[a-z]{2}$`)
				if re.MatchString(fourth) {
					if len(args) == 5 {
						return prompt.FilterHasPrefix([]prompt.Suggest{{Text: "at", Description: "__:__"}}, fifth, true)
					}
					sixth := args[5]
					if fifth == "at" && len(args) == 6 {
						return prompt.FilterHasPrefix(makeTimeSuggest(), sixth, true)
					}
				}
			}
		}
	default:
		return prompt.FilterHasPrefix(scheduleTypeSuggest, in.GetWordBeforeCursor(), true)
	}
	return []prompt.Suggest{}
}
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
func isDayM(s string) bool {
	for i := 1; i < 32; i++ {
		day := translator.OrdinalFromStr(strconv.Itoa(i))
		dayMList = append(dayMList, day)
	}
	for _, a := range dayMList {
		if a == s {
			return true
		}
	}
	return false
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
