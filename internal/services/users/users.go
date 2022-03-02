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

const ResourceNamespace = "User"

const (
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
	"organization",
	"network",
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

func (r Role) GE(role Role) bool {
	return r >= role
}

func (r Role) LE(role Role) bool {
	return r <= role
}

type User struct {
	ID                     string `json:"id,omitempty" gorm:"primaryKey"`
	ResourceID             string `json:"resourceId,omitempty" gorm:"uniqueIndex"`
	OrganizationID         string `json:"organizationId,omitempty"`
	Name                   string `json:"name,omitempty"`
	Email                  string `json:"email,omitempty" gorm:"uniqueIndex"`
	Password               string `json:"-"`
	Role                   string `json:"role,omitempty"`
	ProtectedSigPrivateKey string `json:"protectedSigPrivateKey,omitempty"`
	ProtectedTLSPrivateKey string `json:"protectedTlsPrivateKey,omitempty"`
	Status                 string `json:"status,omitempty"`
	CreatedAt              int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt              int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
}

func newUserByCreateRequest(req *CreateRequest, userCtx *UserContext) (*User, error) {
	if !userCtx.validRole(req.Role) {
		return nil, errors.New("verifying new user role error")
	}

	// 只有 root 用户有权限创建其他组织的 user
	if !userCtx.validOrganization() {
		req.OrganizationID = userCtx.OrganizationID
	}

	return &User{
		ID:             req.ID,
		OrganizationID: req.OrganizationID,
		Email:          req.Email,
		Name:           req.Name,
		Role:           req.Role,
		Password:       req.Password,
	}, nil
}

func (u *User) create() error {
	u.ResourceID = util.GenResourceID(ResourceNamespace)
	u.Password = util.HashPassword(u.Password, u.Email, 10000)
	return storage.Create(u)
}

func FindByIDOrEmail(id string) (*User, error) {
	user := new(User)
	return user, storage.FindByQuery(user, storage.NewQueryOptions().
		SetWhere(&User{ID: id}).
		SetOr(&User{Email: id}))
}

type UserContext struct {
	ID             string `json:"id"`
	ResourceID     string `json:"resourceId"`
	OrganizationID string `json:"organizationId"`
	Role           Role   `json:"role"`
	ExpiresAt      int64  `json:"expiresAt"`
}

func NewUserContext(user *User) *UserContext {
	return &UserContext{
		ID:             user.ID,
		ResourceID:     user.ResourceID,
		OrganizationID: user.OrganizationID,
		Role:           LookRole(user.Role),
	}
}

func (u *UserContext) Valid() error {
	if !u.verifyExpiresAt() {
		return errors.New("token is expired")
	}

	return nil
}

func (u *UserContext) verifyExpiresAt() bool {
	return TimeNowFunc() < u.ExpiresAt
}

func (u *UserContext) SetExpiresAt(expiresAt int64) {
	u.ExpiresAt = expiresAt
}

func (u *UserContext) validRole(role string) bool {
	logger.Tracef("user context role is %v, create user role is: %v", u.Role.String(), role)
	if !u.Role.LE(LookRole(role)) {
		return false
	}
	return true
}

func (u *UserContext) validOrganization() bool {
	if u.Role == RoleRoot {
		return true
	}
	return false
}
