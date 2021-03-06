package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/spf13/cobra"
)

// envListCmd represents the list command
var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "TODO",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments list pre run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("list called")
		environments, err := environment.GetAll()
		if err != nil {
			logger.Fatal().Err(err).Msg("environments get all failed")
		}
		for _, e := range environments {
			if e.Uuid.String() == config.CurrentEnvironment.String() {
				fmt.Printf("* (%s) | %s\n", e.Uuid.String(), e.Name)
			} else {
				fmt.Printf("  (%s) | %s\n", e.Uuid.String(), e.Name)
			}
		}
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
}
