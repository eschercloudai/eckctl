package cmd

import (
	"eckctl/cmd/create"
	"eckctl/cmd/delete"
	"eckctl/cmd/get"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func Execute() {
	cmd := NewRootCommand()
	commands := []*cobra.Command{
		get.NewGetCommand(),
		create.NewCreateCommand(),
		delete.NewDeleteCommand(),
	}
	cmd.AddCommand(commands...)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	url := ""
	username := ""
	password := ""
	project := ""

	rootCmd := &cobra.Command{
		Use:   "eckctl",
		Short: "A CLI for working with the EscherCloud Kubernetes (ECK) Service",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "Username is: ", username)
		},
	}

	rootCmd.PersistentFlags().StringVar(&url, "url", "https://unikorn.nl1.eschercloud.dev", "URL to Unikorn API")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password")
	rootCmd.PersistentFlags().StringVar(&project, "project", "", "Project ID")
	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".eckctl" (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".eckctl")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
