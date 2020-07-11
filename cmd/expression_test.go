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
		day := translator.OrdinalDay(strconv.Itoa(d))
		suggest := prompt.Suggest{Text: day + "_day", Description: "of month"}
		want = append(want, suggest)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error when making month day suggest")
	}
}

func TestmakeMonthdayNumberSuggest(t *testing.T) {
	var want []prompt.Suggest
	checkList := [][]string{
		dayList[:28],
		dayList[:30],
		dayList,
	}
	for d := 1; d < 32; d++ {
		day := translator.OrdinalDay(strconv.Itoa(d))
		suggest := prompt.Suggest{Text: day, Description: "default at 00:00"}
		want = append(want, suggest)
	}
	wantList := [][]prompt.Suggest{
		want[:28],
		want[:30],
		want,
	}
	for i := range checkList {
		got := makeMonthdayNumberSuggest(checkList[i])
		if !reflect.DeepEqual(got, wantList[i]) {
			t.Errorf("Error when making month day suggest")
		}
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

func Example_executeTimeSchedule() {
	checkList := [][]string{
		{"Time_schedule", "at", "00:00"},
		{"Time_schedule", "every_minute", "1_minute"},
		{"Time_schedule", "every_hour", "1_hour"},
		{"Time_schedule", "test"},
		{"Time_schedule", "at", "test"},
	}
	for i := range checkList {
		executeTimeSchedule(checkList[i])
	}

	// Output:
	// 0 0 * * *
	// */1 * * * *
	// * */1 * * *
	// Invalid time schedule
	// Invalid time schedule
}

func Example_executeDailySchedule() {
	checkList := [][]string{
		{"Daily_schedule", "every_day"},
		{"Daily_schedule", "every_day_at", "01:01"},
		{"Daily_schedule", "every_day_at", "01:01"},
		{"Daily_schedule"},
		{"Daily_schedule", "test"},
	}
	for i := range checkList {
		executeDailySchedule(checkList[i])
	}

	// Output:
	// 0 0 */1 * *
	// 1 1 */1 * *
	// 1 1 */1 * *
	// Invalid daily schedule
	// Invalid daily schedule
}

func Example_executeWeeklySchedule() {
	checkList := [][]string{
		{"Weekly_schedule", "on_every", "Sunday"},
		{"Weekly_schedule", "on_every", "Sunday", "at", "01:01"},
		{"Weekly_schedule", "test"},
		{"Weekly_schedule", "on_every", "test"},
		{"Weekly_schedule", "on_every", "Sunday", "at", "test"},
	}
	for i := range checkList {
		executeWeeklySchedule(checkList[i])
	}

	// Output:
	// 0 0 * * 0
	// 1 1 * * 0
	// Invalid weekly schedule
	// Invalid weekly schedule
	// Invalid weekly schedule
}

func Example_executeMonthlySchedule() {
	checkList := [][]string{
		{"Monthly_schedule", "on", "1st_day", "of_every_month"},
		{"Monthly_schedule", "on", "1st_day", "of_every", "2_month"},
		{"Monthly_schedule", "on", "1st_day", "of_every_month", "at", "01:01"},
		{"Monthly_schedule", "on", "1st_day", "of_every", "2_month", "at", "01:01"},
		{"Monthly_schedule", "test"},
		{"Monthly_schedule", "on", "test", "of_every_month"},
		{"Monthly_schedule", "on", "1st_day", "of_every", "test"},
		{"Monthly_schedule", "on", "1st_day", "of_every_month", "at", "test"},
		{"Monthly_schedule", "on", "1st_day", "of_every", "2_month", "at", "test"},
	}
	for i := range checkList {
		executeMonthlySchedule(checkList[i])
	}

	// Output:
	// 0 0 1 */1 *
	// 0 0 1 */2 *
	// 1 1 1 */1 *
	// 1 1 1 */2 *
	// Invalid monthly schedule
	// Invalid monthly schedule
	// Invalid monthly schedule
	// Invalid monthly schedule
	// Invalid monthly schedule
}

func Example_executeYearlySchedule() {
	checkList := [][]string{
		{"Yearly_schedule", "in_every", "January", "1st"},
		{"Yearly_schedule", "in_every", "January", "1st", "at", "01:01"},
		{"Yearly_schedule", "in_every", "January", "test"},
		{"Yearly_schedule", "in_every", "January", "1st", "at", "test"},
		{"Yearly_schedule", "test"},
	}
	for i := range checkList {
		executeYearlySchedule(checkList[i])
	}

	// Output:
	// 0 0 1 1 *
	// 1 1 1 1 *
	// Invalid yearly schedule
	// Invalid yearly schedule
	// Invalid yearly schedule
}

func TestCompleteTimeSchedule(t *testing.T) {
	// var gotSuggest []prompt.Suggest
	// var gotSub string
	checkList := [][]string{
		{"Time_schedule:", "a"},
		{"Time_schedule:", "at", "0"},
		{"Time_schedule:", "every_minute", "1"},
		{"Time_schedule:", "every_hour", "2"},
	}
	wantSuggestList := [][]prompt.Suggest{
		{{Text: "at", Description: "__:__ every day"}, {Text: "every_minute", Description: "per minute"}, {Text: "every_hour", Description: "per hour"}},
		makeTimeSuggest("time"),
		makeTimeSuggest("minute"),
		makeTimeSuggest("hour"),
	}
	wantSubtList := []string{"a", "0", "1", "2"}
	for i := range checkList {
		gotSuggest, gotSub := completeTimeSchedule(checkList[i])
		wantSuggest := wantSuggestList[i]
		wantSub := wantSubtList[i]
		if !reflect.DeepEqual(gotSuggest, wantSuggest) {
			t.Errorf("Got: %v, but want: %s", gotSuggest, wantSuggest)
		}
		if !reflect.DeepEqual(gotSub, wantSub) {
			t.Errorf("Got: %v, but want: %s", gotSub, wantSub)
		}
	}
}

func TestCompleteDailySchedule(t *testing.T) {
	// var gotSuggest []prompt.Suggest
	// var gotSub string
	checkList := [][]string{
		{"Daily_schedule:", "e"},
		{"Daily_schedule:", "every_day_at", "1"},
	}
	wantSuggestList := [][]prompt.Suggest{
		{{Text: "every_day", Description: "every day at 00:00"}, {Text: "every_day_at", Description: "every day at __:__"}},
		makeTimeSuggest("time"),
	}
	wantSubtList := []string{"e", "1"}
	for i := range checkList {
		gotSuggest, gotSub := completeDailySchedule(checkList[i])
		wantSuggest := wantSuggestList[i]
		wantSub := wantSubtList[i]
		if !reflect.DeepEqual(gotSuggest, wantSuggest) {
			t.Errorf("Got: %v, but want: %s", gotSuggest, wantSuggest)
		}
		if !reflect.DeepEqual(gotSub, wantSub) {
			t.Errorf("Got: %v, but want: %s", gotSub, wantSub)
		}
	}
}

func TestCompleteWeeklySchedule(t *testing.T) {
	// var gotSuggest []prompt.Suggest
	// var gotSub string
	checkList := [][]string{
		{"Weekly_schedule:", "o"},
		{"Weekly_schedule:", "on_every", "S"},
		{"Weekly_schedule:", "on_every", "Sunday", "a"},
		{"Weekly_schedule:", "on_every", "Sunday", "at", "1"},
	}
	wantSuggestList := [][]prompt.Suggest{
		{{Text: "on_every", Description: "weekday"}},
		makeWeekdaySuggest(),
		{{Text: "at", Description: "__:__"}},
		makeTimeSuggest("time"),
	}
	wantSubtList := []string{"o", "S", "a", "1"}
	for i := range checkList {
		gotSuggest, gotSub := completeWeeklySchedule(checkList[i])
		wantSuggest := wantSuggestList[i]
		wantSub := wantSubtList[i]
		if !reflect.DeepEqual(gotSuggest, wantSuggest) {
			t.Errorf("Got: %v, but want: %s", gotSuggest, wantSuggest)
		}
		if !reflect.DeepEqual(gotSub, wantSub) {
			t.Errorf("Got: %v, but want: %s", gotSub, wantSub)
		}
	}
}
func TestCompleteMonthlySchedule(t *testing.T) {
	// var gotSuggest []prompt.Suggest
	// var gotSub string
	checkList := [][]string{
		{"Monthly_schedule:", "o"},
		{"Monthly_schedule:", "on", "1"},
		{"Monthly_schedule:", "on", "1st_day", "o"},
		{"Monthly_schedule:", "on", "1st_day", "of_every_month", "a"},
		{"Monthly_schedule:", "on", "1st_day", "of_every_month", "at", "1"},
		{"Monthly_schedule:", "on", "1st_day", "of_every", "1"},
		{"Monthly_schedule:", "on", "1st_day", "of_every", "1_month", "a"},
		{"Monthly_schedule:", "on", "1st_day", "of_every", "1_month", "at", "1"},
	}
	wantSuggestList := [][]prompt.Suggest{
		{{Text: "on", Description: "monthday"}},
		makeMonthdaySuggest(),
		{{Text: "of_every_month", Description: "per month, default at 00:00"}, {Text: "of_every", Description: "period of month"}},
		{{Text: "at", Description: "__:__"}},
		makeTimeSuggest("time"),
		makeMonthNumSuggest(),
		{{Text: "at", Description: "__:__"}},
		makeTimeSuggest("time"),
	}
	wantSubtList := []string{"o", "1", "o", "a", "1", "1", "a", "1"}
	for i := range checkList {
		gotSuggest, gotSub := completeMonthlySchedule(checkList[i])
		wantSuggest := wantSuggestList[i]
		wantSub := wantSubtList[i]
		if !reflect.DeepEqual(gotSuggest, wantSuggest) {
			t.Errorf("Got: %v, but want: %s", gotSuggest, wantSuggest)
		}
		if !reflect.DeepEqual(gotSub, wantSub) {
			t.Errorf("Got: %v, but want: %s", gotSub, wantSub)
		}
	}
}

func TestCompletYearlySchedule(t *testing.T) {
	var day []string
	for i := 1; i < 32; i++ {
		day = append(day, strconv.Itoa(i))
	}
	day28 := day[:28]
	day30 := day[:30]
	f := func(src []string) []prompt.Suggest {
		var suggests []prompt.Suggest
		for _, v := range src {
			suggest := prompt.Suggest{Text: translator.OrdinalDay(v), Description: "default at 00:00"}
			suggests = append(suggests, suggest)
		}
		return suggests
	}
	checkList := [][]string{
		{"Yearly_schedule:", "i"},
		{"Yearly_schedule:", "in_every", "J"},
		{"Yearly_schedule:", "in_every", "January", "1"},
		{"Yearly_schedule:", "in_every", "February", "1"},
		{"Yearly_schedule:", "in_every", "April", "1"},
		{"Yearly_schedule:", "in_every", "January", "1st", "a"},
		{"Yearly_schedule:", "in_every", "January", "1st", "at", "1"},
	}
	wantSuggestList := [][]prompt.Suggest{
		{{Text: "in_every", Description: "month_day"}},
		makeMonthSuggest(),
		f(day),
		f(day28),
		f(day30),
		{{Text: "at", Description: "__:__"}},
		makeTimeSuggest("time"),
	}
	wantSubtList := []string{"i", "J", "1", "1", "1", "a", "1"}
	for i := range checkList {
		gotSuggest, gotSub := completeYearlySchedule(checkList[i])
		wantSuggest := wantSuggestList[i]
		wantSub := wantSubtList[i]
		if !reflect.DeepEqual(gotSuggest, wantSuggest) {
			t.Errorf("Got: %v, but want: %s", gotSuggest, wantSuggest)
		}
		if !reflect.DeepEqual(gotSub, wantSub) {
			t.Errorf("Got: %v, but want: %s", gotSub, wantSub)
		}
	}
}

func TestMakeSuggestByPreWord(t *testing.T) {
	checkList := []string{"at", "every_minute", "every_hour"}
	wantList := [][]prompt.Suggest{makeTimeSuggest("time"), makeTimeSuggest("minute"), makeTimeSuggest("hour")}
	for i := range checkList {
		got := makeSuggestByPreWord(checkList[i])
		want := wantList[i]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got: %v, but want: %s", got, want)
		}
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
