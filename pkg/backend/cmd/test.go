package cmd

import (
	"github.com/doornoc/config-collector/pkg/api/core/get"
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/notify"
	"github.com/spf13/cobra"
	"log"
)

// testCmd represents the start command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test controller without git function",
	Long:  ``,
}

var testOnceCmd = &cobra.Command{
	Use:   "once",
	Short: "test for once without git function",
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

		err = get.GettingDeviceConfig(false)
		if err != nil {
			notify.NotifyErrorToSlack(err)
		}

		log.Println("end")
	},
}

var testCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "test for cron without git function",
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

		err = get.CronExec(false)
		if err != nil {
			notify.NotifyErrorToSlack(err)
		}

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testOnceCmd)
	testCmd.AddCommand(testCronCmd)
	testCmd.PersistentFlags().StringP("config", "c", "", "config path")
	testCmd.PersistentFlags().StringP("template", "t", "", "config path")
}
