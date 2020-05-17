/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package node

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/apierrors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/config"
	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/scheduler"
	"github.com/yakumioto/alkaid/internal/utils/cryptoconfig"
	"github.com/yakumioto/alkaid/internal/vm"
)

const (
	nodeID     = "node_id"
	nodeDetail = "/:" + nodeID
)

var (
	logger *glog.Logger
)

type Service struct{}

func (s *Service) Init(log *glog.Logger, rg *gin.RouterGroup) {
	logger = log.MustGetLogger("node")

	r := rg.Group("/node")
	r.POST("", s.Create)
	r.GET("")
	r.GET(nodeDetail)
	r.PATCH(nodeDetail)
	r.DELETE(nodeDetail)

	logger.Infof("Service initialization success.")
}

func (s *Service) Create(ctx *gin.Context) {
	node := types.NewNode()
	if err := ctx.ShouldBindJSON(node); err != nil {
		logger.Errof("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.BadRequest))
		return
	}

	org, err := db.QueryOrganizationByOrgID(node.OrganizationID)
	if s.checkQueryError(ctx, err) {
		return
	}

	user, err := db.QueryMSPByOrganizationIDAndUserID(node.OrganizationID, node.UserID)
	if s.checkQueryError(ctx, err) {
		return
	}

	network, err := db.QueryNetworkByNetworkID(node.NetworkID)
	if s.checkQueryError(ctx, err) {
		return
	}

	if !org.HasNetwork(network.NetworkID) {
		ctx.JSON(http.StatusBadRequest,
			apierrors.New(apierrors.BadRequest, "The organization is not on the current network"))
		return
	}

	msp, err := cryptoconfig.GetMSPArchive(org, user)
	if err != nil {
		logger.Errof("Get msp archive error: %v", err)
		ctx.JSON(http.StatusInternalServerError, apierrors.New(apierrors.InternalServerError))
		return
	}

	tls, err := cryptoconfig.GetTLSArchive(org, user)
	if err != nil {
		logger.Errof("Get msp archive error: %v", err)
		ctx.JSON(http.StatusInternalServerError, apierrors.New(apierrors.InternalServerError))
		return
	}

	dockerid := fmt.Sprintf("%s.%s", user.UserID, org.Domain)
	switch user.MSPType {
	case types.PeerMSPType:
		if err := s.createPeerNode(dockerid, org.OrganizationID, network.GetNetworkID(), node.CouchDB, msp, tls); err != nil {
			logger.Errof("Create peer node error: %v", err)
			ctx.JSON(http.StatusServiceUnavailable, apierrors.New(apierrors.InternalServerError))
			return
		}
	default:
		logger.Errof("Unknown msp type: %s", user.MSPType)
		ctx.JSON(http.StatusServiceUnavailable, apierrors.New(apierrors.InternalServerError))
		return
	}

	if node.DockerContainerIDs == nil {
		node.DockerContainerIDs = make([]string, 0)
	}

	node.NodeID = dockerid
	node.Type = user.MSPType
	node.Status = types.RunningNodeStatus
	node.DockerContainerIDs = append(node.DockerContainerIDs, dockerid)
	if node.CouchDB {
		node.DockerContainerIDs = append(node.DockerContainerIDs, fmt.Sprintf("couchdb.%s", dockerid))
	}

	if err := db.CreateNode((*db.Node)(node)); err != nil {
		if errors.Is(err, db.ErrNodeExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataAlreadyExists))
			return
		}

		logger.Errof("Save Node to database error: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, node)
}

func (s *Service) createPeerNode(name, mspid, networkMode string, couchdb bool, msp, tls []byte) error {
	requests := make([]*vm.CreateRequest, 0)

	peer := &vm.CreateRequest{
		ContainerName: name,
		ImageName:     config.PeerImageName,
		ImageTag:      config.PeerImageVersion,
		NetworkMode:   networkMode,
		Environment: []string{
			fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspid),
		},
		VolumeMounts: map[string]string{
			name:                        "/var/hyperledger/production",
			fmt.Sprintf("%s.msp", name): "/etc/hyperledger/fabric/msp",
			fmt.Sprintf("%s.tls", name): "/etc/hyperledger/fabric/tls",
		},
		Files: map[string][]byte{
			"/etc/hyperledger/fabric/msp": msp,
			"/etc/hyperledger/fabric/tls": tls,
		},
	}

	if couchdb {
		address := fmt.Sprintf("%s.%s", "couch", name)

		requests = append(requests, &vm.CreateRequest{
			ContainerName: address,
			ImageName:     config.CouchDBImageName,
			ImageTag:      config.CouchDBImageVersion,
			// todo: Custom user name and password
		})

		peer.Environment = append(peer.Environment,
			"CORE_LEDGER_STATE_STATEDATABASE=CouchDB",
			fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=%s:5984", address),
		)
	}

	sche, err := scheduler.NewScheduler(types.DockerNetworkType)
	if err != nil {
		logger.Errof("New scheduler error: %v", err)
		return err
	}
	return sche.CreatePeer(peer, requests...)
}

func (s *Service) checkQueryError(ctx *gin.Context, err error) bool {
	if err != nil {
		switch {
		case errors.Is(err, db.ErrOrganizationNotExist):
			logger.Infof("Organization not exist.")
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataNotExists, "Organization not exist"))
		case errors.Is(err, db.ErrUserNotExist):
			logger.Infof("User not exist.")
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataNotExists, "User not exist"))
		case errors.Is(err, db.ErrNetworkNotExist):
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataNotExists, "Network not exist"))
			logger.Infof("Network not exist.")
		default:
			ctx.JSON(http.StatusInternalServerError, apierrors.New(apierrors.InternalServerError))
			logger.Errof("Query error: %v", err)
		}

		return true
	}

	return false
}
