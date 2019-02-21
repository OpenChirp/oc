package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/openchirp/framework/pubsub"

	"github.com/openchirp/framework/rest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version string = "1.1"
)

const (
	columnPadding = 3
)

var host rest.Host

func main() {
	viper.SetConfigName("occonfig")         // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AddConfigPath("/etc/oc/")         // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/oc") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.oc")        // call multiple times to add many search paths

	var cmdUser = &cobra.Command{
		Use:   "user",
		Short: "Manage the user account",
	}

	var cmdUserInfo = &cobra.Command{
		Use:   "info",
		Short: "Fetch user info",
		Args:  cobra.NoArgs,
		Run:   userInfo,
	}

	var cmdUserLs = &cobra.Command{
		Use:   "ls",
		Short: "Fetch the list of all user",
		Args:  cobra.NoArgs,
		Run:   userLs,
	}

	var cmdUserCreate = &cobra.Command{
		Use:   "create <email> <password> [name]",
		Short: "Create new user",
		Args:  cobra.RangeArgs(2, 3),
		Run:   userCreate,
	}
	cmdUserCreate.Flags().Bool("occonfig", false, "Print out an oc config for the new user")

	var cmdService = &cobra.Command{
		Use:   "service",
		Short: "Manage a service",
	}

	var cmdServiceLs = &cobra.Command{
		Use:   "ls",
		Short: "List all services",
		Long:  `The ls command will print out all services with their respective IDs`,
		Args:  cobra.NoArgs,
		Run:   serviceLs,
	}

	var cmdServiceCreate = &cobra.Command{
		Use:   "create <name> <description>",
		Short: "Create a new service",
		Long: `The create command will create a new service with the given
name and description. Upon success, the service ID is printed.`,
		Args: cobra.ExactArgs(2),
		Run:  serviceCreate,
	}

	var cmdServiceRm = &cobra.Command{
		Use:   "rm <service_id>",
		Short: "Remove a new service",
		Args:  cobra.ExactArgs(1),
		Run:   serviceRm,
	}

	var cmdServiceMonitor = &cobra.Command{
		Use:   "monitor <service_id>",
		Short: "Monitor a service's pubsub traffic",
		Args:  cobra.ExactArgs(1),
		Run:   serviceMonitor,
	}

	var cmdServiceToken = &cobra.Command{
		Use:   "token",
		Short: "Manage the service auth token",
	}

	var cmdServiceTokenGenerate = &cobra.Command{
		Use:   "generate <service_id>",
		Short: "Generate a security token for the service",
		Args:  cobra.ExactArgs(1),
		Run:   serviceTokenGenerate,
	}
	cmdServiceTokenGenerate.Flags().Bool("env", false, "Print out all service environment variables to setup a service")

	var cmdServiceTokenRegenerate = &cobra.Command{
		Use:   "regenerate <service_id>",
		Short: "Regenerate a security token for the service",
		Args:  cobra.ExactArgs(1),
		Run:   serviceTokenRegenerate,
	}
	cmdServiceTokenRegenerate.Flags().Bool("env", false, "Print out all service environment variables to setup a service")

	var cmdServiceTokenRm = &cobra.Command{
		Use:   "rm <service_id>",
		Short: "Remove the security token for the service",
		Args:  cobra.ExactArgs(1),
		Run:   serviceTokenRm,
	}

	var cmdDevice = &cobra.Command{
		Use:   "device",
		Short: "Manage a device",
	}

	var cmdDeviceMonitor = &cobra.Command{
		Use:   "monitor <device_id>",
		Short: "Monitor a device's pubsub traffic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deviceID := args[0]

			device, err := host.RequestDeviceInfo(deviceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to fetch device information:", err)
				os.Exit(1)
			}

			client, err := pubsub.NewMQTTClient(
				viper.GetString("mqtt-server"),
				viper.GetString("auth-id"),
				viper.GetString("auth-token"),
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
		},
	}

	var cmdMonitor = &cobra.Command{
		Use:   "monitor <topics...>",
		Short: "Monitor any mqtt topic",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := pubsub.NewMQTTClient(
				viper.GetString("mqtt-server"),
				viper.GetString("auth-id"),
				viper.GetString("auth-token"),
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
			for _, t := range args {
				fmt.Println("Subscribing to", t)
				client.Subscribe(t, onMessage)
			}

			/* Wait on a signal */
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
			<-signals
		},
	}

	var rootCmd = &cobra.Command{Use: "oc", Version: version}

	// oc
	rootCmd.AddCommand(cmdService)
	rootCmd.AddCommand(cmdDevice)
	rootCmd.AddCommand(cmdUser)
	rootCmd.AddCommand(cmdMonitor)
	// oc service
	cmdService.AddCommand(cmdServiceLs)
	cmdService.AddCommand(cmdServiceCreate)
	cmdService.AddCommand(cmdServiceRm)
	cmdService.AddCommand(cmdServiceMonitor)
	//oc device
	cmdDevice.AddCommand(cmdDeviceMonitor)
	// oc service token
	cmdService.AddCommand(cmdServiceToken)
	cmdServiceToken.AddCommand(cmdServiceTokenGenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRegenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRm)
	// oc user
	cmdUser.AddCommand(cmdUserInfo)
	cmdUser.AddCommand(cmdUserCreate)
	cmdUser.AddCommand(cmdUserLs)

	rootCmd.PersistentFlags().StringP("framework-server", "s", "http://localhost", "Specifies the framework server")
	rootCmd.PersistentFlags().StringP("mqtt-server", "m", "tcp://localhost:1883", "Specifies the mqtt server")
	rootCmd.PersistentFlags().StringP("auth-id", "i", "", "The authentication ID to use with the framework server")
	rootCmd.PersistentFlags().StringP("auth-token", "t", "", "The authentication token to use with the framework server")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	viper.BindPFlag("framework-server", rootCmd.PersistentFlags().Lookup("framework-server"))
	viper.BindPFlag("mqtt-server", rootCmd.PersistentFlags().Lookup("mqtt-server"))
	viper.BindPFlag("auth-id", rootCmd.PersistentFlags().Lookup("auth-id"))
	viper.BindPFlag("auth-token", rootCmd.PersistentFlags().Lookup("auth-token"))
	viper.BindEnv("framework-server", "FRAMEWORK_SERVER")
	viper.BindEnv("auth-id", "AUTH_ID")
	viper.BindEnv("auth-token", "AUTH_TOKEN")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// continue on
		case viper.ConfigParseError:
			// Handle errors reading the config file
			fmt.Fprintln(os.Stderr, "Failed to parse config file:", err)
			os.Exit(1)
		}
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("verbose"); v {
			fmt.Fprintln(os.Stderr, "Framework Server:", viper.GetString("framework-server"))
			fmt.Fprintln(os.Stderr, "MQTT Server:", viper.GetString("mqtt-server"))
			fmt.Fprintln(os.Stderr, "Auth ID:", viper.GetString("auth-id"))
		}
		host = rest.NewHost(viper.GetString("framework-server"))
		host.Login(viper.GetString("auth-id"), viper.GetString("auth-token"))
	}
	rootCmd.Execute()
}
