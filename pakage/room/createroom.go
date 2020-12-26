package room

import (
	"log"
	"math/rand"
)

func randSeq(n int) string {
	var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := make([]rune, n)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

func CreateRoom() string {
	for {
		roomID := randSeq(4)
		element, exist := Rooms[roomID]
		if !exist {
			element = Room{}
			element.Players = make(map[string]*Player)
			Rooms[roomID] = element
			log.Println("Room", roomID, "created")
			return roomID
		} else {
			log.Println("Room", roomID, "already exists")
		}
	}
}
