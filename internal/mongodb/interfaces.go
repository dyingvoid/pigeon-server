package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type CrudRepository[T any] interface {
	Add(T) error
	Get(primitive.ObjectID) (T, error)
	Update(primitive.ObjectID, T) error
	Delete(primitive.ObjectID) error
	List() ([]T, error)
}
