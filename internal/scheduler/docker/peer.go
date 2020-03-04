/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package docker

type PeerNode struct {
}

func (p *PeerNode) CreatePeer() error {
	return nil
}
func (p *PeerNode) RestartPeer() error {
	return nil
}
func (p *PeerNode) StopPeer() error {
	return nil
}
func (p *PeerNode) DeletePeer() error {
	return nil
}
