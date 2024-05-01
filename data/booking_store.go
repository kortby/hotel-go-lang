package data

import (
	"context"

	"github.com/kortbyhotel/reservation/config"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type BookingStore interface {
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, primitive.ObjectID) (*types.Booking, error)
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	UpdateBooking(context.Context, primitive.ObjectID, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll: client.Database(config.New().DBNAME).Collection("bookings"),
	}
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func  (s *MongoBookingStore) GetBookingByID(ctx context.Context, id primitive.ObjectID) (*types.Booking, error) {
	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&booking); err != nil {
		return nil, err
	}
	
	return &booking, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id primitive.ObjectID, update bson.M) error {
    result, err := s.coll.UpdateByID(ctx, id, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return mongo.ErrNoDocuments
    }

    return nil  // Success
}