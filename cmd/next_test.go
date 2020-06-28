package cmd

import (
	"testing"
	"time"
)

func TestNextExecTime(t *testing.T) {
	const layout = "2000-01-01 01:01:01"
	got, _ := nextExecTime("* * * * *")
	want := time.Now().Add(1 * time.Minute).Format(layout)

	if got.Format(layout) != want {
		t.Errorf("Got: %v, but want: %s", got, want)
	}
}
