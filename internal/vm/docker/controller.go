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
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
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
	WorkintDir     string
	Command        []string
}

func (c *Controller) Create(createRequest *CreateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if createRequest.ImageName == "" || createRequest.ImageTag == "" {
		return errors.New("image name and tag cannot be empty")
	}

	has, err := c.hasImage(createRequest.ImageName, createRequest.ImageTag)
	if err != nil {
		return errors.Wrap(err, "check image error")
	}

	if !has {
		if err = c.pullImage(createRequest.ImageName, createRequest.ImageTag); err != nil {
			return errors.Wrap(err, "pull image error")
		}
	}

	// FIXME: check networkmode exist
	has, err = c.hasNetwork(createRequest.NetworkMode)
	if err != nil {
		return errors.Wrap(err, "check network error")
	}

	if !has {
		return errors.New("network not exist")
	}

	mounts := make([]mount.Mount, 0)
	for source, target := range createRequest.Mounts {
		if err1 := c.createVolumeWithDockerMode(source); err1 != nil {
			return err1
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: source,
			Target: target,
		})
	}

	containerConfig := &container.Config{
		Image:      fmt.Sprintf("%s:%s", createRequest.ImageName, createRequest.ImageTag),
		Env:        createRequest.Environment,
		WorkingDir: createRequest.WorkintDir,
		Cmd:        createRequest.Command,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(createRequest.NetworkMode),
		Mounts:      mounts,
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			createRequest.NetworkMode: {
				Aliases: createRequest.NetworkAliases,
			},
		},
	}

	res, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, createRequest.ContainerName)
	if err != nil {
		return errors.Wrap(err, "create container failed")
	}

	for path, content := range createRequest.Files {
		if err := c.copyToContainer(res.ID, path, bytes.NewReader(content)); err != nil {
			return errors.Wrap(err, "cp to container failed")
		}
	}

	return nil
}

func (c *Controller) Start(containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	return c.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (c *Controller) Restart(containerID string) error {
	timeout := 3 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.cli.ContainerRestart(ctx, containerID, &timeout)
}

func (c *Controller) Stop(containerID string) error {
	timeout := 1 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.cli.ContainerStop(ctx, containerID, nil)
}

func (c *Controller) Delete(containerID string) error {
	timeout := 1 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{RemoveVolumes: true})
}

func (c *Controller) CreateNetworkWithDockerMode(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.cli.NetworkCreate(ctx, name, types.NetworkCreate{
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

	return c.cli.CopyToContainer(ctx, id, path, content, types.CopyToContainerOptions{})
}

func (c *Controller) hasImage(name, tag string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := filters.NewArgs()
	args.Add("reference", fmt.Sprintf("%s:%s", name, tag))

	list, err := c.cli.ImageList(ctx, types.ImageListOptions{
		Filters: args,
	})
	if err != nil {
		return false, err
	}

	if len(list) == 0 {
		return false, nil
	}

	return true, nil
}

func (c *Controller) pullImage(name, tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// XXX: Replaceable database?
	pullClose, err := c.cli.ImagePull(ctx,
		fmt.Sprintf("docker.io/library/%s:%s", name, tag), types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, pullClose)

	return pullClose.Close()
}

func (c *Controller) hasNetwork(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := filters.NewArgs()
	args.Add("name", name)
	args.Add("driver", "bridge")

	res, err := c.cli.NetworkList(ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return false, err
	}

	if len(res) > 0 {
		return true, nil
	}

	return false, nil
}
