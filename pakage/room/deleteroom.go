package room

func DeleteRoom(roomID string) {
	delete(Rooms, roomID)
}
