/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/services/users"
)

func testInit() {
	log.Initialize("debug")
	Initialize("secret", time.Hour*24)
}

func TestNewTokenWithUser(t *testing.T) {
	testInit()
	user := &users.User{
		ID:         "yakumioto",
		ResourceID: "users-njoVd5PKVywnZdgmhTC8EV",
		Role:       users.RoleRoot.String(),
	}
	token, err := NewTokenWithUser(user, 1636527720)
	assert.NoError(t, err, "new token error: %v", err)
	t.Logf("token is: %v", token)
}

func TestVerifyTokenWithUser(t *testing.T) {
	testInit()
	tokenString := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Inlha3VtaW90byIsInJlc291cmNlX2lkIjoidXNlcnMtbmpvVmQ1UEtWeXduWmRnbWhUQzhFViIsInJvbGUiOjIsImV4cGlyZXNfYXQiOjE2MzY1Mjc3MjB9.Oyw-59XhRxCLbXEqu7ugUSdPOXoVcp9NBlPbYW4J_Fk`
	users.TimeNowFunc = func() int64 {
		return 1636527721
	}

	_, err := VerifyTokenWithUser(tokenString)
	assert.NoError(t, err)
}
