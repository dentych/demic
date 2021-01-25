package pyramid

import (
	"strings"
	"testing"
)

func TestPyramidIdValid(t *testing.T) {
	pyramid := NewPyramidGame()

	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for _, v := range pyramid.RoomId {
		if !strings.ContainsRune(validChars, v) {
			t.Errorf("Invalid pyramid room ID char: %c\n", v)
		}
	}
}

func TestPyramid_AddPlayer(t *testing.T) {
	pyramid := NewPyramidGame()

	player := NewPlayer("Dennis")
	_ = pyramid.AddPlayer(player)

	result := false
	for _, v := range pyramid.Players {
		if v.Name == player.Name {
			result = true
			break
		}
	}

	if result != true {
		t.Error("Pyramid game doesn't include added player")
	}
}

func TestPyramid_DealCards(t *testing.T) {

}