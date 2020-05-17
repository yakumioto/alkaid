/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/apierrors"
	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/utils/certificate"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

const (
	userID     = "user_id"
	userDetail = "/:" + userID
)

var (
	logger *glog.Logger
)

type Service struct{}

func (s *Service) Init(log *glog.Logger, rg *gin.RouterGroup) {
	logger = log.MustGetLogger("user")
	r := rg.Group("/user")
	r.POST("", s.CreateUser)
	r.GET("")
	r.GET(userDetail, s.GetUserByID)
	r.PATCH(userDetail)
	r.DELETE(userDetail)

	logger.Infof("Service initialization success.")
}

func (s *Service) CreateUser(ctx *gin.Context) {
	user := types.NewUser()
	if err := ctx.ShouldBindJSON(user); err != nil {
		logger.Errof("Bind JSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.BadRequest))
		return
	}

	org, err := db.QueryOrganizationByOrgID(user.OrganizationID)
	if err != nil {
		var notExist *db.OrganizationNotExistError
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.New(apierrors.DataNotExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := s.EnrollCertificate(user, org); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := db.CreateMSP((*db.User)(user)); err != nil {
		var exist *db.UserExistError
		if errors.As(err, &exist) {
			ctx.JSON(http.StatusBadRequest, apierrors.New(apierrors.DataAlreadyExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Service) GetUserByID(ctx *gin.Context) {
	orgid := ctx.Param("organizationID")
	userid := ctx.Param("userID")

	msp, err := db.QueryMSPByOrganizationIDAndUserID(orgid, userid)
	if err != nil {
		var notExist *db.UserNotExistError
		if errors.As(err, &notExist) {
			ctx.JSON(http.StatusNotFound, apierrors.New(apierrors.DataNotExists))
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, msp)
}

func (s *Service) EnrollCertificate(user *types.User, org *types.Organization) error {
	signPriv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("generate sign private key error: %v", err)
	}

	tlsPriv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("generate tls private key error: %v", err)
	}

	commonName := ""
	switch user.MSPType {
	case types.AdminMSPType, types.ClientMSPType:
		commonName = fmt.Sprintf("%s@%s", user.UserID, org.Domain)
	case types.OrdererMSPType, types.PeerMSPType:
		commonName = fmt.Sprintf("%s.%s", user.UserID, org.Domain)
	}

	signCert, err := certificate.SignCertificate(org, commonName, user.MSPType, nil,
		&signPriv.PublicKey, org.SignCAPrivateKey, org.CACertificate)
	if err != nil {
		return fmt.Errorf("sign signature certificate error: %v", err)
	}
	tlsCert, err := certificate.SignCertificate(org, commonName, user.MSPType, user.SANS,
		&tlsPriv.PublicKey, org.SignCAPrivateKey, org.CACertificate)
	if err != nil {
		return fmt.Errorf("sign tls certificate error: %v", err)
	}

	signPrivBytes, err := crypto.PrivateKeyExport(signPriv)
	if err != nil {
		return fmt.Errorf("sign private key export error: %v", err)
	}

	tlsPrivBytes, err := crypto.PrivateKeyExport(tlsPriv)
	if err != nil {
		return fmt.Errorf("tls private key export error: %v", err)
	}

	user.SignPrivateKey = signPrivBytes
	user.TLSPrivateKey = tlsPrivBytes
	user.SignCertificate = crypto.X509Export(signCert)
	user.TLSCertificate = crypto.X509Export(tlsCert)

	return nil
}
