package main

import (
	"fmt"
	"os"

	"github.com/openchirp/framework/rest"

	"github.com/spf13/cobra"
)

const (
	version string = "1.0"
)

func main() {
	var frameworkServer string
	var authID string
	var authToken string
	var host rest.Host

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

	// service
	cmdService.AddCommand(cmdServiceCreate)
	cmdService.AddCommand(cmdServiceDelete)
	// service token
	cmdService.AddCommand(cmdServiceToken)
	cmdServiceToken.AddCommand(cmdServiceTokenGenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenDelete)

	var rootCmd = &cobra.Command{Use: "ocmgr", Version: version}

	rootCmd.PersistentFlags().StringVarP(&frameworkServer, "framework-server", "s", "http://localhost", "Specifies the framework server")
	rootCmd.PersistentFlags().StringVarP(&authID, "auth-id", "i", "", "The authentication ID to use with the framework server")
	rootCmd.PersistentFlags().StringVarP(&authToken, "auth-token", "t", "", "The authentication token to use with the framework server")

	rootCmd.AddCommand(cmdService)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		host = rest.NewHost(frameworkServer)
		host.Login(authID, authToken)
	}
	rootCmd.Execute()
}
