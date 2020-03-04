/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	network2 "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
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

type CreateRequest struct {
	ContainerName  string
	ImageName      string
	ImageTag       string
	Environment    []string
	NetworkMode    string
	NetworkAliases []string
	Mounts         map[string]string
	Files          map[string][]byte
}

func (c *Controller) Create(cr *CreateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mounts := make([]mount.Mount, len(cr.Mounts))
	for source, target := range cr.Mounts {
		if err := c.createVolumeWithDockerMode(source); err != nil {
			return err
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: source,
			Target: target,
		})
	}

	config := &container.Config{
		Image: fmt.Sprintf("%s:%s", cr.ImageName, cr.ImageTag),
		Env:   cr.Environment,
	}

	host := &container.HostConfig{
		NetworkMode: container.NetworkMode(cr.NetworkMode),
		Mounts:      mounts,
	}

	network := &network2.NetworkingConfig{
		EndpointsConfig: map[string]*network2.EndpointSettings{
			cr.NetworkMode: {
				Aliases: cr.NetworkAliases,
			},
		},
	}

	res, err := c.cli.ContainerCreate(ctx, config, host, network, cr.ContainerName)
	if err != nil {
		return errors.Wrap(err, "create container failed")
	}

	for path, content := range cr.Files {
		if err := c.copyToContainer(res.ID, path, bytes.NewReader(content)); err != nil {
			return errors.Wrap(err, "cp to container failed")
		}
	}

	return nil
}

func (c *Controller) Start() error {
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

func (c *Controller) CreateNetworkWithDockerMode(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.cli.NetworkCreate(ctx, name, dockertypes.NetworkCreate{
		CheckDuplicate: true,
	})
	if err != nil {
		return "", err
	}

	return result.ID, err
}

func (c *Controller) createVolumeWithDockerMode(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := c.cli.VolumeCreate(ctx, volume.VolumesCreateBody{
		Name:   name,
		Driver: "local",
	}); err != nil {
		return errors.Wrap(err, "volume create failed")
	}

	return nil
}

func (c *Controller) copyToContainer(id, path string, content io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	return c.cli.CopyToContainer(ctx, id, path, content, dockertypes.CopyToContainerOptions{})
}
