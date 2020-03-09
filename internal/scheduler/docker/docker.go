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
	"github.com/yakumioto/glog"

	dockervm "github.com/yakumioto/alkaid/internal/vm/docker"
)

var (
	logger *glog.Logger
)

func Init(log *glog.Logger) {
	logger = log.MustGetLogger("docker")
}

type Docker struct {
	cli *dockervm.Controller
	*network
	*OrdererNode
	*PeerNode
	*CANode
}

func NewDocker(cli *dockervm.Controller) *Docker {
	d := &Docker{
		cli: cli,
	}

	d.network = newNetwork(cli)

	return d
}
