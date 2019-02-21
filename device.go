package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func deviceMonitor(cmd *cobra.Command, args []string) {
	deviceID := args[0]

	device, err := host.RequestDeviceInfo(deviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch device information:", err)
		os.Exit(1)
	}

	monitor(device.Pubsub.Topic + "/#")
}

func deviceLs(cmd *cobra.Command, args []string) {
	devices, err := host.DeviceAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch devices:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, columnPadding, ' ', 0)
	fmt.Fprintln(w, "NAME\tOWNER\tID\tTOPIC\t")
	for _, d := range devices {
		fmt.Fprintf(w, "%s\t", d.Name)
		fmt.Fprintf(w, "%s (%s)\t", d.Owner.Name, d.Owner.Email)
		fmt.Fprintf(w, "%s\t", d.ID)
		fmt.Fprintf(w, "%s\t", d.Pubsub.Topic)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
