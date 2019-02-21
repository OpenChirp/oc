package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/openchirp/framework/pubsub"
	"github.com/spf13/cobra"
)

func deviceMonitor(cmd *cobra.Command, args []string) {
	deviceID := args[0]

	device, err := host.RequestDeviceInfo(deviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch device information:", err)
		os.Exit(1)
	}

	client, err := pubsub.NewMQTTClient(
		mqttServer,
		authID,
		authToken,
		pubsub.QoSExactlyOnce,
		false,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to MQTT server:", err)
		os.Exit(1)
	}
	defer client.Disconnect()

	deviceTopic := device.Pubsub.Topic + "/#"
	fmt.Println("Subscribing to", deviceTopic)
	client.Subscribe(deviceTopic, func(topic string, payload []byte) {
		fmt.Println(topic, string(payload))
	})

	/* Wait on a signal */
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
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
