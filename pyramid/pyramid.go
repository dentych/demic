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
		player.Output <- Action{ActionType: ActionPlayerJoined, Target: v.Name}
	}
	p.Players = append(p.Players, *player)
	p.output(Action{
		ActionType: ActionPlayerJoined,
		Target:     player.Name,
	})
	return nil
}

func (p *Pyramid) removePlayer(player *Player) error {
	for i := range p.Players {
		if player.Name == p.Players[i].Name {
			p.Players[len(p.Players)-1], p.Players[i] = p.Players[i], p.Players[len(p.Players)-1]
			p.Players = p.Players[:len(p.Players)-1]
			p.output(Action{ActionType: ActionPlayerLeft, Target: player.Name})
			return nil
		}
	}

	return fmt.Errorf("player '%s' wasn't found", player.Name)
}

func (p *Pyramid) play() {
	p.Started = true
	p.dealCards()

	p.waitForContinue()

	for p.boardCardIndex != len(p.board) {
		_, err := p.turnNextCard()
		if err != nil {
			log.Panic(err)
		}
		//p.Output <- "CARD " + string(c.Suit) + c.Rank

		//p.Output <- "ATTACK BEGIN"
		p.attackState = true
		p.waitForContinue()
		p.attackState = false
		//p.Output <- "ATTACK STOP"
	}
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
	var handStr bytes.Buffer
	p.Started = true
	p.deck = card.Shuffle(p.deck)
	p.board, p.deck = card.Deal(p.deck, 10)
	for k := range p.Players {
		p.Players[k].Hand, p.deck = card.Deal(p.deck, 4)
		for _, v := range p.Players[k].Hand {
			handStr.WriteString(v.Rank)
			handStr.WriteRune(v.Suit)
			handStr.WriteString(",")
		}
		p.Players[k].Output <- Action{
			ActionType: ActionDealHand,
			Origin:     p.Players[k].Name,
			Target:     handStr.String()[:len(handStr.String())-1],
		}

	}
}

func (p *Pyramid) turnNextCard() (*card.Card, error) {
	if !p.Started {
		return nil, ErrGameNotStarted
	}

	if p.boardCardIndex < len(p.board) {
		c := p.board[p.boardCardIndex]
		p.boardCardIndex++
		return &c, nil
	} else {
		return nil, ErrNoMoreCards
	}
}

func (p *Pyramid) waitForContinue() {
	for {
		if p.cont {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (p *Pyramid) attack(attacker, attackee *Player) {
	p.output(Action{
		ActionType: ActionAttack,
		Origin:     attacker.Name,
		Target:     attackee.Name,
	})
}

func (p *Pyramid) rejectAttack(attacker, attackee *Player) {
	p.output(Action{
		ActionType: ActionRejectAttack,
		Origin:     attacker.Name,
		Target:     attackee.Name,
	})
}

func (p *Pyramid) acceptAttack(attacker, attackee *Player) {
	p.output(Action{
		ActionType: ActionAcceptAttack,
		Origin:     attacker.Name,
		Target:     attackee.Name,
	})
}

func (p *Pyramid) pickCard(player *Player, idx int) {
	p.output(Action{
		ActionType: ActionAcceptAttack,
		Origin:     player.Name,
		Target:     strconv.Itoa(idx),
	})
}

func (p *Pyramid) turnCard(origin Player, cardIdx int) {
	//p.Output <- attackee.Name + " ACCEPTS ATTACK FROM " + attacker.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
}

func (p *Pyramid) InputHandler() {
	var origin, target *Player

	//Forslag til struktur: Input = [roomid, acting player, action, message]
	for {
		event := <-p.Input
		for _, player := range p.Players {
			switch player.Name {
			case event.Origin:
				origin = &player
			case event.Target:
				target = &player
			}
		}
		switch event.ActionType {
		case ActionStartGame:
			p.play()
		case ActionAttack:
			p.attack(origin, target)
		case ActionAcceptAttack:
			p.rejectAttack(origin, target)
		case ActionRejectAttack:
			p.acceptAttack(origin, target)
		case ActionPickCard:
			target, err := strconv.Atoi(event.Target)
			if err != nil {
				log.Println("Cannot convert string to int: "+event.Target, err)
			}
			p.pickCard(origin, target)
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
