package requests

type NewUserRequest struct {
	Username    string `json:"username"`
	PublicKey   []byte `json:"public_key_der"`
	IssuedNonce []byte `json:"issued_nonce"`
	Signature   []byte `json:"signature"`
}

func (r NewUserRequest) GetSignature() RequestSignature {
	return RequestSignature{
		IssuedNonce: r.IssuedNonce,
		Signature:   r.Signature,
	}
}
