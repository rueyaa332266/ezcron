package translator

import (
	"reflect"
	"testing"
)

func TestMatchCronReg(t *testing.T) {
	checkList := []string{"1 1-2 3,4 */1 *", "* * * * * *"}
	validList := []bool{true, false}
	want := map[string]*CheckResult{
		"minute": {true, "1", "number"},
		"hour":   {true, "1-2", "hyphen"},
		"dayM":   {true, "3,4", "comma"},
		"month":  {true, "*/1", "slash"},
		"dayW":   {true, "*", "asterisk"},
	}
	for i := range checkList {
		gotValid, gotCheckResults := MatchCronReg(checkList[i])
		if gotValid != validList[i] {
			t.Errorf("got: %t; want: %t", gotValid, true)
		}
		// check only in valid pattern
		if i == 0 {
			for key := range want {
				if !reflect.DeepEqual(gotCheckResults[key], want[key]) {
					t.Errorf("Error in field: %s", key)
				}
			}
		}
	}
}

func TestCheckFieldPattern(t *testing.T) {
	asterisk := `\*`
	number := `([0-9]|[1-5][0-9])`
	comma := `(` + number + `\,){1,59}` + number
	hyphen := number + `\-` + number
	slash := `(` + asterisk + `|` + number + `|` + comma + `|` + hyphen + `)` + `\/` + `(` + number + `|` + comma + `)`
	patternList := map[string]string{"asterisk": asterisk, "number": number, "comma": comma, "hyphen": hyphen, "slash": slash}

	gotValid, gotPattern := checkFieldPattern("Invalid", patternList)
	if gotValid != false || gotPattern != "not match" {
		t.Errorf("Error in not match")
	}
	checkList := []string{"*", "1", "1,2", "1-2", "*/1"}
	wantList := []string{"asterisk", "number", "comma", "hyphen", "slash"}
	for i := range checkList {
		_, got := checkFieldPattern(checkList[i], patternList)
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestNumLogicFormat(t *testing.T) {
	type ckeck struct {
		input   string
		pattern string
	}
	type re struct {
		valid bool
		out   string
	}
	checkList := []ckeck{
		{"1-2", "hyphen"},
		{"2-1", "hyphen"},
		{"1,2,3", "comma"},
		{"1,2,2", "comma"},
		{"*/1,2,2", "slash"},
		{"*/2-1", "slash"},
	}
	wantList := []re{
		{true, "1-2"},
		{false, "2-1"},
		{true, "1,2,3"},
		{true, "1,2"},
		{true, "*/1,2"},
		{false, "*/2-1"},
	}
	for i := range checkList {
		gotValid, gotStr := numLogicFormat(checkList[i].input, checkList[i].pattern)
		wantValid := wantList[i].valid
		wantStr := wantList[i].out
		if gotValid != wantValid {
			t.Errorf("Error in pattern: %s", checkList[i].pattern)
		}
		if gotStr != wantStr {
			t.Errorf("Error in pattern: %s", checkList[i].pattern)
		}
	}
}

func TestCheckFieldeReg_Minute(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckFieldeReg(checkList[i], "minute")
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckFieldeReg_Hour(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckFieldeReg(checkList[i], "hour")
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckFieldeReg_DayMonth(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckFieldeReg(checkList[i], "dayM")
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckFieldeReg_Month(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckFieldeReg(checkList[i], "month")
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckFieldeReg_DayWeek(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckFieldeReg(checkList[i], "dayW")
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func ExampleExplain() {
	checkList := []string{
		"* * * * *",
		"30 12 * * *",
		"* * */1 * *",
		"* * * * */1",
		"* * 1 * 1",
	}
	for i := range checkList {
		_, checkResult := MatchCronReg(checkList[i])
		Explain(checkResult)
	}

	// Output:
	// At every minute
	// At every 12:30
	// At every minute, on every 1 day of month
	// At every minute, on every 1 day of week
	// At every minute, on day 1 of the month and on Monday
}

func TestExplainMinute(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1", "1/1,2", "1,2/1", "1-2/1"}
	wantList := []string{
		"At every minute",
		"At every 1st minute",
		"At every 2nd and 3rd and 4th minute",
		"At every 1 through 2 minute",
		"At every 1 minute",
		"At every 1,2 minute from minute 1 through minute 59",
		"At every 1 minute from minute 1 through minute 59 and minute 2 through minute 59",
		"At every 1 minute from minute 1 through minute 2",
	}
	for i := range checkList {
		got := explainMinute(CheckFieldeReg(checkList[i], "minute"))
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestExplainHour(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1", "1/1,2", "1,2/1", "1-2/1"}
	wantList := []string{
		"",
		"past 01:00",
		"past 02:00 and 03:00 and 04:00",
		"past from 01:00-02:00",
		"past every 1 hour",
		"past every 1 hour and every 2 hour from 1:00-24:00",
		"past every 1 hour from 01:00-24:00 and 02:00-24:00",
		"past every 1 hour from 01:00-02:00",
	}
	for i := range checkList {
		got := explainHour(CheckFieldeReg(checkList[i], "hour"))
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestExplainDayMonth(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1", "1/1,2", "1,2/1", "1-2/1"}
	wantList := []string{
		"",
		"on day 1 of the month",
		"on day 2 and 3 and 4 of the month",
		"between day 1 and 2 of the month",
		"on every 1 day of month",
		"on every 1,2 day of month from day 1 of the month",
		"on every 1 day of month from day 1 of the month and day 2 of the month",
		"on every 1 day of month between day 1 and day 2 of the month",
	}
	for i := range checkList {
		got := explainDayMonth(CheckFieldeReg(checkList[i], "dayM"))
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestExplainMonth(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1", "1/1,2", "1,2/1", "1-2/1"}
	wantList := []string{
		"",
		"in January",
		"in February and March and April",
		"in every month from January through February",
		"in every 1 month",
		"in every 1,2 month from January",
		"in every 1 month from January and February",
		"in every 1 month from January through February",
	}
	for i := range checkList {
		got := explainMonth(CheckFieldeReg(checkList[i], "month"))
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}
func TestExplainDayWeek(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1", "1/1,2", "1,2/1", "1-2/1"}
	wantList := []string{
		"",
		"on Monday",
		"on Tuesday and Wednesday and Thursday",
		"on every day from Monday through Tuesday",
		"on every 1 day of week",
		"on every 1,2 day of week from Monday",
		"on every 1 day of week from Monday and Tuesday",
		"on every 1 day of week between Monday and Tuesday",
	}
	for i := range checkList {
		got := explainDayWeek(CheckFieldeReg(checkList[i], "dayW"))
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}
func TestMonthToNum(t *testing.T) {
	checkList := []string{"Jan", "FEB", "March", "APr", "MaY", "june", "July", "Aug", "Sep", "Oct", "Nov", "Dec"}
	wantList := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	for i := range checkList {
		got := MonthToNum(checkList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestWeekDayToNum(t *testing.T) {
	checkList := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	wantList := []string{"0", "1", "2", "3", "4", "5", "6"}
	for i := range checkList {
		got := WeekDayToNum(checkList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestOrdinalFromStr(t *testing.T) {
	checkList := []string{"1", "2", "3", "4"}
	wantList := []string{"1st", "2nd", "3rd", "4th"}
	for i := range checkList {
		got := OrdinalFromStr(checkList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestAddZeorforTenDigit(t *testing.T) {
	checkList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	wantList := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09"}
	for i := range checkList {
		got := AddZeorforTenDigit(checkList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestUniqueSlice(t *testing.T) {
	slice := []string{"foo", "bar", "fizz", "foo", "fizz"}
	got := uniqueSlice(slice)
	want := []string{"foo", "bar", "fizz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got: %s, but want: %s", got, want)
	}
}
