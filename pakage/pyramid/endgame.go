package pyramid

import (
	"gitlab.com/dentych/demic/pakage/room"
	"log"
)

//ends game in room.Room
func EndGame(roomID string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Deck = nil
		element.Board = nil
		for playerName, _ := range element.Players {
			element.Players[playerName].Hand = nil
		}
		room.Rooms[roomID] = element
	}
	log.Println("Room", roomID, ": Game ended")
}
