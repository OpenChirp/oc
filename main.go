package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/openchirp/framework/rest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version string = "1.0"
)

const (
	columnPadding = 3
)

func main() {
	var host rest.Host

	viper.SetConfigName("occonfig")         // name of config file (without extension)
	viper.AddConfigPath("/etc/oc/")         // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/oc") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.oc")        // call multiple times to add many search paths
	viper.AddConfigPath(".")                // optionally look for config in the working directory

	var cmdUser = &cobra.Command{
		Use:   "user",
		Short: "Manage the user account",
	}

	var cmdUserInfo = &cobra.Command{
		Use:   "info",
		Short: "Fetch user info",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

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
		},
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
		Run: func(cmd *cobra.Command, args []string) {
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
		},
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
				fmt.Fprintln(os.Stderr, "Failed to create service:", err)
				os.Exit(1)
			}
			fmt.Println(s.ID)
		},
	}

	var cmdServiceRm = &cobra.Command{
		Use:   "rm <service_id>",
		Short: "Remove a new service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			err := host.ServiceDelete(serviceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to delete service:", err)
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
				fmt.Fprintln(os.Stderr, "Failed to generate token:", err)
				os.Exit(1)
			}
			if envTrue, _ := cmd.Flags().GetBool("env"); envTrue {
				fmt.Printf("FRAMEWORK_SERVER=\"%s\"\n", viper.GetString("framework-server"))
				fmt.Printf("SERVICE_ID=\"%s\"\n", serviceID)
				fmt.Printf("SERVICE_TOKEN=\"%s\"\n", token)
				// missing MQTT server
			} else {
				fmt.Println(token)
			}
		},
	}
	cmdServiceTokenGenerate.Flags().Bool("env", false, "Print out all service environment variables to setup a service")

	var cmdServiceTokenRegenerate = &cobra.Command{
		Use:   "regenerate <service_id>",
		Short: "Regenerate a security token for the service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			token, err := host.ServiceTokenRegenerate(serviceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to regenerate token:", err)
				os.Exit(1)
			}
			if envTrue, _ := cmd.Flags().GetBool("env"); envTrue {
				fmt.Printf("FRAMEWORK_SERVER=\"%s\"\n", viper.GetString("framework-server"))
				fmt.Printf("SERVICE_ID=\"%s\"\n", serviceID)
				fmt.Printf("SERVICE_TOKEN=\"%s\"\n", token)
				// missing MQTT server
			} else {
				fmt.Println(token)
			}
		},
	}
	cmdServiceTokenRegenerate.Flags().Bool("env", false, "Print out all service environment variables to setup a service")

	var cmdServiceTokenRm = &cobra.Command{
		Use:   "rm <service_id>",
		Short: "Remove the security token for the service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceID := args[0]

			err := host.ServiceTokenDelete(serviceID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to delete token:", err)
				os.Exit(1)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "oc", Version: version}

	// oc
	rootCmd.AddCommand(cmdService)
	rootCmd.AddCommand(cmdUser)
	// oc service
	cmdService.AddCommand(cmdServiceLs)
	cmdService.AddCommand(cmdServiceCreate)
	cmdService.AddCommand(cmdServiceRm)
	// oc service token
	cmdService.AddCommand(cmdServiceToken)
	cmdServiceToken.AddCommand(cmdServiceTokenGenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRegenerate)
	cmdServiceToken.AddCommand(cmdServiceTokenRm)
	// oc user
	cmdUser.AddCommand(cmdUserInfo)

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
			fmt.Fprintln(os.Stderr, "Failed to parse config file:", err)
			os.Exit(1)
		}
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		host = rest.NewHost(viper.GetString("framework-server"))
		host.Login(viper.GetString("auth-id"), viper.GetString("auth-token"))
	}
	rootCmd.Execute()
}
