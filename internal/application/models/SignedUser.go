package models

type SignedUser struct {
	User      User      `json:"user"`
	Signature Signature `json:"signature"`
}
