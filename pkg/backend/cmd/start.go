package cmd

import (
	"github.com/doornoc/config-collector/pkg/api/core/get"
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/notify"
	"github.com/spf13/cobra"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start controller server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		templatePath, err := cmd.Flags().GetString("template")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}
		if config.GetTemplate(templatePath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		err = get.GettingDeviceConfig()
		if err != nil {
			notify.NotifyErrorToSlack(err)
		}

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
	startCmd.PersistentFlags().StringP("template", "t", "", "config path")
}
