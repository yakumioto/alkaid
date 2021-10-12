/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package main

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/storage/sqlite3"
	"github.com/yakumioto/alkaid/internal/restful"
	"github.com/yakumioto/alkaid/internal/restful/controllers"
)

func main() {
	var (
		db  storage.Storage
		err error
	)

	initConfig()

	switch viper.GetString("database.use") {
	case sqlite3.Driver:
		db, err = sqlite3.NewDB(viper.GetString("database.sqlite3.path"))
		if err != nil {
			logrus.Panicf("new sqlite3 database error: %v", err)
		}
	}

	storage.Initialization(db)

	service := restful.NewService(
		restful.WithMode(restful.DevelopMode),
		restful.WithRequestTimeout(10*time.Second))
	service.RegisterControllers(
		new(controllers.Health),
	)

	if err := service.Run(viper.GetString("restful.address")); err != nil {

	}
}

func initConfig() {
	viper.SetConfigName("alkaid")
	viper.AddConfigPath("/etc/alkaid")
	viper.AddConfigPath("$HOME/.alkaid")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config, ok := os.LookupEnv("CONFIG_PATH")
	if ok {
		viper.SetConfigFile(config)
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("Fatal error config file: %v", err)
	}
}
