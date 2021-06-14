package cmd

import (
	"Panda/common/log"
	"Panda/conf"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	env     string

	rootCmd = &cobra.Command{
		Use:   "Panda",
		Short: "基于iris MVC的框架整理",
	}
)


func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./conf/config.toml)")
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "prod", "env setting")

}

func initConfig() {
	if _, err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Config)
}


func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}