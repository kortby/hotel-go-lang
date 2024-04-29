package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	minNameLen = 2
	minPasswordLen = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName          string             `bson:"firstName" json:"firstName"`
	LastName           string             `bson:"lastName" json:"lastName"`
	Email              string             `bson:"email" json:"email"`
	EncryptedPassword  string             `bson:"encryptedPassword" json:"-"`
	IsAdmin   		   bool	 			  `bson:"isAdmin" json:"isAdmin"`
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

// emailRegex is a regular expression for validating email addresses.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// validateEmail checks if the email provided matches the RFC 5322 standard.
func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minNameLen {
		errors["firstName"] = fmt.Sprintf("first name length should be at leas %d characters", minNameLen)
	}
	if len(params.LastName) < minNameLen {
		errors["lastName"] = fmt.Sprintf("last name length should be at leas %d characters", minNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at leas %d characters", minPasswordLen)
	}
	if !validateEmail(params.Email) {
		errors["email"] = fmt.Sprintln("please return a valid email")
	}
	return errors
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
