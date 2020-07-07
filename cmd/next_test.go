package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestNextExecTime(t *testing.T) {
	const layout = "2000-01-01 01:01:01"
	checkList := []string{"* * * * *", "* * * *"}
	for i := range checkList {
		got, err := nextExecTime(checkList[i])
		if err != nil {
			want := "Invalid syntax"
			if err.Error() != want {
				t.Errorf("Want: %v", want)
			}
		} else {
			want := time.Now().Add(1 * time.Minute).Format(layout)
			if got.Format(layout) != want {
				t.Errorf("Got: %v, but want: %s", got, want)
			}
		}
	}
}

func Example_getNextTime() {
	var cmd *cobra.Command
	getNextTime(cmd, []string{"*", "*", "*", "*"})

	// Output:
	// Invalid syntax
}
