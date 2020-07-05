/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package start

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/models"
	"github.com/yakumioto/alkaid/internal/services"
	"github.com/yakumioto/alkaid/internal/services/organization"
	"github.com/yakumioto/alkaid/internal/services/user"
)

var (
	logger = glog.MustGetLogger("cmd.start")
)

func Cmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Start RESTful service",
		Run:   serve,
	}

	command.PersistentFlags().StringP("configPath", "c", "", "config file (default is $HOME/.alkaid/config.yaml)")
	command.PersistentFlags().StringP("logLevel", "l", "INFO", "log level")
	_ = viper.BindPFlag("configPath", command.PersistentFlags().Lookup("configPath"))
	_ = viper.BindPFlag("logLevel", command.PersistentFlags().Lookup("logLevel"))
	return command
}

func serve(_ *cobra.Command, _ []string) {
	initConfig()
	setGinMode(setLoggerLevel(viper.GetString("core.loglevel")))

	r := gin.Default()
	services.NewServices(r, organization.NewService(), user.NewService())

	if err := models.InitModels(); err != nil {
		logger.Fatalf("Init Models error: %v", err)
	}

	if err := r.Run(viper.GetString("core.address") + ":" + viper.GetString("core.port")); err != nil {
		logger.Fatalf("Start serivce error: %v", err)
	}
}

func initConfig() {
	configPath := viper.GetString("configPath")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			glog.Fatalf("Get home dir error: %v", err)
		}

		configPath = home + "/.alkaid"

		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("Alkaid")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Debuf("Using config file: %v", viper.ConfigFileUsed())
	}
}

func setLoggerLevel(level string) bool {
	switch level {
	case "WARN":
		glog.SetLevel(glog.LevelWarning)
	case "ERRO":
		glog.SetLevel(glog.LevelError)
	case "DEBU":
		glog.SetLevel(glog.LevelDebug)
		return true
	default:
		glog.SetLevel(glog.LevelInfo)
	}
	return false
}

func setGinMode(debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(gin.ReleaseMode)
}
