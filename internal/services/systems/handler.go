/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package systems

import (
	"net/http"

	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/errors"
	"github.com/yakumioto/alkaid/internal/services/users"
)

var (
	logger = log.GetPackageLogger("services.system")
)

type InitRequest struct {
	ID       string `json:"id,omitempty" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func SystemInit(req *InitRequest) (*System, error) {
	initialized := newSystemByID(KSystemInitialized)
	if err := initialized.findByID(); err != nil {
		if err != storage.ErrNotFound {
			logger.Errorf("query system [%v] error: %v", KSystemInitialized, err)
			return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
				"server unknown error")
		}
	}

	if initialized.Value == VSystemInitialized {
		logger.Warnln("system is initialized")
		return nil, errors.NewError(http.StatusForbidden, errors.ErrForbidden,
			"system is initialized")
	}

	// 处理未初始化情况
	user, err := users.Create(&users.CreateRequest{
		ID:       req.ID,
		Name:     req.Name,
		Email:    req.Email,
		Root:     true,
		Password: req.Password,
	})
	if err != nil {
		logger.Errorf("[%v] initialize root user error: %v", user.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to initialize root user")
	}

	sys := newSystem(KSystemInitialized, VSystemInitialized)
	if err := sys.create(); err != nil {
		logger.Errorf("[%v] create system error: %v", sys.Key, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to create user")
	}
	return sys, nil
}
