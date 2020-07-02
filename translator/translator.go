package translator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/rueyaa332266/multiregexp"
)

type CheckResult struct {
	Valid   bool
	Input   string
	Pattern string
}

var dayList = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var monthList = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

// Field name     Mandatory?   Allowed values    Allowed special characters
// ----------     ----------   --------------    --------------------------
// Minutes        Yes          0-59              * / , -
// Hours          Yes          0-23              * / , -
// Day of month   Yes          1-31              * / , -
// Month          Yes          1-12 or JAN-DEC   * / , -
// Day of week    Yes          0-6 or SUN-SAT    * / , -
func MatchCronReg(cronExpression string) (bool, map[string]*CheckResult) {
	var checkResults = make(map[string]*CheckResult)
	fieldSlice := strings.Split(cronExpression, " ")

	if len(fieldSlice) != 5 {
		return false, checkResults
	}

	for index, field := range fieldSlice {
		switch index {
		case 0: // minute
			checkResults["minute"] = CheckFieldeReg(field, "minute")
		case 1: // hour
			checkResults["hour"] = CheckFieldeReg(field, "hour")
		case 2: // day of the month
			checkResults["dayM"] = CheckFieldeReg(field, "dayM")
		case 3: // month
			checkResults["month"] = CheckFieldeReg(field, "month")
		case 4: // day of the week
			checkResults["dayW"] = CheckFieldeReg(field, "dayW")
		}
	}
	valid := checkResults["minute"].Valid && checkResults["hour"].Valid && checkResults["dayM"].Valid && checkResults["month"].Valid && checkResults["dayW"].Valid
	return valid, checkResults
}

func checkFieldPattern(str string, patternList map[string]string) (bool, string) {
	var pattern string
	var valid bool
	var regs multiregexp.Regexps

	asteriskReg := regexp.MustCompile(`^` + patternList["asterisk"] + `$`)
	numberReg := regexp.MustCompile(`^` + patternList["number"] + `$`)
	commaReg := regexp.MustCompile(`^` + patternList["comma"] + `$`)
	hyphenReg := regexp.MustCompile(`^` + patternList["hyphen"] + `$`)
	slashReg := regexp.MustCompile(`^` + patternList["slash"] + `$`)

	regs = multiregexp.Append(regs, asteriskReg, numberReg, commaReg, hyphenReg, slashReg)
	match := regs.MatchStringWhich(str)

	if len(match) == 0 {
		valid = false
		pattern = "not match"
	} else if len(match) == 1 {
		valid = true
		switch match[0] {
		case 0:
			pattern = "asterisk"
		case 1:
			pattern = "number"
		case 2:
			pattern = "comma"
		case 3:
			pattern = "hyphen"
		case 4:
			pattern = "slash"
		}
	}

	return valid, pattern
}

// check the logic and format in some pattern
func numLogicFormat(input string, pattern string) (bool, string) {
	valid := true
	output := input
	switch pattern {
	case "hyphen":
		slice := strings.Split(input, "-")
		valid = slice[0] < slice[1]
	case "comma":
		unique := uniqueSlice(strings.Split(input, ","))
		output = strings.Join(unique, ",")
	case "slash":
		slice := strings.Split(input, "/")
		validLeft, strLeft := numLogicFormat(slice[0], CheckFieldeReg(slice[0], "minute").Pattern)
		validRight, strRight := numLogicFormat(slice[1], CheckFieldeReg(slice[1], "minute").Pattern)
		valid = validLeft && validRight
		output = strLeft + "/" + strRight
	}
	return valid, output
}

func CheckFieldeReg(str string, feild string) *CheckResult {
	var number, comma string
	switch feild {
	case "minute":
		number = `([0-9]|[1-5][0-9])`
		comma = `(` + number + `\,){1,59}` + number
	case "hour":
		number = `([0-9]|1[0-9]|2[0-3])`
		comma = `(` + number + `\,){1,23}` + number
	case "dayM":
		number = `([1-9]|[1-2][0-9]|3[0-1])`
		comma = `(` + number + `\,){1,30}` + number
	case "month":
		number = `([1-9]|[1][0-2])`
		comma = `(` + number + `\,){1,11}` + number
	case "dayW":
		number = `([0-6])`
		comma = `(` + number + `\,){1,5}` + number
	}
	asterisk := `\*`
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
		explanation += "At every " + AddZeorforTenDigit(hourCheckResult.Input) + ":" + AddZeorforTenDigit(minuteCheckResult.Input)
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
		slachLeftDayM := CheckFieldeReg(sliceDayM[0], "dayM")
		explanationDayMonth := explainDayMonth(dayMonthCheckResult)
		explanationDayWeek := explainDayWeek(dayWeekCheckResult)
		if slachLeftDayM.Pattern == "asterisk" {
			explanation += ", " + explanationDayMonth + " if it's " + explanationDayWeek
		} else {
			explanation += ", " + explanationDayMonth + " and " + explanationDayWeek
		}
	} else if dayWeekCheckResult.Pattern == "slash" && dayMonthCheckResult.Pattern != "asterisk" {
		sliceDayW := strings.Split(dayWeekCheckResult.Input, "/")
		slachLeftDayW := CheckFieldeReg(sliceDayW[0], "dayM")
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
		explanation += " every " + explaiMinuteSlach(c)
	}
	return explanation
}

func explaiMinuteSlach(c *CheckResult) string {
	minute := c.Input
	var output string
	slice := strings.Split(minute, "/")
	for i, v := range slice {
		CheckResult := CheckFieldeReg(v, "minute")
		// check the left side of slash
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
			// check the right side of slash
			// for pattern "number" and "comma"
			output = CheckResult.Input + " minute" + output
		}
	}
	return output
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
		explanation += " " + explainHourSlach(c)
	}
	return explanation
}

func explainHourSlach(c *CheckResult) string {
	hour := c.Input
	var output string
	slice := strings.Split(hour, "/")
	for i, v := range slice {
		CheckResult := CheckFieldeReg(v, "hour")
		// check the left side of slash
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
			// check the right side of slash
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
	return output
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
		explanation += " " + explainDayMonthSlach(c)
	}
	return explanation
}

func explainDayMonthSlach(c *CheckResult) string {
	dayM := c.Input
	var output string
	slice := strings.Split(dayM, "/")
	for i, v := range slice {
		CheckResult := CheckFieldeReg(v, "dayM")
		// check the left side of slash
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
			// check the right side of slash
			// for pattern "number" and "comma"
			output = "every " + CheckResult.Input + " day of month" + output
		}
	}
	return output
}

func explainMonth(c *CheckResult) string {
	explanation := "in"
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
		explanation += " " + explainMonthSlach(c)
	}
	return explanation
}

func explainMonthSlach(c *CheckResult) string {
	month := c.Input
	var output string
	slice := strings.Split(month, "/")
	for i, v := range slice {
		CheckResult := CheckFieldeReg(v, "month")
		// check the left side of slash
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
			// check the right side of slash
			// for pattern "number" and "comma"
			output = "every " + CheckResult.Input + " month" + output
		}
	}
	return output
}

func explainDayWeek(c *CheckResult) string {
	explanation := "on"
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
		explanation += " " + explainDayWeekSlach(c)
	}
	return explanation
}

func explainDayWeekSlach(c *CheckResult) string {
	dayW := c.Input
	var output string
	slice := strings.Split(dayW, "/")
	for i, v := range slice {
		CheckResult := CheckFieldeReg(v, "dayW")
		// check the left side of slash
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
			// check the right side of slash
			// for pattern "number" and "comma"
			output = "every " + CheckResult.Input + " day of week" + output
		}
	}
	return output
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
