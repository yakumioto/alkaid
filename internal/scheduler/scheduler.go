/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package scheduler

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/scheduler/docker"
	dockervm "github.com/yakumioto/alkaid/internal/vm/docker"
)

var (
	logger    *glog.Logger
	scheduler Scheduler
)

func Init() {
	logger = glog.MustGetLogger("scheduler")

	docker.Init(logger)
}

// Scheduler Scheduling layer, can be based on docker, docker swarm, kubernetes
type Scheduler interface {
	Network
	PeerNode
	OrdererNode
	CANode
}

func NewScheduler(typ string) (Scheduler, error) {
	if scheduler != nil {
		return scheduler, nil
	}

	switch typ {
	case types.DockerNetworkType:
		cli, err := dockervm.NewController()
		if err != nil {
			logger.Errof("New vm docker error: %s", err)
			return nil, err
		}

		scheduler = docker.NewDocker(cli)

	default:
		return nil, errors.New(fmt.Sprintf("no scheduler of this type: %s", typ))
	}

	return scheduler, nil
}

// network Operation hyperledger
type Network interface {
	CreateNetwork(network *types.Network) error
	DeleteNetwork() error
}

// PeerNode Operation hyperledger fabric peer node
type PeerNode interface {
	CreatePeer() error
	RestartPeer() error
	StopPeer() error
	DeletePeer() error
}

// OrdererNode Operation hyperledger fabric orderer node
type OrdererNode interface {
	CreateOrderer() error
	RestartOrderer() error
	StopOrderer() error
	DeleteOrderer() error
}

// CANode Operation hyperledger fabric certificate authority node
type CANode interface {
	CreateCA() error
	RestartCA() error
	StopCA() error
	DeleteCA() error
}
