/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package network

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"

	apierrors "github.com/yakumioto/alkaid/internal/api/apierrors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/scheduler"
)

const (
	networkID     = "network_id"
	networkDetail = "/:" + networkID
)

var (
	logger *glog.Logger
)

type Service struct{}

func (s *Service) Init(log *glog.Logger, rg *gin.RouterGroup) {
	logger = log.MustGetLogger("network")

	r := rg.Group("/network")
	r.POST("", s.CreateNetwork)
	r.GET("")
	r.GET(networkDetail, s.GetNetworkByID)
	r.PATCH(networkDetail)
	r.DELETE(networkDetail)

	logger.Infof("Service initialization success.")
}

func (s *Service) CreateNetwork(ctx *gin.Context) {
	network := types.NewNetwork()
	if err := ctx.ShouldBindJSON(network); err != nil {
		logger.Debuf("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.BadRequest))
		return
	}

	sched, err := scheduler.NewScheduler(network.Type)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := db.CreateNetwork((*db.Network)(network)); err != nil {
		var exist *db.NetworkExistError
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataAlreadyExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := sched.CreateNetwork(network); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, network)
}

func (s *Service) GetNetworkByID(ctx *gin.Context) {
	id := ctx.Param("networkID")

	network, err := db.QueryNetworkByNetworkID(id)
	if err != nil {
		var notExist *db.NetworkNotExistError
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataNotExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, network)
}
