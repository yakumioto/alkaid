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

	"github.com/yakumioto/alkaid/internal/vm"
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

func (c *Controller) Create(cr *vm.CreateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if cr.ImageName == "" || cr.ImageTag == "" {
		return errors.New("image name and tag cannot be empty")
	}

	has, err := c.hasImage(cr.ImageName, cr.ImageTag)
	if err != nil {
		return errors.Wrap(err, "check image error")
	}

	if !has {
		if err = c.pullImage(cr.ImageName, cr.ImageTag); err != nil {
			return errors.Wrap(err, "pull image error")
		}
	}

	// FIXME: check networkmode exist
	has, err = c.hasNetwork(cr.NetworkMode)
	if err != nil {
		return errors.Wrap(err, "check network error")
	}

	if !has {
		return errors.New("network not exist")
	}

	mounts := make([]mount.Mount, 0)
	for source, target := range cr.VolumeMounts {
		if err1 := c.createVolumeWithDockerMode(source); err1 != nil {
			return err1
		}

		mounts = append(mounts, mount.Mount{
			Type:          mount.TypeVolume,
			Source:        source,
			Target:        target,
			VolumeOptions: &mount.VolumeOptions{NoCopy: true},
		})
	}
	for source, target := range cr.BindMounts {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		})
	}

	containerConfig := &container.Config{
		Image:      fmt.Sprintf("%s:%s", cr.ImageName, cr.ImageTag),
		Env:        cr.Environment,
		WorkingDir: cr.WorkingDir,
		Cmd:        cr.Command,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(cr.NetworkMode),
		Mounts:      mounts,
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			cr.NetworkMode: {
				Aliases: cr.NetworkAliases,
			},
		},
	}

	res, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, cr.ContainerName)
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

	return c.cli.CopyToContainer(ctx, id, path, content, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
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
