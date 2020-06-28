package translator

import (
	"reflect"
	"testing"
)

func TestMatchCronReg(t *testing.T) {
	gotValid, gotCheckResults := MatchCronReg("1 1-2 3,4 */1 *")
	want := map[string]*CheckResult{
		"minute": {true, "1", "number"},
		"hour":   {true, "1-2", "hyphen"},
		"dayM":   {true, "3,4", "comma"},
		"month":  {true, "*/1", "slash"},
		"dayW":   {true, "*", "asterisk"},
	}
	if gotValid != true {
		t.Errorf("got: %t; want: %t", gotValid, true)
	}
	for key := range want {
		if !reflect.DeepEqual(gotCheckResults[key], want[key]) {
			t.Errorf("Error in field: %s", key)
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

func TestCheckMinuteReg(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckMinuteReg(checkList[i])
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckHourReg(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckHourReg(checkList[i])
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckDayMonthReg(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckDayMonthReg(checkList[i])
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckMonthReg(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckMonthReg(checkList[i])
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
		}
	}
}

func TestCheckDayWeekReg(t *testing.T) {
	checkList := []string{"*", "1", "2,3,4", "1-2", "*/1"}
	wantList := []CheckResult{
		{true, "*", "asterisk"},
		{true, "1", "number"},
		{true, "2,3,4", "comma"},
		{true, "1-2", "hyphen"},
		{true, "*/1", "slash"},
	}
	for i := range checkList {
		got := CheckDayWeekReg(checkList[i])
		want := &wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error in pattern: %s", wantList[i].Pattern)
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
