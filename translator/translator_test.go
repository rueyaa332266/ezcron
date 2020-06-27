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
	gotList := []string{"Jan", "FEB", "March", "APr", "MaY", "june", "July", "Aug", "Sep", "Oct", "Nov", "Dec"}
	wantList := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	for i := range gotList {
		got := MonthToNum(gotList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestWeekDayToNum(t *testing.T) {
	gotList := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	wantList := []string{"0", "1", "2", "3", "4", "5", "6"}
	for i := range gotList {
		got := WeekDayToNum(gotList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestOrdinalFromStr(t *testing.T) {
	gotList := []string{"1", "2", "3", "4"}
	wantList := []string{"1st", "2nd", "3rd", "4th"}
	for i := range gotList {
		got := OrdinalFromStr(gotList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("Got: %s, but want: %s", got, want)
		}
	}
}

func TestAddZeorforTenDigit(t *testing.T) {
	gotList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	wantList := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09"}
	for i := range gotList {
		got := AddZeorforTenDigit(gotList[i])
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
