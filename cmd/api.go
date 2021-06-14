package cmd

import (
	"Panda/api"
	"Panda/common/log"
	"Panda/conf"
	"context"
	"github.com/spf13/cobra"
)

var (
	frontendCmd = &cobra.Command{
		Use: "api",
	}

	frontendStartCmd = &cobra.Command{
		Use: "start",
		Run: apiStart,
	}

	switchCron int
)

func init() {
	rootCmd.AddCommand(frontendCmd)
	frontendCmd.AddCommand(frontendStartCmd)

	frontendStartCmd.PersistentFlags().IntVarP(&switchCron, "switchCron", "c", 0, "开启定时任务")
}

func apiStart(cmd *cobra.Command, args []string) {
	fe := api.New(conf.Config)
	if switchCron == 1 {
		fe.StartCron()
	}
	if err := fe.Start(context.Background()); err != nil {
		log.Error(err.Error())
	}
}
