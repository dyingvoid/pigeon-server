package requests

type SignedRequest interface {
	GetSignature() RequestSignature
}
