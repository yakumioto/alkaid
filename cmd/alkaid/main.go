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

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/storage/sqlite3"
	"github.com/yakumioto/alkaid/internal/restful"
	"github.com/yakumioto/alkaid/internal/restful/controllers"
	"github.com/yakumioto/alkaid/internal/restful/middlewares"
)

func main() {
	var (
		db  storage.Storage
		err error
	)

	initConfig()

	log.Initialize(viper.GetString("logging.level"))

	switch viper.GetString("database.use") {
	case sqlite3.Driver:
		db, err = sqlite3.NewDB(viper.GetString("database.sqlite3.path"))
		if err != nil {
			logrus.Panicf("new sqlite3 database error: %v", err)
		}
	}

	storage.Initialize(db)

	service := restful.NewService(
		restful.WithMode(viper.GetString("restful.mode")),
		restful.WithRequestTimeout(viper.GetDuration("restful.request.timeout")),
	)
	service.RegisterMiddlewares(
		new(middlewares.Logger),
		new(middlewares.Recovery),
		new(middlewares.ResolveVersion),
	)
	service.RegisterControllers(
		new(controllers.Health),
	)

	if err := service.Run(viper.GetString("restful.address")); err != nil {
		logrus.Panicf("running restful service error: %v", err)
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
