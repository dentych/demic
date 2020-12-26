package room

import "log"

//Delete Room from Rooms
func DeleteRoom(roomID string) {
	delete(Rooms, roomID)
	log.Println("Room", roomID, ": Deleted")
}
