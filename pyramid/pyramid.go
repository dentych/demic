package pyramid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/dentych/demic/card"
	"gitlab.com/dentych/demic/models"
	"gitlab.com/dentych/demic/util"
	"log"
	"sort"
	"strconv"
	"time"
)

var (
	ErrGameStarted      = fmt.Errorf("game is Started")
	ErrGameNotStarted   = fmt.Errorf("game is not Started")
	ErrPlayerNameExists = fmt.Errorf("player name already exists")
	ErrNoMoreCards      = fmt.Errorf("no more cards to turn")
	ErrAttackState      = fmt.Errorf("wrong state of the game")
)

var Rooms map[string]*Pyramid

type Pyramid struct {
	Input chan models.IncomingMessage `json:"-"`

	RoomId         string   `json:"room_id"`
	Started        bool     `json:"started"`
	Players        []Player `json:"players"`
	Attacks        Attacks  `json:"attacks"`
	boardCardIndex int
	board          []card.Card
	deck           []card.Card
	cont           bool
	attackState    bool
	gameEnd        bool
}

func init() {
	Rooms = make(map[string]*Pyramid)
}

func Create(clientID string) (roomID string, input chan<- models.IncomingMessage, output <-chan models.OutgoingMessage, err error) {
	p := Pyramid{}
	p.RoomId = GenerateId(4)
	p.Players = make([]Player, 0)
	p.deck = card.NewDeck()
	p.boardCardIndex = 0
	p.attackState = false

	p.Input = make(chan models.IncomingMessage, 25)

	player := NewPlayer(clientID, "HOST")
	p.Players = append(p.Players, *player)

	Rooms[p.RoomId] = &p

	go p.play()
	go p.inputHandler()

	return p.RoomId, p.Input, player.Output, nil
}

func Join(roomID, clientID, playerName string) (chan<- models.IncomingMessage, <-chan models.OutgoingMessage, error) {
	room, ok := Rooms[roomID]
	if !ok {
		return nil, nil, fmt.Errorf("room with ID '%s' not found", roomID)
	}

	player, err := room.joinPlayer(clientID, playerName)
	if err != nil {
		return nil, nil, err
	}

	return room.Input, player.Output, nil
}

func (p *Pyramid) joinPlayer(clientID, playerName string) (player Player, err error) {
	for index := range p.Players {
		if clientID == p.Players[index].ClientID {
			p.Players[index] = *NewPlayer(clientID, playerName)
			player = p.Players[index]
			return
		}
	}

	player = *NewPlayer(clientID, playerName)
	p.Players = append(p.Players, player)
	if len(p.Players) > 1 {
		player.Output <- actionToOutgoing(PayloadAction{
			ActionType: ActionHost,
			Target:     p.Players[1].Name,
		})
	}
	p.output(PayloadAction{
		ActionType: ActionPlayerJoined,
		Target:     player.Name,
	})
	return
}

func (p *Pyramid) removePlayer(player *Player) error {
	for i := range p.Players {
		if player.Name == p.Players[i].Name {
			p.Players[len(p.Players)-1], p.Players[i] = p.Players[i], p.Players[len(p.Players)-1]
			p.Players = p.Players[:len(p.Players)-1]
			p.output(PayloadAction{ActionType: ActionPlayerQuit, Target: player.Name})
			return nil
		}
	}

	return fmt.Errorf("player '%s' wasn't found", player.Name)
}

func (p *Pyramid) play() {
	for !p.Started {
		time.Sleep(500 * time.Millisecond)
	}
	p.dealCards()
	p.waitForContinue()
	p.output(PayloadAction{ActionType: ActionStartGame})

	// main loop
	for p.boardCardIndex != len(p.board) {
		err := p.turnNextCard()
		if err != nil {
			log.Panic(err)
		}
		p.output(PayloadAction{ActionType: ActionNewRound})
		p.waitForContinue()
	}

	// Game end
	p.gameEnd = true
	p.output(PayloadAction{ActionType: ActionGameEnd})

	select {}
}

func (p *Pyramid) updateAttackState() {
	switch p.Attacks.Len() {
	case 0:
		p.attackState = false
	case 1:
		p.attackState = true
	default:
		p.attackState = true
		return
	}

	p.Players[1].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionAttackState,
		Origin:     p.Players[0].Name,
		Target:     strconv.FormatBool(p.attackState),
	})

}

func actionToOutgoing(action PayloadAction) models.OutgoingMessage {
	return models.OutgoingMessage{
		ActionType: action.ActionType,
		Payload:    action,
	}
}

func (p *Pyramid) continueGame() {
	p.cont = true
}

func (p *Pyramid) output(action PayloadAction, excluded ...string) {
	outgoingMessage := actionToOutgoing(action)
	for _, player := range p.Players {
		if util.ArrayContains(excluded, player.Name) {
			continue
		}
		player.Output <- outgoingMessage
	}
}

func (p *Pyramid) dealCards() {
	p.Started = true
	card.Shuffle(&p.deck)
	p.board = card.Deal(&p.deck, 15)
	for k := range p.Players {
		if p.Players[k].Name == "HOST" {
			continue
		}
		p.Players[k].Hand = card.Deal(&p.deck, 4)
		sort.SliceStable(p.Players[k].Hand, func(i, j int) bool {
			var a, b int
			aRank, bRank := p.Players[k].Hand[i].Rank, p.Players[k].Hand[j].Rank
			if len(aRank) == 2 {
				a = 10
			} else {
				if aRank[0] == 'A' {
					a = 0
				} else {
					a = int(aRank[0] - 48)
				}
			}

			if len(bRank) == 2 {
				b = 10
			} else {
				if bRank[0] == 'A' {
					b = 0
				} else {
					b = int(bRank[0] - 48)
				}
			}
			return a < b
		})
		var cardByte bytes.Buffer
		for _, v := range p.Players[k].Hand {
			cardByte.WriteString(v.String())
			cardByte.WriteString(",")
		}
		cardStr := cardByte.String()
		p.Players[k].Output <- actionToOutgoing(PayloadAction{
			ActionType: ActionDealHand,
			Origin:     p.Players[k].Name,
			Target:     cardStr[:len(cardStr)-1],
		})
	}
}

func (p *Pyramid) turnNextCard() error {
	var cardByte bytes.Buffer
	if !p.Started {
		return ErrGameNotStarted
	}
	if p.boardCardIndex >= len(p.board) {
		return ErrNoMoreCards
	}

	c := p.board[p.boardCardIndex]
	p.boardCardIndex++

	cardByte.WriteString(c.Rank)
	cardByte.WriteRune(c.Suit)
	cardStr := cardByte.String()

	p.Players[0].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionDealHand,
		Origin:     p.Players[0].Name,
		Target:     cardStr,
	})
	return nil
}

func (p *Pyramid) waitForContinue() {
	for {
		if p.cont {
			p.cont = false
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (p *Pyramid) attack(event PayloadAction) error {
	var originIdx, targetIdx int
	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Attacks.Add(Attack{
		Attacker: p.Players[originIdx],
		Target:   p.Players[targetIdx],
	})
	p.updateAttackState()

	msg := actionToOutgoing(PayloadAction{
		ActionType: ActionAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	})
	p.Players[0].Output <- msg
	p.Players[targetIdx].Output <- msg
	return nil
}

func (p *Pyramid) rejectAttack(event PayloadAction) error {
	var originIdx, targetIdx int
	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Attacks.Remove(Attack{Attacker: p.Players[targetIdx], Target: p.Players[originIdx]})
	p.updateAttackState()

	msg := actionToOutgoing(PayloadAction{
		ActionType: ActionRejectAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	})
	p.Players[0].Output <- msg
	p.Players[targetIdx].Output <- msg
	return nil
}

func (p *Pyramid) acceptAttack(event PayloadAction) error {
	var originIdx, targetIdx int
	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Attacks.Remove(Attack{Attacker: p.Players[targetIdx], Target: p.Players[originIdx]})
	p.updateAttackState()

	msg := actionToOutgoing(PayloadAction{
		ActionType: ActionAcceptAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	})

	p.Players[0].Output <- msg
	p.Players[targetIdx].Output <- msg
	return nil
}

func (p *Pyramid) pickCard(event PayloadAction) error {
	var originIdx int

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		}
	}

	handIdx, err := strconv.Atoi(event.Target)
	if err != nil {
		log.Println("Error converting player pick card target to int:", err)
		return err
	}

	chosenCard := p.Players[originIdx].Hand[handIdx]
	p.deck = append(p.deck, chosenCard)
	card.Shuffle(&p.deck)

	p.Players[0].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionPickCard,
		Origin:     p.Players[originIdx].Name,
		Target:     chosenCard.String(),
	})
	newCard := card.Deal(&p.deck, 1)[0]
	p.Players[originIdx].Hand[handIdx] = newCard
	var hand bytes.Buffer
	for _, v := range p.Players[originIdx].Hand {
		hand.WriteString(v.String())
		hand.WriteString(",")
	}
	cardStr := hand.String()
	p.Players[originIdx].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionDealHand,
		Target:     cardStr[:len(cardStr)-1],
	})
	return nil
}

func (p *Pyramid) newCard(event PayloadAction) {
	var originIdx, handIdx int
	var cardByte bytes.Buffer
	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			handIdx = k
		}
	}
	newCard := card.Deal(&p.deck, 1)[0]
	p.Players[originIdx].Hand[handIdx] = newCard

	for _, v := range p.Players[originIdx].Hand {
		cardByte.WriteString(v.Rank)
		cardByte.WriteRune(v.Suit)
		cardByte.WriteString(",")
	}

	cardStr := cardByte.String()[:len(cardByte.String())-1]

	p.Players[originIdx].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionDealHand,
		Origin:     p.Players[originIdx].Name,
		Target:     cardStr,
	})
}

func (p *Pyramid) inputHandler() {
	var err error
	for {
		event := <-p.Input
		var payload PayloadAction
		err = json.Unmarshal(event.Payload, &payload)
		if err != nil {
			log.Println("failed to unmarshal payload in inputhandler:", err)
			return
		}
		switch event.ActionType {
		case ActionStartGame:
			p.Started = true
		case ActionAttack:
			p.attack(payload)
		case ActionAcceptAttack:
			p.acceptAttack(payload)
		case ActionRejectAttack:
			p.rejectAttack(payload)
		case ActionPickCard:
			p.pickCard(payload)
		//case ActionNewCard:
		//	p.newCard(event)
		case ActionContinue:
			p.continueGame()
		case ActionShowCard:
			p.showCard(payload)
		}
	}
}

func (p *Pyramid) showCard(event PayloadAction) error {
	if !p.gameEnd {
		return fmt.Errorf("game not ended yet")
	}

	var player Player
	for _, v := range p.Players {
		if v.Name == event.Origin {
			player = v
			break
		}
	}

	cardIndex, err := strconv.Atoi(event.Target)
	if err != nil {
		return err
	}
	p.Players[0].Output <- actionToOutgoing(PayloadAction{
		ActionType: ActionShowCard,
		Origin:     event.Origin,
		Target:     player.Hand[cardIndex].String(),
	})

	return nil
}
