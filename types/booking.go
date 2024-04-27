package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	UserID			   primitive.ObjectID 	`bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID			   primitive.ObjectID 	`bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons		   int 					`bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate 		   time.Time			`bson:"fromDate" json:"fromDate"`
	UntilDate 		   time.Time			`bson:"untilDate" json:"untilDate"`
}