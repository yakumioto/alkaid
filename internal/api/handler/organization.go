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

func CreateOrganization(ctx *gin.Context) {
	org := types.NewOrganization()
	if err := ctx.ShouldBindJSON(org); err != nil {
		logger.Debuf("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.BadRequestData))
		return
	}

	signca, err := types.NewCA(org, types.SignCAType)
	if err != nil {
		returnInternalServerError(ctx, "New sign CA error: %v", err)
		return
	}
	tlsca, err := types.NewCA(org, types.TLSCAType)
	if err != nil {
		returnInternalServerError(ctx, "New TLS CA error: %v", err)
		return
	}

	if err := db.CreateOrganization((*db.Organization)(org)); err != nil {
		var exist *db.ErrOrganizationExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Insert organization error: %s", err)
		return
	}

	if err := db.CreateCA((*db.CA)(signca)); err != nil {
		var exist *db.ErrCAExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Create sign CA error: %v", err)
		return
	}

	if err := db.CreateCA((*db.CA)(tlsca)); err != nil {
		var exist *db.ErrCAExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Create TLS CA error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func GetOrganizationByID(ctx *gin.Context) {
	id := ctx.Param("organizationID")

	org, err := db.QueryOrganizationByOrgID(id)
	if err != nil {
		var notExist *db.ErrOrganizationNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query organization by organization_id error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, org)
}
