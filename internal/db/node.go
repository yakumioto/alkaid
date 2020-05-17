/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package db

import (
	"fmt"
	"time"

	"github.com/yakumioto/alkaid/internal/api/types"
)

var (
	ErrNodeExist    = new(NodeExistError)
	ErrNodeNotExist = new(NodeNotExistError)
)

type nodeError struct {
	NodeID string
}

type NodeExistError struct {
	nodeError
}

func (e *NodeExistError) Error() string {
	return fmt.Sprintf("node already exists [node_id: %s]", e.NodeID)
}

type NodeNotExistError struct {
	nodeError
}

func (e *NodeNotExistError) Error() string {
	return fmt.Sprintf("node not exists [node_id: %s]", e.NodeID)
}

type Node struct {
	ID                 int64    `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	NodeID             string   `xorm:"'node_id' UNIQUE INDEX NOT NULL"`
	OrganizationID     string   `xorm:"'organization_id'"`
	UserID             string   `xorm:"'user_id'"`
	NetworkID          string   `xorm:"'network_id'"`
	Type               string   `xorm:"'type'"`
	CouchDB            bool     `xorm:"'couch_db'"`
	Status             string   `xorm:"'status'"`
	DockerContainerIDs []string `xorm:"'docker_container_ids'"`
	CreatedAt          int64    `xorm:"'created_at'"`
	UpdatedAt          int64    `xorm:"'updated_at'"`
}

func (*Node) TableName() string {
	return "node"
}

func (n *Node) BeforeInsert() {
	n.CreatedAt = time.Now().Unix()
	n.UpdatedAt = n.CreatedAt
}

func isNodeExist(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	return x.Get(&Node{NodeID: id})
}

func CreateNode(node *Node) error {
	exist, err := isNodeExist(node.NodeID)
	if err != nil {
		return err
	}

	if exist {
		return &NodeExistError{nodeError{NodeID: node.NodeID}}
	}

	_, err = x.Insert(node)
	if err != nil {
		return err
	}

	return nil
}

func QueryNodeByNodeID(id string) (*types.Node, error) {
	node := &Node{NodeID: id}

	has, err := x.Get(node)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &NodeNotExistError{nodeError{NodeID: id}}
	}

	return (*types.Node)(node), nil
}
