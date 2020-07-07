package cmd

import (
	"github.com/spf13/cobra"
)

func Example_translate() {
	var cmd *cobra.Command
	checkList := [][]string{
		{"*", "*", "*", "*", "*"},
		{"*", "*", "*", "*"},
	}
	for i := range checkList {
		translate(cmd, checkList[i])
	}

	// Output:
	// At every minute
	// Invalid syntax
}
