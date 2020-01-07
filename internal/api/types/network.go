/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package types

type Network struct {
	ID          int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	NetworkID   string `xorm:"'network_id' UNIQUE INDEX NOT NULL"`
	Name        string `xorm:"'name'"`
	Description string `xorm:"'description'"`
	CreatedAt   int64  `xorm:"'created_at'"`
	UpdatedAt   int64  `xorm:"'updated_at'"`
}
