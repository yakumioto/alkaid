/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "v0.3.0"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Alkaid",
		Long:  `All software has versions. This is Alkaid's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
