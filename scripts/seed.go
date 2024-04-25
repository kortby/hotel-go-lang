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

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(data.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := data.NewMongoHotelStore(client)
	roomStore := data.NewMongoRoomStore(client, hotelStore)

	hotel := types.Hotel{
		Name: "Marriott",
		Location: "Italy",
		Romms: []primitive.ObjectID{},
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
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Room ", insertedRoom)
	}
}