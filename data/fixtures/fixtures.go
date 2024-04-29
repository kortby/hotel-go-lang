package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *data.Store, fname, lname string, isAdmin bool) (*types.User, string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: fmt.Sprintf("%s@%s.com",fname, lname),
		FirstName: fname,
		LastName: lname,
		Password: "test1234",
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	usersTokens, err := api.CreateTokenFromUser(user)
	if err != nil {
		fmt.Printf("cannot create token")
	}
	return insertedUser, usersTokens
}

func AddHotel(store *data.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name: name,
		Location: loc,
		Romms: roomIDs,
		Rating: rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(store *data.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size: size,
		Seeside: ss,
		Price: price,
		HotelID: hotelID,
	}
	insertesRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertesRoom
}

func AddBooking(store *data.Store, userID primitive.ObjectID, roomID primitive.ObjectID, from time.Time, unitl time.Time) *types.Booking {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		UntilDate: unitl,
	}
	res, err := store.Booking.InsertBooking(context.Background(), booking); 
	if err != nil {
		log.Fatal(err)
	}
	return res
}
