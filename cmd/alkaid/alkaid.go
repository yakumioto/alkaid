/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/handler"
	"github.com/yakumioto/alkaid/internal/api/routers"
	"github.com/yakumioto/alkaid/internal/config"
	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/scheduler"
)

func main() {
	rand.Seed(time.Now().Unix())

	cmd := config.InitConfig(run)
	if err := cmd.Execute(); err != nil {
		glog.Fatalln(err)
	}
}

func run(_ *cobra.Command, _ []string) {
	var debug bool

	gin.SetMode(gin.ReleaseMode)

	switch config.LogLevel {
	case "INFO":
		glog.SetLevel(glog.LevelInfo)
	case "WARN":
		glog.SetLevel(glog.LevelWarning)
	case "ERRO":
		glog.SetLevel(glog.LevelError)
	case "DEBU":
		debug = true
		glog.SetLevel(glog.LevelDebug)
	default:
		glog.SetLevel(glog.LevelInfo)
	}

	if debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	routers.Init(r)
	handler.Init()
	scheduler.Init()

	if err := db.Init(config.DBPath, "cache=shared&mode=rwc"); err != nil {
		glog.Fatalf("db initialization failed: %s", err)
	}

	if err := r.Run(config.Address + ":" + strconv.Itoa(config.Port)); err != nil {
		glog.Fatalf("service startup failed: %s", err)
	}
}
