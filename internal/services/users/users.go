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

// User 实体用户，每个用户会生成两对公私密钥，用于签名以及通讯认证。
// 所以在创建用户时需要填入一个交易密码，此交易密码用来加解密上述两对密钥。
type User struct {
	ResourceID              string `json:"resourceID,omitempty" gorm:"primaryKey"`
	ID                      string `json:"id,omitempty" gorm:"uniqueIndex"`
	Name                    string `json:"name,omitempty"`
	Email                   string `json:"email,omitempty" gorm:"uniqueIndex"`
	Password                string `json:"-"`
	Root                    bool   `json:"root,omitempty"`
	ProtectedSignPrivateKey string `json:"protectedSignPrivateKey,omitempty"`
	ProtectedTLSPrivateKey  string `json:"protectedTlsPrivateKey,omitempty"`
	Deactivate              bool   `json:"deactivate,omitempty"`
	CreatedAt               int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt               int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
	DeactivateAt            int64  `json:"deactivateAt,omitempty"`
}

func newUserByCreateRequest(req *CreateRequest, userCtx *UserContext) (*User, error) {
	return &User{
		ID:       req.ID,
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}, nil
}

func (u *User) Create() error {
	u.ResourceID = util.GenResourceID(ResourceNamespace)
	u.Password = util.HashPassword(u.Password, u.Email, 10000)
	return storage.Create(u)
}

func FindByIDOrEmailOrResourceID(id string) (*User, error) {
	user := new(User)
	return user, storage.FindByQuery(user,
		storage.NewQueryOptions().
			Or(&User{ID: id}).
			Or(&User{Email: id}).
			Or(&User{ResourceID: id}))
}

type UserContext struct {
	ID            string `json:"id,omitempty"`
	ResourceID    string `json:"resourceID,omitempty"`
	Root          bool   `json:"root,omitempty"`
	Organizations struct {
		OrganizationResourceID string `json:"organizationResourceID,omitempty"`
		Role                   Role   `json:"role"`
	}
	ExpiredAt int64 `json:"expiredAt"`
}

func NewUserContext(user *User) *UserContext {
	return &UserContext{
		ID:         user.ID,
		ResourceID: user.ResourceID,
	}
}

func (u *UserContext) Valid() error {
	if !u.verifyExpiresAt() {
		return errors.New("token is expired")
	}

	return nil
}

func (u *UserContext) verifyExpiresAt() bool {
	return TimeNowFunc() < u.ExpiredAt
}

func (u *UserContext) SetExpiresAt(expiredAt int64) {
	u.ExpiredAt = expiredAt
}

type UserOrganizations struct {
	ResourceID             string `json:"resourceId,omitempty" gorm:"primaryKey"`
	UserResourceID         string `json:"userResourceID,omitempty" gorm:"index"`
	OrganizationResourceID string `json:"organizationResourceID,omitempty" gorm:"index"`
	Role                   Role   `json:"role,omitempty"`
	Status                 string `json:"status,omitempty"`
	Deactivate             bool   `json:"deactivate,omitempty"`
	CreatedAt              int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt              int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
	DeactivateAt           int64  `json:"deactivateAt,omitempty"`
}
