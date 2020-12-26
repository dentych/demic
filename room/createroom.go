package room

import (
	"log"
	"math/rand"
	"time"
)

//adds time to generate unique randomness
func init() {
	rand.Seed(time.Now().UnixNano())
}

//generates a random roomID
func generateRoomID(n int) string {
	var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := make([]rune, n)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

//Creates a Room with a unique roomID in Rooms
func CreateRoom() string {
	for {
		roomID := generateRoomID(4)
		element, exist := Rooms[roomID]
		if !exist {
			element = Room{}
			element.Players = make(map[string]*Player)
			Rooms[roomID] = element
			log.Println("Room", roomID, ": Created")
			return roomID
		}
	}
}
