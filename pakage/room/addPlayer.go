package room

import "log"

//add Player to Room
func AddPlayerToRoom(roomID string, playerName string) {
	element, exist := Rooms[roomID]
	if exist {
		element.Players[playerName] = &Player{}
		Rooms[roomID] = element
		log.Println("Room", roomID, ":", playerName, "Joined the room")
	} else {
		log.Println("Room", roomID, ": RoomID not found")
	}
}
