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

func organizationRouter(route *gin.RouterGroup) {
	organizationDetail := "/:organizationID"
	userDetail := "/:userID"
	user := "/user"
	signca := "/signca"
	tlsca := "/tlsca"

	route.POST("/", handler.CreateOrganization)
	route.GET(organizationDetail, handler.GetOrganizationByID)

	// sign ca
	route.GET(organizationDetail+signca, handler.GetCAByOrganizationID)
	route.GET(organizationDetail+tlsca, handler.GetCAByOrganizationID)

	// msp user
	route.POST(organizationDetail+user, handler.CreateMSP)
	route.GET(organizationDetail+user+userDetail, handler.GetMSPByUserID)
}

func networkRouter(route *gin.RouterGroup) {
	networkDetail := "/:networkID"

	route.POST("", handler.CreateNetwork)
	route.GET(networkDetail, handler.GetNetworkByID)
}
