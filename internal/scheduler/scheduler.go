/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package scheduler

// Scheduler Scheduling layer, can be based on docker, docker swarm, kubernetes
type Scheduler interface {
	PeerNode
	OrdererNode
	CANode
}

// PeerNode Operation hyperledger peer node
type PeerNode interface {
	CreatePeer() error
	RestartPeer() error
	StopPeer() error
	DeletePeer() error
}

// OrdererNode Operation hyperledger orderer node
type OrdererNode interface {
	CreateOrderer() error
	RestartOrderer() error
	StopOrderer() error
	DeleteOrderer() error
}

// CANode Operation hyperledger Certificate Authority node
type CANode interface {
	CreateCA() error
	RestartCA() error
	StopCA() error
	DeleteCA() error
}
