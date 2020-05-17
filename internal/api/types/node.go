/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package types

import (
	"fmt"
)

const (
	RunningNodeStatus = "running"
	StopedNodeStatus  = "stoped"

	PeerNodeType    = "peer"
	OrdererNodeType = "orderer"
)

// Image docker image
type Image struct {
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
}

func (i *Image) String() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

// Node is a docker container or k8s pod
type Node struct {
	ID                 int64    `json:"-"`
	NodeID             string   `json:"node_id,omitempty"`
	OrganizationID     string   `json:"organization_id,omitempty" binding:"required"`
	UserID             string   `json:"user_id,omitempty" binding:"required"`
	NetworkID          string   `json:"network_id,omitempty" binding:"required"`
	Type               string   `json:"type,omitempty"`
	CouchDB            bool     `json:"couch_db,omitempty"`
	Status             string   `json:"status,omitempty"`
	DockerContainerIDs []string `json:"-"`
	CreatedAt          int64    `json:"created_at,omitempty"`
	UpdatedAt          int64    `json:"updated_at,omitempty"`
}

func NewNode() *Node {
	return &Node{}
}
