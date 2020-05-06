/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/yakumioto/alkaid/internal/api/handler"
)

func Init(r *gin.Engine) {
	network := r.Group("/network")
	networkRouter(network)

	organization := r.Group("/organization")
	organizationRouter(organization)
}

func organizationRouter(r *gin.RouterGroup) {
	organization := "/"
	organizationDetail := "/:organizationID"

	signca := "/signca"
	tlsca := "/tlsca"

	user := "/user"
	userDetail := "/:userID"

	peer := "/peer"
	peerDetail := "/:peerID"

	orderer := "/orderer"
	ordererDetail := "/:ordererID"

	r.POST(organization, handler.CreateOrganization)
	r.GET(organizationDetail, handler.GetOrganizationByID)

	// sign ca
	r.GET(organizationDetail+signca, handler.GetCAByOrganizationID)
	r.GET(organizationDetail+tlsca, handler.GetCAByOrganizationID)

	// msp user
	r.POST(organizationDetail+user, handler.CreateMSP)
	r.GET(organizationDetail+user+userDetail, handler.GetMSPByUserID)

	// peer
	r.POST(organizationDetail+peer, nil)
	r.GET(organizationDetail+peer+peerDetail, nil)

	// orderer
	r.POST(organizationDetail+orderer, nil)
	r.GET(organizationDetail+orderer+ordererDetail, nil)
}

func networkRouter(r *gin.RouterGroup) {
	network := "/"
	networkDetail := "/:networkID"

	r.POST(network, handler.CreateNetwork)
	r.GET(networkDetail, handler.GetNetworkByID)
}
