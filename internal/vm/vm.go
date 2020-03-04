/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package vm

import (
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/vm/docker"
)

var (
	logger *glog.Logger
	vm     VirtualMachine
)

type VirtualMachine interface {
	Create(node *types.Node) error
	Restart() error
	Stop() error
	Delete() error
}

func Init() {
	logger = glog.MustGetLogger("vm")

	if vm == nil {
		c, err := docker.NewController()
		if err != nil {
			logger.Panicf("New docker controller error: %v", err)
		}

		vm = c
	}

	logger.Debuf("Virtual machine initialized successfully")
}
