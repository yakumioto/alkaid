/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package routers

import "github.com/gin-gonic/gin"

const (
	network   = "/network"
	networkID = "/:networkID"
	orderer   = network + networkID + "/orderer"
	ordererID = "/:ordererID"
	peer      = network + networkID + "/peer"
	peerID    = "/:peerID"
)

func AddRouters(r *gin.Engine) {
	r.POST(network)
	r.GET(network)
	r.GET(network + networkID)
	r.DELETE(network + networkID)

	r.POST(orderer)
	r.GET(orderer)
	r.GET(orderer + ordererID)
	r.DELETE(orderer + ordererID)

	r.POST(peer)
	r.GET(peer)
	r.GET(peer + peerID)
	r.DELETE(peer + peerID)
}
