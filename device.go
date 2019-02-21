package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
