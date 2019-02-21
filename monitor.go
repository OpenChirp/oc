package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/openchirp/framework/pubsub"
)

func monitor(topics ...string) {
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

	onMessage := func(topic string, payload []byte) {
		fmt.Println(topic, string(payload))
	}
	for _, t := range topics {
		fmt.Println("Subscribing to", t)
		client.Subscribe(t, onMessage)
	}

	/* Wait on a signal */
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
}
