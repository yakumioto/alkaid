/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package users

import (
	"errors"
	"strconv"
	"time"

	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/util"
)

const (
	Namespace = "User"

	RoleRoot Role = iota
	RoleOrganization
	RoleNetworkAdmin
	RoleUser
)

var TimeNowFunc = func() int64 {
	return time.Now().Unix()
}

var roleNames = []string{
	"root",
	"organizationAdmin",
	"networkAdmin",
	"user",
}

var roleMap = map[string]Role{
	"root":              RoleRoot,
	"organizationAdmin": RoleOrganization,
	"networkAdmin":      RoleNetworkAdmin,
	"user":              RoleUser,
}

func LookRole(str string) Role {
	role, ok := roleMap[str]
	if !ok {
		return -1
	}

	return role
}

type Role int

func (r Role) String() string {
	if RoleRoot <= r && r <= RoleUser {
		return roleNames[r]
	}

	return "%!Role(" + strconv.Itoa(int(r)) + ")"
}

func (r Role) Less(role Role) bool {
	if r == -1 {
		return true
	}

	return role < r
}

type User struct {
	ID                     string `json:"id,omitempty"`
	ResourceID             string `json:"resourceId,omitempty"`
	Name                   string `json:"name,omitempty"`
	Email                  string `json:"email,omitempty"`
	Password               string `json:"-"`
	Role                   string `json:"role,omitempty"`
	ProtectedSigPrivateKey string `json:"protectedSigPrivateKey,omitempty"`
	ProtectedTLSPrivateKey string `json:"protectedTlsPrivateKey,omitempty"`
	Status                 string `json:"status,omitempty"`
	CreatedAt              int64  `json:"createdAt,omitempty"`
	UpdatedAt              int64  `json:"updatedAt,omitempty"`
}

func (u *User) initByCreateRequest(req *CreateRequest) {
	u.ID = req.ID
	u.Email = req.Email
	u.Name = req.Name
	u.Role = req.Role
	u.Password = util.HashPassword(req.Password, req.Email, 10000)
}

func (u *User) initUserByID(id string) {
	u.ID = id
}

func (u *User) create() error {
	u.ResourceID = util.GenResourceID(Namespace)
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	return storage.Create(u)
}

func (u *User) findByID() error {
	return storage.FindByID(u, u.ID)
}

func NewUserContext(user *User) *UserContext {
	return &UserContext{
		ID:         user.ID,
		ResourceID: user.ResourceID,
		Role:       LookRole(user.Role),
	}
}

type UserContext struct {
	ID         string `json:"id,omitempty"`
	ResourceID string `json:"resource_id,omitempty"`
	Role       Role   `json:"role,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
}

func (u *UserContext) Valid() error {
	if !u.verifyExpiresAt() {
		return errors.New("token is expired")
	}

	return nil
}

func (u *UserContext) verifyExpiresAt() bool {
	return TimeNowFunc() > u.ExpiresAt
}

func (u *UserContext) SetExpiresAt(expiresAt int64) {
	u.ExpiresAt = expiresAt
}
