package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	roomStore data.RoomStore
	hotelStore data.HotelStore
	ctx = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Romms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Type: types.SeaSideRoomType,
			BasePrice: 230.99,
		},
		{
			Type: types.SingleRoomType,
			BasePrice: 49.99,
		},
		{
			Type: types.DeluxRoomType,
			BasePrice: 149.99,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hotel name ", insertedHotel)

	

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}


func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(data.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(data.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = data.NewMongoHotelStore(client)
	roomStore = data.NewMongoRoomStore(client, hotelStore)
}


func main() {

	seedHotel("Kortby INN", "Italy", 4)
	seedHotel("Marriott", "Spain", 3)
	seedHotel("Hyatt", "US", 5)

}