package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func userInfo(cmd *cobra.Command, args []string) {

	user, err := host.RequestUserInfo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch user info:", err)
		os.Exit(1)
	}
	fmt.Println("Name:", user.Name)
	fmt.Println("Email:", user.Email)
	fmt.Println("UserID:", user.UserID)
	fmt.Print("Groups: ")
	for _, g := range user.Groups {
		var access = "execute"
		if g.WriteAccess {
			access = "write"
		}
		fmt.Printf("%s-%s ", g.Name, access)
	}
	fmt.Println()
}

func userLs(cmd *cobra.Command, args []string) {
	users, err := host.UserAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch users:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, columnPadding, ' ', 0)
	fmt.Fprintln(w, "NAME\tUSERID\tEMAIL\tID\t")
	for _, u := range users {
		fmt.Fprintf(w, "%s\t", u.Name)
		fmt.Fprintf(w, "%s\t", u.UserID)
		fmt.Fprintf(w, "%s\t", u.Email)
		fmt.Fprintf(w, "%s\t", u.ID)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}

func userCreate(cmd *cobra.Command, args []string) {
	email := args[0]
	password := args[1]
	name := ""
	if len(args) > 2 {
		name = args[2]
	}
	if err := host.UserCreate(email, name, password); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create user:", err)
		os.Exit(1)
	}
	if occonfigTrue, _ := cmd.Flags().GetBool("occonfig"); occonfigTrue {
		fmt.Printf("framework-server = \"%s\"\n", frameworkHost)
		fmt.Printf("mqtt-server = \"%s\"\n", mqttServer)
		fmt.Printf("auth-id = \"%s\"\n", email)
		fmt.Printf("auth-token = \"%s\"\n", "MUST_GENERATE")
	}
}
