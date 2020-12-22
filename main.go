package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var rooms []room

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hello frontend"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		var message message

		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		if message.Action == "new-room" {
			if len(rooms) > 0 {
				strconv.Itoa(len(rooms) + 1)
				rooms = append(rooms, room{strconv.Itoa(len(rooms) + 1), nil})
			} else {
				rooms = append(rooms, room{strconv.Itoa(1), nil})
			}
			fmt.Println(rooms)
		}
		if message.Action == "join-room" {
			for k, v := range rooms {
				if v.id == message.RoomID {
					rooms[k].players = append(v.players, player{message.PlayerName, nil})
				}
			}
			fmt.Println(rooms)
		}
	}
}

type message struct {
	Action     string `json:"action"`
	PlayerName string `json:"name"`
	RoomID     string `json:"id"`
}

type room struct {
	id      string
	players []player
}

type player struct {
	playerName string
	hand       []Card
}

// card:

// Card holds the card suits and types in the deck
type Card struct {
	Type string
	Suit string
}

// Deck holds the cards in the deck to be shuffled
type Deck []Card

// New creates a deck of cards to be used
func NewDeck() (deck Deck) {

	// Valid types include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	types := []string{"Two", "Three", "Four", "Five", "Six", "Seven",
		"Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace"}

	// Valid suits include Heart, Diamond, Club & Spade
	suits := []string{"Heart", "Diamond", "Club", "Spade"}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(types); i++ {
		for n := 0; n < len(suits); n++ {
			card := Card{
				Type: types[i],
				Suit: suits[n],
			}
			deck = append(deck, card)
		}
	}
	return
}

// Shuffle the deck
func Shuffle(d Deck) Deck {
	for i := 1; i < len(d); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d[r], d[i] = d[i], d[r]
		}
	}
	return d
}

// Deal a specified amount of cards
func Deal(d Deck, n int) Deck {
	for k1, v1 := range rooms {
		for k2, _ := range v1.players {
			for i := 0; i < n; i++ {
				rooms[k1].players[k2].hand = append(rooms[k1].players[k2].hand, d[i])
				d = d[1:]
			}
		}
	}
	return d
}

// Debug helps debugging the deck of cards
func Debug(d Deck) {
	if os.Getenv("DEBUG") != "" {
		for i := 0; i < len(d); i++ {
			fmt.Printf("Card #%d is a %s of %ss\n", i+1, d[i].Type, d[i].Suit)
		}
	}
}

// Seed our randomness with the current time
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//fmt.Println("Hello World")
	//setupRoutes()
	//log.Fatal(http.ListenAndServe(":8081", nil))
	rooms = append(rooms, room{"0", nil})
	rooms[0].players = append(rooms[0].players, player{"Noer", nil})
	rooms[0].players = append(rooms[0].players, player{"Dennis", nil})

	deck := NewDeck()
	Shuffle(deck)
	deck = Deal(deck, 4)
	fmt.Println(rooms)
	fmt.Println(len(deck))
}
