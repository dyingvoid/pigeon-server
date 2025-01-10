package repositories

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailRepository struct {
	collection *mongo.Collection
}

type MailMetadata struct {
	Ids   []primitive.ObjectID `json:"ids"`
	Count int64                `json:"count"`
	Space int64                `json:"space"`
}

/*func (r *MailRepository) SendMail(mail models.Mail) error {

}

func (r *MailRepository) ReceiveAllMails(receiver models.User) ([]models.Mail, error) {

}

// TODO authorization ( great :) )
func (r *MailRepository) ReceiveMail(id primitive.ObjectID) error {

}

func (r *MailRepository) PeakMails(mail models.Mail) (MailMetadata, error) {

}*/
