package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12 // Good practice to make this configurable based on environment
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName          string             `bson:"firstName" json:"firstName"`
	LastName           string             `bson:"lastName" json:"lastName"`
	Email              string             `bson:"email" json:"email"`
	EncryptedPassword  string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                 primitive.NewObjectID(),
		FirstName:          params.FirstName,
		LastName:           params.LastName,
		Email:              params.Email,
		EncryptedPassword:  string(encpw),
	}, nil
}
