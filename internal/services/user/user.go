/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package user

import (
	"crypto/x509"
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

const (
	MSPTypeAdmin   = "admin"
	MSPTypeClient  = "client"
	MSPTypeOrderer = "orderer"
	MSPTypePeer    = "peer"
)

var (
	logger = glog.MustGetLogger("services.user")
)

type User struct{}

func NewService() services.Service {
	return &User{}
}

func (u *User) RootPath() string {
	return "/user"
}

func (u *User) DetailID() string {
	return "userID"
}

func (u *User) Post(ctx *gin.Context) {
	user := models.NewUser()
	if err := ctx.ShouldBindJSON(user); err != nil {
		logger.Errof("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, services.NewError(services.BadRequest))
		return
	}

	// get organization by id
	org := models.NewOrganization()
	if err := org.Get(user.OrganizationID); err != nil {
		logger.Errof("Get organization error: %v", err)
		if errors.Is(err, storage.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, services.NewError(
				services.NotFound, "Organization not found"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	// enroll sign and tls certificate
	if err := u.enrollCertificate(user, org); err != nil {
		logger.Errof("Enroll certificate error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	if err := user.Create(); err != nil {
		logger.Errof("Save database error: %v", err)
		if errors.Is(err, storage.ErrAlreadyExist) {
			ctx.JSON(http.StatusBadRequest, services.NewError(services.DataAlreadyExists))
			return
		}

		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}
}

func (u *User) Patch(ctx *gin.Context) {
	id := ctx.Param(u.DetailID())
	user := new(models.User)

	if err := ctx.ShouldBindJSON(user); err != nil {
		if !errors.As(err, &validator.ValidationErrors{}) {
			logger.Debuf("Bind JSON error: %v", err)
			ctx.JSON(http.StatusBadRequest, services.NewError(services.BadRequest))
			return
		}
	}

	if err := user.Update(id); err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, services.NewError(services.DataNotExists))
			return
		}

		logger.Errof("Update database error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *User) Get(ctx *gin.Context) {
	id := ctx.Param(u.DetailID())
	user := new(models.User)

	if err := user.Get(id); err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, services.NewError(services.DataNotExists))
			return
		}

		logger.Errof("Get error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *User) GetList(ctx *gin.Context) {
	users := new(models.Users)

	if err := users.Query(); err != nil {
		logger.Errof("Get list error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (u *User) Delete(ctx *gin.Context) {
	id := ctx.Param(u.DetailID())
	user := new(models.User)

	if err := user.Delete(id); err != nil {
		logger.Errof("Delete error: %v", err)
		ctx.JSON(http.StatusInternalServerError, services.NewError(services.InternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *User) enrollCertificate(user *models.User, org *models.Organization) error {
	signPriv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("generate sign private key error: %v", err)
	}

	tlsPriv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("generate tls private key error: %v", err)
	}

	var (
		commonName  string
		keyUsage    x509.KeyUsage
		extKeyUsage []x509.ExtKeyUsage
	)
	switch user.MSPType {
	case MSPTypeAdmin, MSPTypeClient:
		commonName = fmt.Sprintf("%s@%s", user.ID, org.Domain)
		keyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		extKeyUsage = []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}
	case MSPTypeOrderer, MSPTypePeer:
		commonName = fmt.Sprintf("%s.%s", user.ID, org.Domain)
		keyUsage = x509.KeyUsageDigitalSignature
	}

	signCert, err := certificate.SignCertificate(
		utils.GetPkixName(org, commonName),
		nil,
		&signPriv.PublicKey,
		keyUsage,
		extKeyUsage,
		org.CAPrivateKey,
		org.CACertificate)
	if err != nil {
		return fmt.Errorf("enroll sign certificate error: %v", err)
	}
	user.SignPrivateKey, _ = crypto.PrivateKeyExport(signPriv)
	user.SignCertificate = crypto.X509Export(signCert)

	tlsCert, err := certificate.SignCertificate(
		utils.GetPkixName(org, commonName),
		user.SANs,
		&signPriv.PublicKey,
		keyUsage,
		extKeyUsage,
		org.CAPrivateKey,
		org.CACertificate)
	if err != nil {
		return fmt.Errorf("enroll tls certificate error: %v", err)
	}
	user.TLSPrivateKey, _ = crypto.PrivateKeyExport(tlsPriv)
	user.TLSCertificate = crypto.X509Export(tlsCert)

	return nil
}
