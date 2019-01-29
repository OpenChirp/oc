package main

import (
	"fmt"
	"os"

	"github.com/openchirp/framework/rest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version string = "1.0"
)

func main() {
	var host rest.Host

	viper.SetConfigName("occonfig")         // name of config file (without extension)
	viper.AddConfigPath("/etc/oc/")         // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/oc") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.oc")        // call multiple times to add many search paths
	viper.AddConfigPath(".")                // optionally look for config in the working directory

	var cmdService = &cobra.Command{
		Use:   "service",
		Short: "Manage a service",
	}

	var cmdServiceCreate = &cobra.Command{
		Use:   "create <name> <description>",
		Short: "Create a new service",
		Long: `The create command will create a new service with the given
name and description. Upon success, the service ID is printed.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			description := args[1]
			s, err := host.ServiceCreate(name, description, nil, nil)
			if err != nil {
				fmt.Println("Failed to create service:", err)
				os.Exit(1)
			}
			fmt.Println(s.ID)
		},
	}

	var cmdServiceDelete = &cobra.Command{
		Use:   "delete <service_id>",
		Short: "Delete a new service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			err := host.ServiceDelete(serviceID)
			if err != nil {
				fmt.Println("Failed to delete service:", err)
				os.Exit(1)
			}
		},
	}

	var cmdServiceToken = &cobra.Command{
		Use:   "token",
		Short: "Manage the service auth token",
	}

	var cmdServiceTokenGenerate = &cobra.Command{
		Use:   "generate <service_id>",
		Short: "Generate a security token for the service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			token, err := host.ServiceTokenGenerate(serviceID)
			if err != nil {
				fmt.Println("Failed to generate token:", err)
				os.Exit(1)
			}
			fmt.Println(token)
		},
	}

	var cmdServiceTokenRegenerate = &cobra.Command{
		Use:   "regenerate <service_id>",
		Short: "Regenerate a security token for the service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			token, err := host.ServiceTokenRegenerate(serviceID)
			if err != nil {
				fmt.Println("Failed to regenerate token:", err)
				os.Exit(1)
			}
			fmt.Println(token)
		},
	}

	var cmdServiceTokenDelete = &cobra.Command{
		Use:   "delete <service_id>",
		Short: "Delete the security token for the service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			err := host.ServiceTokenDelete(serviceID)
			if err != nil {
				fmt.Println("Failed to delete token:", err)
				os.Exit(1)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "oc", Version: version}

	// oc
	rootCmd.AddCommand(cmdService)
	// oc service
	cmdService.AddCommand(cmdServiceCreate)
	cmdService.AddCommand(cmdServiceDelete)
	// oc service token
	cmdService.AddCommand(cmdServiceToken)
	cmdServiceToken.AddCommand(cmdServiceTokenGenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRegenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenDelete)

	rootCmd.PersistentFlags().StringP("framework-server", "s", "http://localhost", "Specifies the framework server")
	rootCmd.PersistentFlags().StringP("auth-id", "i", "", "The authentication ID to use with the framework server")
	rootCmd.PersistentFlags().StringP("auth-token", "t", "", "The authentication token to use with the framework server")
	viper.BindPFlag("framework-server", rootCmd.PersistentFlags().Lookup("framework-server"))
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
			panic(fmt.Errorf("Failed to parse config file: %v\n", err))
		}
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		host = rest.NewHost(viper.GetString("framework-server"))
		host.Login(viper.GetString("auth-id"), viper.GetString("auth-token"))
	}
	rootCmd.Execute()
}
