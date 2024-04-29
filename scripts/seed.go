package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kortbyhotel/reservation/api"
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
	userStore data.UserStore
	ctx = context.Background()
)

func seedUser(fname,lname, email string, isAdmin bool) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "test1234",
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	usersTokens, err := api.CreateTokenFromUser(user)
	if err != nil {
		fmt.Printf("cannot create token")
	}
	fmt.Printf("%s -> %s\n", user.Email, usersTokens)
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Romms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Size: "Small",
			Price: 230.99,
		},
		{
			Size: "normal",
			Price: 49.99,
		},
		{
			Size: "kingsize",
			Price: 149.99,
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
	userStore = data.NewMongoUserStore(client)
}


func main() {

	seedHotel("Kortby INN", "Italy", 4)
	seedHotel("Marriott", "Spain", 3)
	seedHotel("Hyatt", "US", 5)
	seedUser("john", "Doe", "test@test.com", false)
	seedUser("admin", "Admin", "admin@test.com", true)
}