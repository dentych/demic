package room

import "log"

//Remove Player from Room
func RemovePlayer(roomID string, playerName string) {
	delete(Rooms[roomID].Players, playerName)
	log.Println("Room", roomID, ":", playerName, "Left the room")
}
