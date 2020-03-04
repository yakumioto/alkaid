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
	"github.com/yakumioto/alkaid/internal/scheduler"
)

func CreateNetwork(ctx *gin.Context) {
	network := types.NewNetwork()
	if err := ctx.ShouldBindJSON(network); err != nil {
		logger.Debuf("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.BadAuthenticationData))
		return
	}

	sched, err := scheduler.NewScheduler(network.Type)
	if err != nil {
		returnInternalServerError(ctx, "New scheduler error: %v", err)
		return
	}

	if err := db.CreateNetwork((*db.Network)(network)); err != nil {
		var exist *db.ErrNetworkExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataAlreadyExists))
			return
		}

		returnInternalServerError(ctx, "Insert network error: %v", err)
		return
	}

	if err := sched.CreateNetwork(network); err != nil {
		returnInternalServerError(ctx, "Scheduler create network error: %s", err)
		return
	}

	ctx.JSON(http.StatusOK, network)
}

func GetNetworkByID(ctx *gin.Context) {
	id := ctx.Param("networkID")

	network, err := db.QueryNetworkByNetworkID(id)
	if err != nil {
		var notExist *db.ErrNetworkNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.NewErrors(apierrors.DataNotExists))
			return
		}

		returnInternalServerError(ctx, "Query network by network_id error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, network)
}
