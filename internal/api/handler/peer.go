/*
 *  Copyright 2020. The Alkaid Authors. All rights reserved.
 *  Use of this source code is governed by a MIT-style
 *  license that can be found in the LICENSE file.
 *  Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apierrors "github.com/yakumioto/alkaid/internal/api/errors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/db"
)

type createPeerNodeRequest struct {
	OrganizationID string `json:"organization_id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
}

func CreatePeerNode(ctx *gin.Context) {
	req := new(createPeerNodeRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Debuf("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.BadRequestData))
		return
	}

	// query organization by id
	org, err := db.QueryOrganizationByOrgID(req.OrganizationID)
	if err != nil {
		var notExist *db.ErrOrganizationNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query organization by organization_id error: %v", err)
		return
	}

	// query user by id
	msp, err := db.QueryMSPByOrganizationIDAndUserID(req.OrganizationID, req.UserID)
	if err != nil {
		var notExist *db.ErrMSPNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query MSP by organization_id and user_id error: %v", err)
		return
	}

	// todo: Accurate error message
	if org.Type != types.PeerOrgType || msp.Type != types.PeerMSPType {
		ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.BadRequestData))
		return
	}
}
