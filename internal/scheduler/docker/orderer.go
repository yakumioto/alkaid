/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package docker

import (
	"github.com/yakumioto/alkaid/internal/vm"
	dockervm "github.com/yakumioto/alkaid/internal/vm/docker"
)

type ordererNode struct {
	cli *dockervm.Controller
}

func (o *ordererNode) CreateOrderer(crs ...*vm.CreateRequest) error {
	for _, cr := range crs {
		if err := o.cli.Create(cr); err != nil {
			logger.Errof("Create orderer error: %s", err)
			return err
		}
	}
	return nil
}

func (o *ordererNode) RestartOrderer(ids ...string) error {
	for _, id := range ids {
		if err := o.cli.Restart(id); err != nil {
			logger.Errof("Restart orderer error: %s", err)
			return err
		}
	}
	return nil
}

func (o *ordererNode) StopOrderer(ids ...string) error {
	for _, id := range ids {
		if err := o.cli.Restart(id); err != nil {
			logger.Errof("Stop orderer error: %s", err)
			return err
		}
	}
	return nil
}

func (o *ordererNode) DeleteOrderer(ids ...string) error {
	for _, id := range ids {
		if err := o.cli.Restart(id); err != nil {
			logger.Errof("Delete orderer error: %s", err)
			return err
		}
	}
	return nil
}
