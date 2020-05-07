/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
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

func CreateMSP(ctx *gin.Context) {
	orgid := ctx.Param("organizationID")

	msp := types.NewMSP()
	if err := ctx.ShouldBindJSON(msp); err != nil {
		logger.Debuf("Bind JSON error: %v", err)

		ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.BadRequestData))
		return
	}
	msp.OrganizationID = orgid

	org, err := db.QueryOrganizationByOrgID(orgid)
	if err != nil {
		var notExist *db.ErrOrganizationNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Get Organization error: %v", err)
		return
	}

	signCA, err := db.QueryCAByOrganizationIDAndType(orgid, types.SignCAType)
	if err != nil {
		var notExist *db.ErrCANotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Get Organization error: %v", err)
		return
	}

	tlsCA, err := db.QueryCAByOrganizationIDAndType(orgid, types.TLSCAType)
	if err != nil {
		var notExist *db.ErrCANotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Get Organization error: %v", err)
		return
	}

	if err := msp.Initialize(org, signCA, tlsCA); err != nil {
		returnInternalServerError(ctx, "MSP initialize error: %s", err)
		return
	}

	if err := db.CreateMSP((*db.MSP)(msp)); err != nil {
		var exist *db.ErrMSPExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Insert MSP error: %s", err)
		return
	}

	ctx.JSON(http.StatusOK, msp)
}

func GetMSPByUserID(ctx *gin.Context) {
	orgid := ctx.Param("organizationID")
	userid := ctx.Param("userID")

	msp, err := db.QueryMSPByOrganizationIDAndUserID(orgid, userid)
	if err != nil {
		var notExist *db.ErrMSPNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query MSP by organization_id and user_id error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, msp)
}
