package tests

import (
	"context"
	"log"
	"testing"

	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TESTDBURI = "mongodb://localhost:27017"
	TESTDBNAME = "hotel-reservation-test"
)

type testdb struct {
	client *mongo.Client
	*data.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(TESTDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(TESTDBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := data.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		Store: &data.Store{
			User: data.NewMongoUserStore(client),
			Hotel: hotelStore,
			Room: data.NewMongoRoomStore(client, hotelStore),
			Booking: data.NewMongoBookingStore(client),
		},
	}
}
