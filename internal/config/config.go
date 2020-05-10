/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package config

import (
	"github.com/spf13/cobra"
)

var (
	Path           string
	LogLevel       string
	DBPath         string
	Address        string
	Port           int
	FileSystemPath string
)

func InitConfig(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "alkaid",
		Short: "",
		Long:  "",
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&Path, "config", "c", "/var/alkaid/config", "config path.")
	rootCmd.Flags().StringVar(&LogLevel, "logLevel", "INFO", "log level.")
	rootCmd.Flags().StringVar(&DBPath, "dbPath", "/var/alkaid/alkaid.db", "sqlite3 db path.")
	rootCmd.Flags().StringVarP(&Address, "address", "l", "0.0.0.0", "listening address.")
	rootCmd.Flags().IntVarP(&Port, "port", "p", 8080, "listening port.")
	rootCmd.Flags().StringVar(&FileSystemPath, "fileSystemPath", "/var/alkaid", "Project data files.")

	return rootCmd
}
