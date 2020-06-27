package cmd

import "testing"

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
