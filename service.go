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

func serviceLs(cmd *cobra.Command, args []string) {
	services, err := host.ServiceList()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch services:", err)
		os.Exit(1)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, columnPadding, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tID\tOWNER NAME\tOWNER EMAIL\t")
	for _, s := range services {
		fmt.Fprintf(w, "%s\t", s.Name)
		fmt.Fprintf(w, "%s\t", s.Description)
		fmt.Fprintf(w, "%s\t", s.ID)
		fmt.Fprintf(w, "%s\t", s.Owner.Name)
		fmt.Fprintf(w, "%s\t", s.Owner.Email)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}

func serviceCreate(cmd *cobra.Command, args []string) {
	name := args[0]
	description := args[1]
	s, err := host.ServiceCreate(name, description, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create service:", err)
		os.Exit(1)
	}
	fmt.Println(s.ID)
}

func serviceRm(cmd *cobra.Command, args []string) {
	serviceID := args[0]

	err := host.ServiceDelete(serviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to delete service:", err)
		os.Exit(1)
	}
}

func serviceMonitor(cmd *cobra.Command, args []string) {
	serviceID := args[0]

	service, err := host.ServiceGet(serviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch service information:", err)
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

	serviceTopic := service.Pubsub.Topic + "/#"
	fmt.Println("Subscribing to", serviceTopic)
	client.Subscribe(serviceTopic, func(topic string, payload []byte) {
		fmt.Println(topic, string(payload))
	})

	/* Wait on a signal */
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
}

func serviceTokenGenerate(cmd *cobra.Command, args []string) {
	serviceID := args[0]

	token, err := host.ServiceTokenGenerate(serviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to generate token:", err)
		os.Exit(1)
	}
	if envTrue, _ := cmd.Flags().GetBool("env"); envTrue {
		fmt.Printf("FRAMEWORK_SERVER=\"%s\"\n", frameworkHost)
		fmt.Printf("MQTT_SERVER=\"%s\"\n", mqttServer)
		fmt.Printf("SERVICE_ID=\"%s\"\n", serviceID)
		fmt.Printf("SERVICE_TOKEN=\"%s\"\n", token)
	} else {
		fmt.Println(token)
	}
}

func serviceTokenRegenerate(cmd *cobra.Command, args []string) {
	serviceID := args[0]

	token, err := host.ServiceTokenRegenerate(serviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to regenerate token:", err)
		os.Exit(1)
	}
	if envTrue, _ := cmd.Flags().GetBool("env"); envTrue {
		fmt.Printf("FRAMEWORK_SERVER=\"%s\"\n", frameworkHost)
		fmt.Printf("MQTT_SERVER=\"%s\"\n", mqttServer)
		fmt.Printf("SERVICE_ID=\"%s\"\n", serviceID)
		fmt.Printf("SERVICE_TOKEN=\"%s\"\n", token)
	} else {
		fmt.Println(token)
	}
}

func serviceTokenRm(cmd *cobra.Command, args []string) {
	serviceID := args[0]

	err := host.ServiceTokenDelete(serviceID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to delete token:", err)
		os.Exit(1)
	}
}
