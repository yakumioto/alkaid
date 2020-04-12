/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package docker

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/stretchr/testify/assert"

	"github.com/yakumioto/alkaid/internal/utils/targz"
)

func testInit(t *testing.T) *Controller {
	c, err := NewController()
	assert.NoError(t, err, "new docker controller error")

	return c
}

func testDeleteController(t *testing.T, containerID, imageID string) {
	c := testInit(t)

	err := c.Delete(containerID)
	assert.NoError(t, err)

	_, err = c.cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{})
	assert.NoError(t, err)
}

func equalVolumeName(t *testing.T, name string) {
	c := testInit(t)

	filter := filters.NewArgs()
	filter.Add("name", name)

	res, err := c.cli.VolumeList(context.Background(), filter)
	assert.NoError(t, err)

	assert.Equal(t, res.Volumes[0].Name, name)
}

func TestController(t *testing.T) {
	c := testInit(t)

	createRequest := new(CreateRequest)

	// tset empty image name and image tag
	err := c.Create(createRequest)
	assert.EqualError(t, err, "image name and tag cannot be empty")

	files := make([]*targz.File, 1)
	files[0] = &targz.File{
		Name: "test",
		Mode: 0664,
		Body: []byte("hello world\n"),
	}

	filesBytes, err := targz.Generate(files)
	assert.NoError(t, err)

	createRequest.ContainerName = "test"
	createRequest.ImageName = "alpine"
	createRequest.ImageTag = "3.7"
	createRequest.Environment = []string{"TEST_FILE=/home/mioto/test"}
	createRequest.WorkingDir = "/home/mioto"
	createRequest.NetworkMode = "bridge"
	createRequest.Mounts = map[string]string{"test.work": "/home/mioto"}
	createRequest.Files = map[string][]byte{"/home/mioto": filesBytes}
	createRequest.Command = []string{"sleep", "1000"}

	// test create container
	err = c.Create(createRequest)
	assert.NoError(t, err)

	inspect, err := c.cli.ContainerInspect(context.Background(), createRequest.ContainerName)
	assert.NoError(t, err)
	assert.Equal(t, inspect.Name, "/"+createRequest.ContainerName)
	assert.Equal(t, inspect.Config.Image, fmt.Sprintf("%s:%s", createRequest.ImageName, createRequest.ImageTag))
	assert.Equal(t, inspect.Config.Env[0], createRequest.Environment[0])
	assert.Equal(t, inspect.Config.WorkingDir, createRequest.WorkingDir)
	assert.Equal(t, inspect.HostConfig.NetworkMode, container.NetworkMode(createRequest.NetworkMode))
	assert.Equal(t, inspect.HostConfig.Mounts[0].Source, "test.work")
	assert.Equal(t, inspect.HostConfig.Mounts[0].Target, "/home/mioto")
	equalVolumeName(t, "test.work")

	// test start container
	err = c.Start(createRequest.ContainerName)
	assert.NoError(t, err)

	inspect, err = c.cli.ContainerInspect(context.Background(), createRequest.ContainerName)
	assert.NoError(t, err)
	assert.Equal(t, inspect.State.Running, true)

	// test stop container
	err = c.Stop(createRequest.ContainerName)
	assert.NoError(t, err)

	inspect, err = c.cli.ContainerInspect(context.Background(), createRequest.ContainerName)
	assert.NoError(t, err)
	assert.Equal(t, inspect.State.Running, false)

	// test restart container
	err = c.Restart(createRequest.ContainerName)
	assert.NoError(t, err)

	inspect, err = c.cli.ContainerInspect(context.Background(), createRequest.ContainerName)
	assert.NoError(t, err)
	assert.Equal(t, inspect.State.Running, true)

	err = c.Stop(createRequest.ContainerName)
	assert.NoError(t, err)

	testDeleteController(t, createRequest.ContainerName, fmt.Sprintf("%s:%s", createRequest.ImageName, createRequest.ImageTag))
}
