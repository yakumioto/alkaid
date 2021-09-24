package users

type CreateRequest struct {
	ID                  string `json:"id,omitempty" validate:"required"`
	Email               string `json:"email,omitempty" validate:"required,email"`
	Password            string `json:"password,omitempty" validate:"required"`
	TransactionPassword string `json:"transactionPassword,omitempty" validate:"required"` // 交易密码仅用来加解密 PrivateKey
	Role                string `json:"role,omitempty" validate:"required"`
}

func CreateUser(req *CreateRequest) (*User, error) {
	//user := NewUser(req)
	//
	//sigPrivateKey, err := crypto.KeyGen(&crypto.ECDSAP256KeyGenOpts{})
	//if err != nil {
	//	return nil, err
	//}
	//sigPrivateKeyPem, err := sigPrivateKey.Bytes()
	//if err != nil {
	//	return nil, err
	//}
	//tlsPrivateKey, err := crypto.KeyGen(&crypto.ECDSAP256KeyGenOpts{})
	//if err != nil {
	//	return nil, err
	//}
	//tlsPrivateKeyPem, err := tlsPrivateKey.Bytes()
	//if err != nil {
	//	return nil, err
	//}
	//
	//ak, err := crypto.KeyImport([]byte(req.TransactionPassword), &crypto.AES256KeyImportOpts{})
	//if err != nil {
	//	return nil, err
	//}
	//protectedSigPrivateKey, err := ak.Encrypt(sigPrivateKeyPem)
	//if err != nil {
	//	return nil, err
	//}
	//protectedTLSPrivateKey, err := ak.Encrypt(tlsPrivateKeyPem)
	//if err != nil {
	//	return nil, err
	//}
	//
	//user.ProtectedSigPrivateKey = base64.StdEncoding.EncodeToString(protectedSigPrivateKey)
	//user.ProtectedTLSPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)
	return nil, nil
}
