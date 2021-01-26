package pyramid

import (
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

type Pyramid struct {
	Input  chan Action `json:"-"`
	Output chan Action `json:"-"`

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

func NewPyramidGame() *Pyramid {
	p := Pyramid{}
	p.RoomId = GenerateId(4)
	p.Players = make([]Player, 0)
	p.deck = card.NewDeck()
	p.boardCardIndex = 0
	p.attackState = false

	p.Input = make(chan Action, 25)
	p.Output = make(chan Action, 25)

	return &p
}

func (p *Pyramid) AddPlayer(player *Player) error {
	if p.Started {
		return ErrGameStarted
	}

	for _, v := range p.Players {
		if player.Name == v.Name {
			return ErrPlayerNameExists
		}
	}

	p.Players = append(p.Players, *player)
	return nil
}

func (p *Pyramid) Play() {
	p.Started = true
	go p.inputHandler()
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

func (p *Pyramid) dealCards() {
	p.Started = true
	p.deck = card.Shuffle(p.deck)
	p.board, p.deck = card.Deal(p.deck, 10)
	for k := range p.Players {
		p.Players[k].Hand, p.deck = card.Deal(p.deck, 4)
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

func (p *Pyramid) attack(attacker, attackee *Player, dmg int) {
	//p.Output <- attacker.Name + " ATTACKS " + attackee.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
}

func (p *Pyramid) acceptAttack(attacker, attackee *Player, dmg int) {
	//p.Output <- attackee.Name + " ACCEPTS ATTACK FROM " + attacker.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
	attackee.Sips += dmg
}

func (p *Pyramid) inputHandler() {
	//Forslag til struktur: Input = [roomid, acting player, action, message]
	for {
		event := <-p.Input
		var s []string
		if len(s) < 1 {
			fmt.Println("MESSAGE NOT UNDERSTOOD: " + event.ActionType)
			continue
		}
		switch s[2] {
		//Forslag til struktur: Input = [roomid, acting player, action = "ATTACK", target, dmg]
		case "ATTACK":
			if len(s) < 4 {
				fmt.Println("SHIT ATTACK MESSAGE: " + event.ActionType)
				continue
			}
			if !p.attackState {
				//p.Output <- "ATTACKING FAILED, NOT IN ATTACKING STATE"
			} else {
				var attackingPlayer, targetPlayer *Player
				//p.Output <- "ATTACKING " + s[1] + " " + s[2]
				for k := range p.Players {
					switch p.Players[k].Name {
					case s[1]:
						attackingPlayer = &p.Players[k]
					case s[3]:
						targetPlayer = &p.Players[k]
					}
				}
				k, _ := strconv.Atoi(s[4])
				p.attack(attackingPlayer, targetPlayer, k)
			}
		//Forslag til struktur: Input = [roomid, acting player, action = "ATTACK", target, dmg]
		case "ACCEPT_ATTACK":
			if len(s) < 3 {
				fmt.Println("SHIT ACCEPT_ATTACK MESSAGE: " + event.ActionType)
				continue
			}
			if !p.attackState {
				//p.Output <- "ACCEPT_ATTACK FAILED, NOT IN ATTACKING STATE"
			} else {
				var attackingPlayer, targetPlayer *Player
				//p.Output <- "ACCEPT_ATTACK " + s[1] + " " + s[2]
				for k := range p.Players {
					switch p.Players[k].Name {
					case s[1]:
						attackingPlayer = &p.Players[k]
					case s[3]:
						targetPlayer = &p.Players[k]
					}
				}
				k, _ := strconv.Atoi(s[4])
				p.acceptAttack(attackingPlayer, targetPlayer, k)

			}
		case "REJECT_ATTACK":
			if len(s) < 3 {
				fmt.Println("SHIT REJECT_ATTACK MESSAGE: " + event.ActionType)
				continue
			}
			if !p.attackState {
				//p.Output <- "REJECT_ATTACK FAILED, NOT IN ATTACKING STATE"
			} else {
				//p.Output <- "REJECT_ATTACK " + s[1] + " " + s[2]
			}
		}
	}
}
