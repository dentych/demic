package pyramid

import (
	"gitlab/dentych/demic/pakage/room"
)

func EndGame(roomID string) {
	delete(room.Rooms, roomID)
}
