package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Mail is a temporary stored message, which should be deleted upon consume.
// Client must not know whether it was consumed or not. TODO ??
// Server does not care about payload contents, file extension etc., it's client's work.
// Signature is intended to be used on the whole Mail struct, not just payload.
type Mail struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Sender    User               `json:"sender"`
	Receiver  User               `json:"receiver"`
	Signature Signature          `json:"signature"`
	Payload   []byte             `json:"payload"`
}
