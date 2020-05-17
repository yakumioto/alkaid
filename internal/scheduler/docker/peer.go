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
	"fmt"

	"github.com/yakumioto/alkaid/internal/vm"
	dockervm "github.com/yakumioto/alkaid/internal/vm/docker"
)

type peerNode struct {
	cli *dockervm.Controller
}

func newPeer(cli *dockervm.Controller) *peerNode {
	return &peerNode{cli: cli}
}

func (p *peerNode) CreatePeer(peer *vm.CreateRequest, crs ...*vm.CreateRequest) error {
	ids := make([]string, 0)

	// other nodes eg: couchdb
	for _, cr := range crs {
		if err := p.cli.Create(cr); err != nil {
			logger.Errof("Create crs error: %s", err)
			return err
		}
		ids = append(ids, cr.ContainerName)
	}

	if peer.Environment == nil {
		peer.Environment = make([]string, 0)
	}

	peer.Environment = append(peer.Environment,
		"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
		"FABRIC_LOGGING_SPEC=INFO",
		"CORE_PEER_TLS_ENABLED=true",
		"CORE_PEER_GOSSIP_USELEADERELECTION=true",
		"CORE_PEER_GOSSIP_ORGLEADER=false",
		"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
		"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
		"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
		fmt.Sprintf("CORE_PEER_ID=%s", peer.ContainerName),
		fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peer.ContainerName),
		"CORE_PEER_LISTENADDRESS=0.0.0.0:7051",
		fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=%s:7052", peer.ContainerName),
		"CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052",
	)

	if peer.BindMounts == nil {
		peer.BindMounts = make(map[string]string)
	}
	peer.BindMounts["/var/run/"] = "/var/run/"

	if err := p.cli.Create(peer); err != nil {
		logger.Errof("Create peer error: %s", err)
		return err
	}

	ids = append(ids, peer.ContainerName)

	// start nodes
	for _, id := range ids {
		if err := p.cli.Start(id); err != nil {
			logger.Errof("Start nodes error: %s", err)
			return err
		}
	}

	return nil
}

func (p *peerNode) RestartPeer(ids ...string) error {
	for _, id := range ids {
		if err := p.cli.Restart(id); err != nil {
			logger.Errof("Restart peer error: %s", err)
			return err
		}
	}
	return nil
}

func (p *peerNode) StopPeer(ids ...string) error {
	for _, id := range ids {
		if err := p.cli.Stop(id); err != nil {
			logger.Errof("Stop peer error: %s", err)
			return err
		}
	}
	return nil
}

func (p *peerNode) DeletePeer(ids ...string) error {
	for _, id := range ids {
		if err := p.cli.Delete(id); err != nil {
			logger.Errof("Delete peer error: %s", err)
			return err
		}
	}
	return nil
}
