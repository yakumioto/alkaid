package users

import (
	"github.com/yakumioto/alkaid/internal/common/util"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

type CreateRequest struct {
	ID                  string `json:"id,omitempty" validate:"required"`
	Email               string `json:"email,omitempty" validate:"required,email"`
	Password            string `json:"password,omitempty" validate:"required"`
	TransactionPassword string `json:"transactionPassword,omitempty" validate:"required"` // 交易密码仅用来加解密 PrivateKey
	Role                string `json:"role,omitempty" validate:"required"`
}

func CreateUser(request *CreateRequest) (*User, error) {
	user := &User{}
	user.ID = request.ID
	user.Email = request.Email
	user.Password = util.PBKDF2L32Sha256WithBase64(request.Password, request.Email, 10000)

	sigPrivateKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	sigPrivateKeyPem, err := crypto.PrivateKeyExport(sigPrivateKey)
	if err != nil {
		return nil, err
	}

	tlsPrivateKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

}
