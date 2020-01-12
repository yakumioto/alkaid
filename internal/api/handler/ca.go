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
	"strings"

	"github.com/gin-gonic/gin"

	apierrors "github.com/yakumioto/alkaid/internal/api/errors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/db"
)

func CreateCA(ctx *gin.Context) {
	org, err := db.QueryOrganizationByOrgID(ctx.Param("organizationID"))
	if err != nil {
		var notExist *db.ErrOrganizationNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Get Organization error: %v", err)
		return
	}

	ca, err := types.NewCA(org, getType(ctx.Request.RequestURI))
	if err != nil {
		returnInternalServerError(ctx, "New CA error: %v", err)
		return
	}

	if err := db.CreateCA((*db.CA)(ca)); err != nil {
		var exist *db.ErrCAExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Create CA error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, ca)
}

func GetCAByOrganizationID(ctx *gin.Context) {
	id := ctx.Param("organizationID")
	ca, err := db.QueryCAByOrganizationIDAndType(id, getType(ctx.Request.RequestURI))
	if err != nil {
		var notExist *db.ErrCANotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query CA unknown error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, ca)
}

func getType(uri string) string {
	if strings.Contains(uri, "signca") {
		return types.SignCAType
	}

	if strings.Contains(uri, "tlsca") {
		return types.TLSCAType
	}

	return ""
}
