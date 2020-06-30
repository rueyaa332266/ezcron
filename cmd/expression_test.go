package cmd

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/rueyaa332266/ezcron/translator"
)

func TestMakeTimeSuggest_Time(t *testing.T) {
	var want []prompt.Suggest
	got := makeTimeSuggest("time")
	for h := 0; h <= 23; h++ {
		for m := 0; m <= 59; m++ {
			min := translator.AddZeorforTenDigit(strconv.Itoa(m))
			hour := translator.AddZeorforTenDigit(strconv.Itoa(h))
			suggest := prompt.Suggest{Text: hour + ":" + min}
			want = append(want, suggest)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making time suggest")
	}
}

func TestMakeTimeSuggest_Minute(t *testing.T) {
	var want []prompt.Suggest
	got := makeTimeSuggest("minute")
	for m := 1; m < 61; m++ {
		min := strconv.Itoa(m)
		suggest := prompt.Suggest{Text: min + "_minute"}
		want = append(want, suggest)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making minute suggest")
	}
}

func TestMakeTimeSuggest_Hour(t *testing.T) {
	var want []prompt.Suggest
	got := makeTimeSuggest("hour")
	for h := 1; h < 25; h++ {
		hour := strconv.Itoa(h)
		suggest := prompt.Suggest{Text: hour + "_hour"}
		want = append(want, suggest)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making hour suggest")
	}
}

func TestMakeWeekdaySuggest(t *testing.T) {
	want := []prompt.Suggest{
		{Text: "Sunday", Description: "default at 00:00"},
		{Text: "Monday", Description: "default at 00:00"},
		{Text: "Tuesday", Description: "default at 00:00"},
		{Text: "Wednesday", Description: "default at 00:00"},
		{Text: "Thursday", Description: "default at 00:00"},
		{Text: "Friday", Description: "default at 00:00"},
		{Text: "Saturday", Description: "default at 00:00"},
	}
	got := makeWeekdaySuggest()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making week day suggest")
	}
}

func TestMakeMonthdaySuggest(t *testing.T) {
	var want []prompt.Suggest
	got := makeMonthdaySuggest()
	for d := 1; d < 32; d++ {
		day := translator.OrdinalFromStr(strconv.Itoa(d))
		suggest := prompt.Suggest{Text: day + "_day", Description: "of month"}
		want = append(want, suggest)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making month day suggest")
	}
}

func TestMakeMonthNumSuggest(t *testing.T) {
	want := []prompt.Suggest{
		{Text: "1_month", Description: "default at 00:00"},
		{Text: "2_month", Description: "default at 00:00"},
		{Text: "3_month", Description: "default at 00:00"},
		{Text: "4_month", Description: "default at 00:00"},
		{Text: "5_month", Description: "default at 00:00"},
		{Text: "6_month", Description: "default at 00:00"},
		{Text: "7_month", Description: "default at 00:00"},
		{Text: "8_month", Description: "default at 00:00"},
		{Text: "9_month", Description: "default at 00:00"},
		{Text: "10_month", Description: "default at 00:00"},
		{Text: "11_month", Description: "default at 00:00"},
		{Text: "12_month", Description: "default at 00:00"},
	}
	got := makeMonthNumSuggest()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making month number suggest")
	}
}

func TestMakeMonthSuggest(t *testing.T) {
	want := []prompt.Suggest{
		{Text: "January"},
		{Text: "February"},
		{Text: "March"},
		{Text: "April"},
		{Text: "May"},
		{Text: "June"},
		{Text: "July"},
		{Text: "August"},
		{Text: "September"},
		{Text: "October"},
		{Text: "November"},
		{Text: "December"},
	}
	got := makeMonthSuggest()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making month suggest")
	}
}

func TestComplete(t *testing.T) {
	in := prompt.Document{Text: "Wait aaa aaa"}
	got := completer(in)
	want := scheduleTypeSuggest
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got: %v, but want: %s", got, want)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"foo", "bar"}
	checkList := []string{"foo", "buzz"}
	wantList := []bool{true, false}
	for i := range checkList {
		got := contains(slice, checkList[i])
		want := wantList[i]
		if got != want {
			t.Errorf("got: %t; want: %t", got, want)
		}
	}
}
