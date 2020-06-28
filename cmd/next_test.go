package cmd

import (
	"testing"
)

func TestNextExecTime(t *testing.T) {
	got := nextExecTime("* * * * *")
	if got != nil {
		t.Errorf("Error when check next execute time")
	}
}
