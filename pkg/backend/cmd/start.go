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
	Short: "start controller",
	Long:  ``,
}

var startOnceCmd = &cobra.Command{
	Use:   "once",
	Short: "start for once",
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
			log.Fatal("getting config", err)
		}
		if config.GetTemplate(templatePath) != nil {
			log.Fatal("getting template", err)
		}

		err = get.GettingDeviceConfig()
		if err != nil {
			notify.NotifyErrorToSlack(err)
		}

		log.Println("end")
	},
}

var startCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "start for cron",
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
			log.Fatal("getting config", err)
		}
		if config.GetTemplate(templatePath) != nil {
			log.Fatal("getting template", err)
		}

		err = get.CronExec()
		if err != nil {
			notify.NotifyErrorToSlack(err)
		}

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startOnceCmd)
	startCmd.AddCommand(startCronCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
	startCmd.PersistentFlags().StringP("template", "t", "", "config path")
}
