package pyramid

import (
	"bytes"
	"fmt"
	"gitlab.com/dentych/demic/card"
	"log"
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

var PyramidRooms map[string]*Pyramid

type Pyramid struct {
	Input       chan Action  `json:"-"`
	PlayerJoin  chan *Player `json:"-"`
	PlayerLeave chan *Player `json:"-"`

	RoomId         string   `json:"room_id"`
	Started        bool     `json:"started"`
	Players        []Player `json:"players"`
	attacks        []Attack
	boardCardIndex int
	board          []card.Card
	deck           []card.Card
	cont           bool
	attackState    bool
}

type Attack struct {
	Attacker Player
	Target   Player
	Accepted bool
	Dmg      int
}

func init() {
	PyramidRooms = make(map[string]*Pyramid)
}

func NewPyramidGame() *Pyramid {
	p := Pyramid{}
	p.RoomId = GenerateId(4)
	p.Players = make([]Player, 0)
	p.deck = card.NewDeck()
	p.boardCardIndex = 0
	p.attackState = false

	p.Input = make(chan Action, 25)
	p.PlayerJoin = make(chan *Player)
	p.PlayerLeave = make(chan *Player)

	go p.play()
	go p.InputHandler()
	go p.playerJoinHandler()
	go p.playerLeaveHandler()

	return &p
}

func (p *Pyramid) addPlayer(player *Player) error {
	if p.Started {
		return ErrGameStarted
	}

	for _, v := range p.Players {
		if player.Name == v.Name {
			return ErrPlayerNameExists
		}
	}

	for _, v := range p.Players {
		if v.Name != "HOST" {
			player.Output <- Action{ActionType: ActionPlayerJoin, Target: v.Name}
		}
	}
	p.Players = append(p.Players, *player)
	p.output(Action{
		ActionType: ActionPlayerJoin,
		Target:     player.Name,
	})
	return nil
}

func (p *Pyramid) removePlayer(player *Player) error {
	for i := range p.Players {
		if player.Name == p.Players[i].Name {
			p.Players[len(p.Players)-1], p.Players[i] = p.Players[i], p.Players[len(p.Players)-1]
			p.Players = p.Players[:len(p.Players)-1]
			p.output(Action{ActionType: ActionPlayerQuit, Target: player.Name})
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
	p.output(Action{ActionType: ActionStartGame})

	// main loop
	for p.boardCardIndex != len(p.board) {
		err := p.turnNextCard()
		if err != nil {
			log.Panic(err)
		}
		p.setAttackState(true)
		p.waitForContinue()
		p.setAttackState(false)
	}
}

func (p *Pyramid) setAttackState(value bool) {
	p.attackState = value
	p.output(Action{
		ActionType: ActionAttackState,
		Target:     strconv.FormatBool(value),
	})

}

func (p *Pyramid) Continue() {
	p.cont = true
}

func (p *Pyramid) output(action Action) {
	for _, player := range p.Players {
		player.Output <- action
	}
}

func (p *Pyramid) dealCards() {
	var cardByte bytes.Buffer
	p.Started = true
	card.Shuffle(&p.deck)
	p.board = card.Deal(&p.deck, 15)
	for k := range p.Players {
		if p.Players[k].Name == "HOST" {
			continue
		}
		p.Players[k].Hand = card.Deal(&p.deck, 4)
		for _, v := range p.Players[k].Hand {
			cardByte.WriteString(v.Rank)
			cardByte.WriteRune(v.Suit)
			cardByte.WriteString(",")
		}
		cardStr := cardByte.String()[:len(cardByte.String())-1]
		p.Players[k].Output <- Action{
			ActionType: ActionDealHand,
			Origin:     p.Players[k].Name,
			Target:     cardStr,
		}

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

	p.Players[0].Output <- Action{
		ActionType: ActionDealHand,
		Origin:     p.Players[0].Name,
		Target:     cardStr,
	}
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

func (p *Pyramid) attack(event Action) error {
	var originIdx, targetIdx int
	if !p.attackState {
		return ErrAttackState
	}

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) rejectAttack(event Action) error {
	var originIdx, targetIdx int
	if !p.attackState {
		return ErrAttackState
	}

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionRejectAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) acceptAttack(event Action) error {
	var originIdx, targetIdx int
	if !p.attackState {
		return ErrAttackState
	}

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		case event.Target:
			targetIdx = k
		}
	}

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionAcceptAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) pickCard(event Action) error {
	var originIdx, handIdx int
	var cardByte bytes.Buffer
	if !p.attackState {
		return ErrAttackState
	}

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		}
	}

	cardByte.WriteString(p.Players[originIdx].Hand[handIdx].Rank)
	cardByte.WriteRune(p.Players[originIdx].Hand[handIdx].Suit)
	cardStr := cardByte.String()

	p.output(Action{
		ActionType: ActionPickCard,
		Origin:     p.Players[originIdx].Name,
		Target:     cardStr,
	})
	p.Players[originIdx].Hand[handIdx] = card.Deal(&p.deck, 1)[0]
	return nil
}

func (p *Pyramid) turnCard(event Action) {
	//p.Output <- attackee.Name + " ACCEPTS ATTACK FROM " + attacker.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
}

func (p *Pyramid) InputHandler() {
	for {
		event := <-p.Input
		switch event.ActionType {
		case ActionStartGame:
			p.Started = true
		case ActionAttack:
			p.attack(event)
		case ActionAcceptAttack:
			p.rejectAttack(event)
		case ActionRejectAttack:
			p.acceptAttack(event)
		case ActionPickCard:
			p.pickCard(event)
		case ActionContinue:
			p.Continue()
		}
	}
}

func (p *Pyramid) playerJoinHandler() {
	for {
		pl := <-p.PlayerJoin
		err := p.addPlayer(pl)
		if err != nil {
			log.Println("Failed to add player "+pl.Name, err)
		}
		if len(p.Players) > 1 {
			pl.Output <- Action{
				ActionType: ActionHost,
				Origin:     "",
				Target:     p.Players[1].Name,
			}
		}
	}
}

func (p *Pyramid) playerLeaveHandler() {
	for {
		pl := <-p.PlayerLeave
		err := p.removePlayer(pl)
		if err != nil {
			log.Println("Failed to remove player: "+pl.Name, err)
		}
	}
}
