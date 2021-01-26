package pyramid

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
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

func TestPyramidGame(t *testing.T) {
	p := NewPyramidGame()

	p1 := NewPlayer("Dennis")
	p2 := NewPlayer("Noer")
	err := p.AddPlayer(p1)
	if err != nil {
		log.Panic(err)
	}
	err = p.AddPlayer(p2)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for {
			fmt.Println("OUTPUT: " + <-p.Output)
		}
	}()
	go p.Play()

	timer := time.NewTimer(1 * time.Second)
	<-timer.C
	p.Input <- p.RoomId + " " + p1.Name + " ATTACK " + p2.Name + " 4"
	p.Input <- p.RoomId + " " + p1.Name + " ACCEPT_ATTACK " + p2.Name + " 4"
	fmt.Println(p.Players)
	for {
	}
}
