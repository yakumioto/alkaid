/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package docker

import (
	"context"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Controller struct {
	cli *client.Client
}

func NewController() (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Controller{cli: cli}, err
}

func (c *Controller) Create() error {
	return nil
}

func (c *Controller) Restart() error {
	return nil
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) Delete() error {
	return nil
}

func (c *Controller) CreateNetworkWithDockerMode(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.cli.NetworkCreate(ctx, name, dockertypes.NetworkCreate{
		CheckDuplicate: true,
	})

	return err
}
