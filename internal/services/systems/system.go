/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package systems

import "github.com/yakumioto/alkaid/internal/common/storage"

const (
	KSystemInitialized = "system_initialized"
)

var (
	VSystemInitialized = "initialized"
)

type System struct {
	Key       string `json:"key,omitempty" gorm:"primaryKey"`
	Value     string `json:"value,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
}

func newSystem(k, v string) *System {
	return &System{
		Key:   k,
		Value: v,
	}
}

func newSystemByID(id string) *System {
	return &System{
		Key: id,
	}
}

func (s *System) create() error {
	return storage.Create(s)
}

func (s *System) findByID() error {
	return storage.FindByID(s, s.Key)
}
