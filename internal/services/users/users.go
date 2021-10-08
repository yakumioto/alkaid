/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package users

import (
	"time"

	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/util"
)

const (
	Namespace = "User"
)

type User struct {
	ID                     string `json:"id,omitempty"`
	ResourceID             string `json:"resourceId,omitempty"`
	Name                   string `json:"name,omitempty"`
	Email                  string `json:"email,omitempty"`
	Password               string `json:"password,omitempty"`
	Role                   string `json:"role,omitempty"`
	ProtectedSigPrivateKey string `json:"protectedSigPrivateKey,omitempty"`
	ProtectedTLSPrivateKey string `json:"protectedTlsPrivateKey,omitempty"`
	Status                 string `json:"status,omitempty"`
	CreatedAt              int64  `json:"createdAt,omitempty"`
	UpdatedAt              int64  `json:"updatedAt,omitempty"`
}

func newUser(req *CreateRequest) *User {
	return &User{
		ID:       req.ID,
		Email:    req.Email,
		Password: util.HashPassword(req.Password, req.Email, 10000),
	}
}

func (u *User) create() error {
	u.ResourceID = util.GenResourceID(Namespace)
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	return storage.Create(u)
}
