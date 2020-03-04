/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package docker

type CANode struct{}

func (c *CANode) CreateCA() error {
	return nil
}
func (c *CANode) RestartCA() error {
	return nil
}
func (c *CANode) StopCA() error {
	return nil
}
func (c *CANode) DeleteCA() error {
	return nil
}
