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
	ResourceNamespace = "User"

	RoleRoot Role = iota
	RoleOrganization
	RoleNetwork
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
	"root":         RoleRoot,
	"organization": RoleOrganization,
	"network":      RoleNetwork,
	"user":         RoleUser,
}

func LookRole(str string) Role {
	role, ok := roleMap[str]
	if !ok {
		return RoleUser
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

func (r Role) Ge(role Role) bool {
	return role >= r
}

func (r Role) Le(role Role) bool {
	return role <= r
}

type User struct {
	ID                     string `json:"id,omitempty"`
	ResourceID             string `json:"resourceId,omitempty"`
	OrganizationID         string `json:"organizationId,omitempty"`
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

func (u *User) initByCreateRequest(req *CreateRequest, userCtx *UserContext) error {
	if !userCtx.validRole(req.Role) {
		return errors.New("verifying new user role error")
	}

	// 只有 root 用户有权限创建其他组织的 user
	if !userCtx.validOrganization() {
		req.OrganizationID = userCtx.OrganizationID
	}

	u.ID = req.ID
	u.OrganizationID = req.OrganizationID
	u.Email = req.Email
	u.Name = req.Name
	u.Role = req.Role
	u.Password = util.HashPassword(req.Password, req.Email, 10000)

	return nil
}

func (u *User) initUserByID(id string) {
	u.ID = id
}

func (u *User) create() error {
	u.ResourceID = util.GenResourceID(ResourceNamespace)
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	return storage.Create(u)
}

func (u *User) findByID() error {
	return storage.FindByID(u, u.ID)
}

func NewUserContext(user *User) *UserContext {
	return &UserContext{
		ID:             user.ID,
		ResourceID:     user.ResourceID,
		OrganizationID: user.OrganizationID,
		Role:           LookRole(user.Role),
	}
}

type UserContext struct {
	ID             string `json:"id"`
	ResourceID     string `json:"resourceId"`
	OrganizationID string `json:"organizationId"`
	Role           Role   `json:"role"`
	ExpiresAt      int64  `json:"expiresAt"`
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

func (u *UserContext) validRole(role string) bool {
	if !u.Role.Le(LookRole(role)) {
		return false
	}

	return true
}

func (u *UserContext) validOrganization() bool {
	if u.Role == RoleRoot {
		return true
	}

	return true
}
