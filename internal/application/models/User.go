package models

type User struct {
	Name      string    `json:"name"`
	PublicKey PublicKey `json:"public_key"`
}
