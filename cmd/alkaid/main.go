/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package main

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/cmd/alkaid/start"
	"github.com/yakumioto/alkaid/cmd/alkaid/version"
)

func main() {
	rand.Seed(time.Now().Unix())

	manCmd := &cobra.Command{
		Use:  "alkaid",
		Long: `Alkaid is a BaaS(Blockchan as a Service) service based on Hyperledger Fabric.`,
	}

	manCmd.AddCommand(start.Cmd())
	manCmd.AddCommand(version.Cmd())

	if err := manCmd.Execute(); err != nil {
		glog.Fatalln(err)
	}
}
