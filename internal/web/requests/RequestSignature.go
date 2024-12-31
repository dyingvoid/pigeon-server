package requests

type RequestSignature struct {
	IssuedNonce []byte `json:"issued_nonce"`
	Signature   []byte `json:"signature"`
	PublicKey   []byte `json:"public_key_der"`
}
