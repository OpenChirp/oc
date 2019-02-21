package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func groupCreate(cmd *cobra.Command, args []string) {
	name := args[0]

	if err := host.GroupCreate(name); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create group:", err)
		os.Exit(1)
	}
}

func groupLs(cmd *cobra.Command, args []string) {
	groups, err := host.GroupAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch groups:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, columnPadding, ' ', 0)
	fmt.Fprintln(w, "NAME\tID\t")
	for _, g := range groups {
		fmt.Fprintf(w, "%s\t", g.Name)
		fmt.Fprintf(w, "%s\t", g.ID)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
