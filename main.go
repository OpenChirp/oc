package main

import (
	"fmt"
	"os"

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

var (
	frameworkHost string
	mqttServer    string
	authID        string
	authToken     string
)
var host rest.Host

func printConfig(tomlFormat bool) {
	if tomlFormat {
		fmt.Printf("framework-server = \"%s\"\n", frameworkHost)
		fmt.Printf("mqtt-server = \"%s\"\n", mqttServer)
		fmt.Printf("auth-id = \"%s\"\n", authID)
		fmt.Printf("auth-token = \"%s\"\n", authToken)
	} else {
		fmt.Fprintln(os.Stderr, "Framework Server:", frameworkHost)
		fmt.Fprintln(os.Stderr, "MQTT Server:", mqttServer)
		fmt.Fprintln(os.Stderr, "Auth ID:", authID)
	}
}

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
	cmdUserCreate.Flags().BoolP("occonfig", "c", false, "Print out an oc config for the new user")

	var cmdGroup = &cobra.Command{
		Use:   "group",
		Short: "Manage groups",
	}

	var cmdGroupCreate = &cobra.Command{
		Use:   "create <name>",
		Short: "Create new group",
		Args:  cobra.ExactArgs(1),
		Run:   groupCreate,
	}

	var cmdGroupLs = &cobra.Command{
		Use:   "ls",
		Short: "List all groups",
		Args:  cobra.NoArgs,
		Run:   groupLs,
	}

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
		Run:   deviceMonitor,
	}

	var cmdDeviceLs = &cobra.Command{
		Use:   "ls",
		Short: "List all devices",
		Args:  cobra.NoArgs,
		Run:   deviceLs,
	}

	var cmdMonitor = &cobra.Command{
		Use:   "monitor <topics...>",
		Short: "Monitor any mqtt topic",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			monitor(args...)
		},
	}

	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Print out current config settings in use",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			tomlFormat, _ := cmd.Flags().GetBool("occonfig")
			printConfig(tomlFormat)
		},
	}
	cmdConfig.Flags().BoolP("occonfig", "c", false, "Print out an oc config in toml format")

	var rootCmd = &cobra.Command{Use: "oc", Version: version}

	// oc
	rootCmd.AddCommand(cmdService)
	rootCmd.AddCommand(cmdDevice)
	rootCmd.AddCommand(cmdUser)
	rootCmd.AddCommand(cmdGroup)
	rootCmd.AddCommand(cmdMonitor)
	rootCmd.AddCommand(cmdConfig)
	// oc service
	cmdService.AddCommand(cmdServiceLs)
	cmdService.AddCommand(cmdServiceCreate)
	cmdService.AddCommand(cmdServiceRm)
	cmdService.AddCommand(cmdServiceMonitor)
	//oc device
	cmdDevice.AddCommand(cmdDeviceMonitor)
	cmdDevice.AddCommand(cmdDeviceLs)
	// oc service token
	cmdService.AddCommand(cmdServiceToken)
	cmdServiceToken.AddCommand(cmdServiceTokenGenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRegenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRm)
	// oc user
	cmdUser.AddCommand(cmdUserInfo)
	cmdUser.AddCommand(cmdUserCreate)
	cmdUser.AddCommand(cmdUserLs)
	// oc group
	cmdGroup.AddCommand(cmdGroupCreate)
	cmdGroup.AddCommand(cmdGroupLs)

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
		frameworkHost = viper.GetString("framework-server")
		mqttServer = viper.GetString("mqtt-server")
		authID = viper.GetString("auth-id")
		authToken = viper.GetString("auth-token")

		if v, _ := cmd.Flags().GetBool("verbose"); v {
			printConfig(false)
		}
		host = rest.NewHost(frameworkHost)
		host.Login(authID, authToken)
	}
	rootCmd.Execute()
}
