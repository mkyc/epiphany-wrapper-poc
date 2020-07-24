/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"

	"github.com/spf13/cobra"
)

// environmentsRunCmd represents the run command
var environmentsRunCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		if len(args) == 2 {
			config, err := configuration.GetConfig()
			if err != nil {
				panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
			}
			e, err := environment.Get(config.CurrentEnvironment)
			if err != nil {
				panic(fmt.Sprintf("environemtns details failed: %v\n", err)) //TODO err
			}
			c, err := e.GetComponentByName(args[0])
			if err != nil {
				panic(fmt.Sprintf("environemtns get component failed: %v\n", err)) //TODO err
			}
			err = c.Run(args[1])
			if err != nil {
				panic(fmt.Sprintf("run command failed: %v\n", err)) //TODO err
			}
			fmt.Printf("running command completed!")
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsRunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsRunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsRunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
