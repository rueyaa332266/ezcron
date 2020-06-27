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
	for key, _ := range want {
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
