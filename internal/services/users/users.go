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

	"github.com/yakumioto/alkaid/internal/common/crypto/utils"
	"github.com/yakumioto/alkaid/internal/common/storage"
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

type Role int

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

// User 实体用户，每个用户会生成三对公私密钥，用于签名，通讯认证以及组织对称密钥管理。
type User struct {
	ResourceID              string `json:"resourceId,omitempty" gorm:"primaryKey"`
	UserID                  string `json:"userId,omitempty" gorm:"uniqueIndex"`
	Name                    string `json:"name,omitempty"`
	Email                   string `json:"email,omitempty" gorm:"uniqueIndex"`
	Password                string `json:"-"`
	Root                    bool   `json:"root,omitempty"`
	ProtectedSymmetricKey   string `json:"protectedSymmetricKey,omitempty"`
	ProtectedSignPrivateKey string `json:"protectedSignPrivateKey,omitempty"`
	SignPublicKey           string `json:"signPublicKey,omitempty"`
	ProtectedTLSPrivateKey  string `json:"protectedTlsPrivateKey,omitempty"`
	TLSPublicKey            string `json:"tlsPublicKey,omitempty"`
	ProtectedRSAPrivateKey  string `json:"protectedRSAPrivateKey,omitempty"`
	RSAPublicKey            string `json:"rsaPublicKey,omitempty"`
	Deactivate              bool   `json:"deactivate,omitempty"`
	CreatedAt               int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt               int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
	DeactivateAt            int64  `json:"deactivateAt,omitempty"`
}

func newUserByCreateRequest(req *CreateRequest) *User {
	return &User{
		UserID:   req.UserID,
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}
}

func (u *User) Create() error {
	u.ResourceID = utils.GenResourceID(ResourceNamespace)
	u.Password = utils.HashPassword(
		string(utils.GetMasterKey(u.Password, u.Email)), u.Password, 1)
	return storage.Create(u)
}

func FindUserByID(id string) (*User, error) {
	user := new(User)
	return user, storage.FindByQuery(user,
		storage.NewQueryOptions().
			Or(&User{UserID: id}).
			Or(&User{Email: id}).
			Or(&User{ResourceID: id}))
}

type UserOrganizations struct {
	ResourceID     string `json:"resourceId,omitempty" gorm:"primaryKey"`
	UserID         string `json:"UserId,omitempty" gorm:"index"`
	OrganizationID string `json:"OrganizationId,omitempty" gorm:"index"`
	Role           Role   `json:"role,omitempty"`
	Status         string `json:"status,omitempty"`
	Deactivate     bool   `json:"deactivate,omitempty"`
	CreatedAt      int64  `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt      int64  `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
	DeactivateAt   int64  `json:"deactivateAt,omitempty"`
}

func FindUserOrganizationsByUserID(id string) ([]*UserOrganizations, error) {
	organizations := make([]*UserOrganizations, 0)
	return organizations, storage.FindByQuery(&organizations,
		storage.NewQueryOptions().
			Where(UserOrganizations{UserID: id}))
}

type UserContext struct {
	ID            string          `json:"id,omitempty"`
	Root          bool            `json:"root,omitempty"`
	Organizations []*organization `json:"organizations,omitempty"`
	ExpiredAt     int64           `json:"expiredAt"`
}

type organization struct {
	OrganizationID string `json:"organizationId,omitempty"`
	Role           Role   `json:"role"`
}

func NewUserContext(user *User, orgs []*UserOrganizations) *UserContext {
	organizations := make([]*organization, 0)
	for _, org := range orgs {
		organizations = append(organizations, &organization{
			OrganizationID: org.OrganizationID,
			Role:           org.Role,
		})
	}

	return &UserContext{
		ID:            user.UserID,
		Root:          user.Root,
		Organizations: organizations,
	}
}

func (u *UserContext) Role(orgID string) string {
	if u.Root {
		return RoleRoot.String()
	}

	for _, org := range u.Organizations {
		if org.OrganizationID == orgID {
			return org.Role.String()
		}
	}

	return RoleUser.String()
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
