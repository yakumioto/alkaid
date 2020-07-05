/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package organization

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/models"
	"github.com/yakumioto/alkaid/internal/services"
	"github.com/yakumioto/alkaid/internal/services/utils"
	"github.com/yakumioto/alkaid/internal/storage"
	"github.com/yakumioto/alkaid/internal/utils/certificate"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

var (
	logger = glog.MustGetLogger("services.organization")
)

type Service struct{}

func NewService() services.Service {
	return &Service{}
}

func (o *Service) RootPath() string {
	return "/organization"
}

func (o *Service) DetailID() string {
	return "organizationID"
}

func (o *Service) Post(ctx *gin.Context) {
	org := models.NewOrganization()
	if err := ctx.ShouldBindJSON(org); err != nil {
		logger.Errof("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, services.NewError(services.BadRequest))
		return
	}

	// generate root ca
	priv, cert, err := certificate.NewCA(utils.GetPkixName(org, fmt.Sprintf("ca.%s", org.Domain)))
	if err != nil {
		logger.Errof("Generating signature root certificate error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(
			services.InternalServerError, "Generating signature root certificate error"))
		return
	}
	org.CAPrivateKey, _ = crypto.PrivateKeyExport(priv)
	org.CACertificate = crypto.X509Export(cert)

	// generate root tlsca
	priv, cert, err = certificate.NewCA(utils.GetPkixName(org, fmt.Sprintf("tlsca.%s", org.Domain)))
	if err != nil {
		logger.Errof("Generating TLS root certificate error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(
			services.InternalServerError, "Generating tls root certificate error"))
		return
	}
	org.TLSCAPrivateKey, _ = crypto.PrivateKeyExport(priv)
	org.TLSCACertificate = crypto.X509Export(cert)

	if err := org.Create(); err != nil {
		logger.Errof("Save database error: %v", err)
		if errors.Is(err, storage.ErrAlreadyExist) {
			ctx.JSON(http.StatusBadRequest, services.NewError(services.DataAlreadyExists))
			return
		}

		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusCreated, org)
}

func (o *Service) Patch(ctx *gin.Context) {
	id := ctx.Param(o.DetailID())
	org := new(models.Organization)

	if err := ctx.ShouldBindJSON(org); err != nil {
		if !errors.As(err, &validator.ValidationErrors{}) {
			logger.Debuf("Bind JSON error: %v", err)
			ctx.JSON(http.StatusBadRequest, services.NewError(services.BadRequest))
			return
		}
	}

	if err := org.Update(id); err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, services.NewError(services.DataNotExists))
			return
		}

		logger.Errof("Update database error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (o *Service) Get(ctx *gin.Context) {
	id := ctx.Param(o.DetailID())
	org := new(models.Organization)

	if err := org.Get(id); err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, services.NewError(services.DataNotExists))
			return
		}

		logger.Errof("Get error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (o *Service) GetList(ctx *gin.Context) {
	orgs := new(models.Organizations)

	if err := orgs.Query(); err != nil {
		logger.Errof("Get list error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, orgs)
}

func (o *Service) Delete(ctx *gin.Context) {
	id := ctx.Param(o.DetailID())
	org := new(models.Organization)

	if err := org.Delete(id); err != nil {
		logger.Errof("Delete error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, org)
}
