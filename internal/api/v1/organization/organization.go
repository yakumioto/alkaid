/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package organization

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yakumioto/glog"

	apierrors "github.com/yakumioto/alkaid/internal/api/apierrors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/utils/certificate"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

const (
	organizationID     = "orgnaization_id"
	organizationDetail = "/:" + organizationID
)

var (
	logger *glog.Logger
)

type Service struct{}

func (s *Service) Init(log *glog.Logger, rg *gin.RouterGroup) {
	logger = log.MustGetLogger("organization")

	r := rg.Group("/organization")
	r.POST("", s.CreateOrganization)
	r.GET("")
	r.GET(organizationDetail, s.GetOrganizationByID)
	r.PATCH(organizationDetail)
	r.DELETE(organizationDetail)

	logger.Infof("Service initialization success.")
}

func (s *Service) CreateOrganization(ctx *gin.Context) {
	org := types.NewOrganization()
	if err := ctx.ShouldBindJSON(org); err != nil {
		logger.Debuf("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.BadRequest))
		return
	}

	priv, cert, err := certificate.NewCA(org, fmt.Sprintf("ca.%s", org.Domain))
	if err != nil {
		logger.Errof("Generating signature root certificate error: %v", err)
		ctx.JSON(http.StatusInternalServerError, apierrors.New(
			apierrors.InternalServerError, "Generating signature root certificate error"))
		return
	}
	org.SignCAPrivateKey, _ = crypto.PrivateKeyExport(priv)
	org.SignCACertificate = crypto.X509Export(cert)

	priv, cert, err = certificate.NewCA(org, fmt.Sprintf("tlsca.%s", org.Domain))
	if err != nil {
		logger.Errof("Generating TLS root certificate error: %v", err)
		ctx.JSON(http.StatusInternalServerError, apierrors.New(
			apierrors.InternalServerError, "Generating TLS root certificate error"))
		return
	}
	org.TLSCAPrivateKey, _ = crypto.PrivateKeyExport(priv)
	org.TLSCACertificate = crypto.X509Export(cert)

	if err := db.CreateOrganization((*db.Organization)(org)); err != nil {
		var exist *db.ErrOrganizationExist
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataAlreadyExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (s *Service) GetOrganizationByID(ctx *gin.Context) {
	id := ctx.Param("organizationID")

	org, err := db.QueryOrganizationByOrgID(id)
	if err != nil {
		var notExist *db.ErrOrganizationNotExist
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataNotExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, org)
}
