package room

func Deleteplayer(playerName string, roomID string) {
	delete(Rooms[roomID].Players, playerName)
}
