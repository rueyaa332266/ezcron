package translator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

type CheckResult struct {
	Valid   bool
	Input   string
	Pattern string
}

// 	Field name     Mandatory?   Allowed values    Allowed special characters
// ----------     ----------   --------------    --------------------------
// Minutes        Yes          0-59              * / , -
// Hours          Yes          0-23              * / , -
// Day of month   Yes          1-31              * / , -
// Month          Yes          1-12 or JAN-DEC   * / , -
// Day of week    Yes          0-6 or SUN-SAT    * / , -
func MatchCronReg(cronExpression string) (bool, map[string]*CheckResult) {
	var checkResults = make(map[string]*CheckResult)
	fieldSlice := strings.Split(cronExpression, " ")

	for index, field := range fieldSlice {
		switch index {
		case 0: // minute
			checkResults["minute"] = CheckMinuteReg(field)
		case 1: // hour
			checkResults["hour"] = CheckHourReg(field)
		case 2: // day of the month
			checkResults["dayM"] = CheckDayMonthReg(field)
		case 3: // month
			checkResults["month"] = CheckMonthReg(field)
		case 4: // day of the week
			checkResults["dayW"] = CheckDayWeekReg(field)
		}
	}
	valid := checkResults["minute"].Valid && checkResults["hour"].Valid && checkResults["dayM"].Valid && checkResults["month"].Valid && checkResults["dayW"].Valid
	return valid, checkResults
}

func checkFieldPattern(str string, patternList map[string]string) (bool, string) {
	var pattern string
	valid := true
	asteriskReg := regexp.MustCompile(`^` + patternList["asterisk"] + `$`)
	numberReg := regexp.MustCompile(`^` + patternList["number"] + `$`)
	commaReg := regexp.MustCompile(`^` + patternList["comma"] + `$`)
	hyphenReg := regexp.MustCompile(`^` + patternList["hyphen"] + `$`)
	slashReg := regexp.MustCompile(`^` + patternList["slash"] + `$`)

	switch {
	case asteriskReg.MatchString(str) == true:
		pattern = "asterisk"
	case numberReg.MatchString(str) == true:
		pattern = "number"
	case commaReg.MatchString(str) == true:
		pattern = "comma"
	case hyphenReg.MatchString(str) == true:
		pattern = "hyphen"
	case slashReg.MatchString(str) == true:
		pattern = "slash"
	default:
		pattern = "not match"
		valid = false
	}

	return valid, pattern
}

// check the logic and format in some pattern
func numLogicFormat(input string, pattern string) (bool, string) {
	valid := true
	switch pattern {
	case "hyphen":
		slice := strings.Split(input, "-")
		valid = slice[0] < slice[1]
	case "comma":
		unigueStr := uniqueSlice(strings.Split(input, ","))
		input = strings.Join(unigueStr, ",")
	case "slash":
		slice := strings.Split(input, "/")
		validLeft, strLeft := numLogicFormat(slice[0], CheckMinuteReg(slice[0]).Pattern)
		validRight, strRight := numLogicFormat(slice[1], CheckMinuteReg(slice[1]).Pattern)
		valid = validLeft && validRight
		input = strLeft + "/" + strRight
	}
	return valid, input
}

func CheckMinuteReg(str string) *CheckResult {
	asterisk := `\*`
	number := `([0-9]|[1-5][0-9])`
	comma := `(` + number + `\,){1,59}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}
	valid, pattern := checkFieldPattern(str, patternList)
	if valid {
		valid, str = numLogicFormat(str, pattern)
	}
	return &CheckResult{Valid: valid, Input: str, Pattern: pattern}
}

func CheckHourReg(str string) *CheckResult {
	asterisk := `\*`
	number := `([0-9]|1[0-9]|2[0-3])`
	comma := `(` + number + `\,){1,23}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}
	valid, pattern := checkFieldPattern(str, patternList)
	if valid {
		valid, str = numLogicFormat(str, pattern)
	}
	return &CheckResult{Valid: valid, Input: str, Pattern: pattern}
}

func CheckDayMonthReg(str string) *CheckResult {
	asterisk := `\*`
	number := `([1-9]|[1-2][0-9]|3[0-1])`
	comma := `(` + number + `\,){1,30}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}
	valid, pattern := checkFieldPattern(str, patternList)
	if valid {
		valid, str = numLogicFormat(str, pattern)
	}
	return &CheckResult{Valid: valid, Input: str, Pattern: pattern}
}

func CheckMonthReg(str string) *CheckResult {
	str = MonthToNum(str)
	asterisk := `\*`
	number := `([1-9]|[1][0-2])`
	comma := `(` + number + `\,){1,11}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}
	valid, pattern := checkFieldPattern(str, patternList)
	if valid {
		valid, str = numLogicFormat(str, pattern)
	}
	return &CheckResult{Valid: valid, Input: str, Pattern: pattern}
}

func CheckDayWeekReg(str string) *CheckResult {
	str = WeekDayToNum(str)
	asterisk := `\*`
	number := `([0-6])`
	comma := `(` + number + `\,){1,5}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}
	valid, pattern := checkFieldPattern(str, patternList)
	if valid {
		valid, str = numLogicFormat(str, pattern)
	}
	return &CheckResult{Valid: valid, Input: str, Pattern: pattern}
}

func Explain(cronCheckResults map[string]*CheckResult) {
	var explanation string
	minuteCheckResult := cronCheckResults["minute"]
	hourCheckResult := cronCheckResults["hour"]
	dayMonthCheckResult := cronCheckResults["dayM"]
	monthCheckResult := cronCheckResults["month"]
	dayWeekCheckResult := cronCheckResults["dayW"]
	// explain in HH:MM
	if minuteCheckResult.Pattern == hourCheckResult.Pattern && minuteCheckResult.Pattern == "number" {
		explanation += "at every " + AddZeorforTenDigit(hourCheckResult.Input) + ":" + AddZeorforTenDigit(minuteCheckResult.Input)
	} else {
		explanation += explainMinute(minuteCheckResult)
		explanationHour := explainHour(hourCheckResult)
		if explanationHour != "" {
			explanation += ", " + explanationHour
		}
	}
	// explain when one of day month and day week use "*/num"
	if dayMonthCheckResult.Pattern == "slash" && dayWeekCheckResult.Pattern != "asterisk" {
		sliceDayM := strings.Split(dayMonthCheckResult.Input, "/")
		slachLeftDayM := CheckDayMonthReg(sliceDayM[0])
		explanationDayMonth := explainDayMonth(dayMonthCheckResult)
		explanationDayWeek := explainDayWeek(dayWeekCheckResult)
		if slachLeftDayM.Pattern == "asterisk" {
			explanation += ", " + explanationDayMonth + " if it's " + explanationDayWeek
		} else {
			explanation += ", " + explanationDayMonth + " and " + explanationDayWeek
		}
	} else if dayWeekCheckResult.Pattern == "slash" && dayMonthCheckResult.Pattern != "asterisk" {
		sliceDayW := strings.Split(dayWeekCheckResult.Input, "/")
		slachLeftDayW := CheckDayMonthReg(sliceDayW[0])
		explanationDayMonth := explainDayMonth(dayMonthCheckResult)
		explanationDayWeek := explainDayWeek(dayWeekCheckResult)
		if slachLeftDayW.Pattern == "asterisk" {
			explanation += ", " + explanationDayMonth + " if it's " + explanationDayWeek
		} else {
			explanation += ", " + explanationDayMonth + " and " + explanationDayWeek
		}
	} else {
		explanationDayMonth := explainDayMonth(dayMonthCheckResult)
		explanationDayWeek := explainDayWeek(dayWeekCheckResult)
		if explanationDayMonth != "" && explanationDayWeek != "" {
			explanation += ", " + explanationDayMonth + " and " + explanationDayWeek
		} else if explanationDayMonth != "" {
			explanation += ", " + explanationDayMonth
		} else if explanationDayWeek != "" {
			explanation += ", " + explanationDayWeek
		}
	}
	explanationMonth := explainMonth(monthCheckResult)
	if explanationMonth != "" {
		explanation += ", " + explanationMonth
	}
	fmt.Println(explanation)
}

func explainMinute(c *CheckResult) string {
	explanation := "At"
	switch c.Pattern {
	case "asterisk":
		explanation += " every minute"
	case "number":
		minute, _ := strconv.Atoi(c.Input)
		explanation += " every " + humanize.Ordinal(minute) + " minute"
	case "comma":
		minute := c.Input
		slice := strings.Split(minute, ",")
		for i, v := range slice {
			slice[i] = OrdinalFromStr(v)
		}
		explanation += " every " + strings.Join(slice, " and ") + " minute"
	case "hyphen":
		minute := c.Input
		explanation += " every " + strings.Replace(minute, "-", " through ", 1) + " minute"
	case "slash":
		// "asterisk","number","comma","hyphen" / "number","comma"
		minute := c.Input
		var output string
		slice := strings.Split(minute, "/")
		for i, v := range slice {
			CheckResult := CheckMinuteReg(v)
			if i == 0 {
				switch CheckResult.Pattern {
				case "asterisk":
					output += ""
				case "number":
					if CheckResult.Input == "59" {
						output += " from minute 59"
					} else {
						output += " from minute " + CheckResult.Input + " through minute 59"
					}
				case "comma":
					minute := CheckResult.Input
					slice := strings.Split(minute, ",")
					for i, v := range slice {
						if v == "59" {
							slice[i] = "minute 59"
						} else {
							slice[i] = "minute " + v + " through minute 59"
						}
					}
					output += " from " + strings.Join(slice, " and ")
				case "hyphen":
					minute := CheckResult.Input
					output += " from minute " + strings.Replace(minute, "-", " through minute ", 1)
				}
			} else {
				// for pattern "number" and "comma"
				output = CheckResult.Input + " minute" + output
			}
		}
		explanation += " every " + output
	}
	return explanation
}

func explainHour(c *CheckResult) string {
	explanation := "past"
	switch c.Pattern {
	case "asterisk":
		explanation = ""
	case "number":
		hour := AddZeorforTenDigit(c.Input)
		explanation += " " + hour + ":00"
	case "comma":
		hour := c.Input
		slice := strings.Split(hour, ",")
		for i, v := range slice {
			slice[i] = AddZeorforTenDigit(v) + ":00"
		}
		explanation += " " + strings.Join(slice, " and ")
	case "hyphen":
		hour := c.Input
		slice := strings.Split(hour, "-")
		for i, v := range slice {
			slice[i] = AddZeorforTenDigit(v) + ":00"
		}
		explanation += " from " + strings.Join(slice, "-")
	case "slash":
		// "asterisk","number","comma","hyphen" / "number","comma"
		hour := c.Input
		var output string
		slice := strings.Split(hour, "/")
		for i, v := range slice {
			CheckResult := CheckHourReg(v)
			if i == 0 {
				switch CheckResult.Pattern {
				case "asterisk":
					output += ""
				case "number":
					output += " from " + CheckResult.Input + ":00-24:00"
				case "comma":
					hour := CheckResult.Input
					slice := strings.Split(hour, ",")
					for i, v := range slice {
						slice[i] = " " + AddZeorforTenDigit(v) + ":00-24:00"
					}
					output += " from" + strings.Join(slice, " and")
				case "hyphen":
					hour := CheckResult.Input
					slice := strings.Split(hour, "-")
					for i, v := range slice {
						slice[i] = AddZeorforTenDigit(v) + ":00"
					}
					output += " from " + strings.Join(slice, "-")
				}
			} else {
				// for pattern "number" and "comma"
				switch CheckResult.Pattern {
				case "number":
					output = "every " + CheckResult.Input + " hour" + output
				case "comma":
					hour := CheckResult.Input
					slice := strings.Split(hour, ",")
					for i, v := range slice {
						slice[i] = "every " + v + " hour"
					}
					output = strings.Join(slice, " and ") + output
				}
			}
		}
		explanation += " " + output
	}
	return explanation
}

func explainDayMonth(c *CheckResult) string {
	explanation := "on"
	switch c.Pattern {
	case "asterisk":
		explanation = ""
	case "number":
		dayM := c.Input
		explanation += " day " + dayM + " of the month"
	case "comma":
		dayM := c.Input
		explanation += " day " + strings.Replace(dayM, ",", " and ", -1) + " of the month"
	case "hyphen":
		dayM := c.Input
		explanation = "between day " + strings.Replace(dayM, "-", " and ", -1) + " of the month"
	case "slash":
		// "asterisk","number","comma","hyphen" / "number","comma"
		dayM := c.Input
		var output string
		slice := strings.Split(dayM, "/")
		for i, v := range slice {
			CheckResult := CheckDayMonthReg(v)
			if i == 0 {
				switch CheckResult.Pattern {
				case "asterisk":
					output += ""
				case "number":
					output += " from day " + CheckResult.Input + " of the month"
				case "comma":
					dayM := CheckResult.Input
					slice := strings.Split(dayM, ",")
					for i, v := range slice {
						slice[i] = "day " + v + " of the month"
					}
					output += " from " + strings.Join(slice, " and ")
				case "hyphen":
					dayM := CheckResult.Input
					slice := strings.Split(dayM, "-")
					for i, v := range slice {
						slice[i] = "day " + v
					}
					output += " between " + strings.Join(slice, " and ") + " of the month"
				}
			} else {
				// for pattern "number" and "comma"
				output = "every " + CheckResult.Input + " day of month" + output
			}
		}
		explanation += " " + output
	}
	return explanation
}

func explainMonth(c *CheckResult) string {
	explanation := "in"
	monthList := [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	switch c.Pattern {
	case "asterisk":
		explanation = ""
	case "number":
		monthNum, _ := strconv.Atoi(c.Input)
		explanation += " " + monthList[monthNum-1]
	case "comma":
		slice := strings.Split(c.Input, ",")
		for i, v := range slice {
			monthNum, _ := strconv.Atoi(v)
			slice[i] = monthList[monthNum-1]
		}
		explanation += " " + strings.Join(slice, " and ")
	case "hyphen":
		slice := strings.Split(c.Input, "-")
		for i, v := range slice {
			monthNum, _ := strconv.Atoi(v)
			slice[i] = monthList[monthNum-1]
		}
		explanation += " every month from " + strings.Join(slice, " through ")
	case "slash":
		// "asterisk","number","comma","hyphen" / "number","comma"
		month := c.Input
		var output string
		slice := strings.Split(month, "/")
		for i, v := range slice {
			CheckResult := CheckMonthReg(v)
			if i == 0 {
				switch CheckResult.Pattern {
				case "asterisk":
					output += ""
				case "number":
					monthNum, _ := strconv.Atoi(CheckResult.Input)
					output += " from " + monthList[monthNum-1]
				case "comma":
					slice := strings.Split(CheckResult.Input, ",")
					for i, v := range slice {
						monthNum, _ := strconv.Atoi(v)
						slice[i] = monthList[monthNum-1]
					}
					output += " from " + strings.Join(slice, " and ")
				case "hyphen":
					slice := strings.Split(CheckResult.Input, "-")
					for i, v := range slice {
						monthNum, _ := strconv.Atoi(v)
						slice[i] = monthList[monthNum-1]
					}
					output += " from " + strings.Join(slice, " through ")
				}
			} else {
				// for pattern "number" and "comma"
				output = "every " + CheckResult.Input + " month" + output
			}
		}
		explanation += " " + output
	}
	return explanation
}

func explainDayWeek(c *CheckResult) string {
	explanation := "on"
	dayList := [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	switch c.Pattern {
	case "asterisk":
		explanation = ""
	case "number":
		dayNum, _ := strconv.Atoi(c.Input)
		explanation += " " + dayList[dayNum]
	case "comma":
		slice := strings.Split(c.Input, ",")
		for i, v := range slice {
			dayNum, _ := strconv.Atoi(v)
			slice[i] = dayList[dayNum]
		}
		explanation += " " + strings.Join(slice, " and ")
	case "hyphen":
		slice := strings.Split(c.Input, "-")
		for i, v := range slice {
			dayNum, _ := strconv.Atoi(v)
			slice[i] = dayList[dayNum]
		}
		explanation += " every day from " + strings.Join(slice, " through ")
	case "slash":
		// "asterisk","number","comma","hyphen" / "number","comma"
		dayW := c.Input
		var output string
		slice := strings.Split(dayW, "/")
		for i, v := range slice {
			CheckResult := CheckDayWeekReg(v)
			if i == 0 {
				switch CheckResult.Pattern {
				case "asterisk":
					output += ""
				case "number":
					dayNum, _ := strconv.Atoi(CheckResult.Input)
					output += " from " + dayList[dayNum]
				case "comma":
					dayW := CheckResult.Input
					slice := strings.Split(dayW, ",")
					for i, v := range slice {
						dayNum, _ := strconv.Atoi(v)
						slice[i] = dayList[dayNum]
					}
					output += " from " + strings.Join(slice, " and ")
				case "hyphen":
					dayW := CheckResult.Input
					slice := strings.Split(dayW, "-")
					for i, v := range slice {
						dayNum, _ := strconv.Atoi(v)
						slice[i] = dayList[dayNum]
					}
					output += " between " + strings.Join(slice, " and ")
				}
			} else {
				// for pattern "number" and "comma"
				output = "every " + CheckResult.Input + " day of week" + output
			}
		}
		explanation += " " + output
	}
	return explanation
}

func MonthToNum(str string) string {
	pattern := [12]string{`(?i)Jan(uary)?`, `(?i)Feb(ruary)?`, `(?i)Mar(ch)?`, `(?i)Apr(il)?`, `(?i)May`, `(?i)June?`, `(?i)July?`, `(?i)Aug(ust)?`, `(?i)Sep(tember)?`, `(?i)Oct(ober)?`, `(?i)Nov(ember)?`, `(?i)Dec(ember)?`}
	for i, v := range pattern {
		re := regexp.MustCompile(v)
		str = re.ReplaceAllString(str, strconv.Itoa(i+1))
	}
	return str
}

func WeekDayToNum(str string) string {
	pattern := [7]string{`(?i)Sun(day)?`, `(?i)Mon(day)?`, `(?i)Tue(sday)?`, `(?i)Wed(nesday)?`, `(?i)Thu(rsday)?`, `(?i)Fri(day)?`, `(?i)Sat(urday)?`}
	for i, v := range pattern {
		re := regexp.MustCompile(v)
		str = re.ReplaceAllString(str, strconv.Itoa(i))
	}
	return str
}

func OrdinalFromStr(str string) string {
	i, _ := strconv.Atoi(str)
	return humanize.Ordinal(i)
}

func AddZeorforTenDigit(str string) string {
	Re := regexp.MustCompile(`^\d{1}$`)
	if Re.MatchString(str) == true {
		return "0" + str
	}
	return str
}

func uniqueSlice(slice []string) (unique []string) {
	m := map[string]bool{}
	for _, v := range slice {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}
	return unique
}
