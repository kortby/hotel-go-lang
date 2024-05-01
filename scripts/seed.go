package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/config"
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
	bookingStore data.BookingStore
	ctx = context.Background()
)

func seedUser(fname,lname, email string, isAdmin bool) *types.User {
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
	insertedUser, err := userStore.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	usersTokens, err := api.CreateTokenFromUser(user)
	if err != nil {
		fmt.Printf("cannot create token")
	}
	fmt.Printf("%s -> %s\n", user.Email, usersTokens)
	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size: size,
		Seeside: ss,
		Price: price,
		HotelID: hotelID,
	}
	insertesRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertesRoom
}

func seedHotel(name string, location string, rating int) (*types.Hotel) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Romms: []primitive.ObjectID{},
		Rating: rating,
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func seedBooking(userID primitive.ObjectID, roomID primitive.ObjectID, from time.Time, unitl time.Time) {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		UntilDate: unitl,
	}
	res, err := bookingStore.InsertBooking(context.Background(), booking); 
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking id -> ", res.ID)
}


func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.New().DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(config.New().TESTDBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = data.NewMongoHotelStore(client)
	roomStore = data.NewMongoRoomStore(client, hotelStore)
	userStore = data.NewMongoUserStore(client)
	bookingStore = data.NewMongoBookingStore(client)
}


func main() {

	guest := seedUser("john", "Doe", "test@test.com", false)
	seedUser("admin", "Admin", "admin@test.com", true)
	hotel := seedHotel("Hyatt", "US", 5)
	room := seedRoom("Queen Room", false, 109, hotel.ID)
	seedBooking(guest.ID, room.ID, time.Now(), time.Now().AddDate(0,0,2))

	for i := 0; i < 500; i++ {
		hotel := seedHotel(faker.Name(), faker.Century(), rand.Intn(5))
		
		for i := 0; i < 50; i++ {
			price := rand.Intn(235)
			seedRoom(faker.Name(), false, float64(price), hotel.ID)
		}
	}
}