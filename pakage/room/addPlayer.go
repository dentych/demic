package room

import "log"

func AddPlayerToRoom(playerName string, roomID string) {
	element, exist := Rooms[roomID]
	if exist {
		element.Players[playerName] = &Player{}
		Rooms[roomID] = element
		log.Println(playerName, "joined room ID:", roomID)
	} else {
		log.Println("Room", roomID, "pla", playerName, "not found")
	}
}
