package pyramid

import (
	"fmt"
	"gitlab.com/dentych/demic/card"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	ErrGameStarted      = fmt.Errorf("game is started")
	ErrGameNotStarted   = fmt.Errorf("game is not started")
	ErrPlayerNameExists = fmt.Errorf("player name already exists")
	ErrNoMoreCards      = fmt.Errorf("no more cards to turn")
)

type Pyramid struct {
	Input  chan string
	Output chan string

	roomId         string
	started        bool
	players        []Player
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
	p.roomId = GenerateId(4)
	p.players = make([]Player, 0)
	p.deck = card.NewDeck()
	p.boardCardIndex = 0
	p.attackState = false

	p.Input = make(chan string, 25)
	p.Output = make(chan string, 25)

	return &p
}

func (p *Pyramid) AddPlayer(player *Player) error {
	if p.started {
		return ErrGameStarted
	}

	for _, v := range p.players {
		if player.Name == v.Name {
			return ErrPlayerNameExists
		}
	}

	p.players = append(p.players, *player)
	return nil
}

func (p *Pyramid) Play() {
	p.started = true
	go p.inputHandler()
	p.dealCards()

	p.waitForContinue()

	for p.boardCardIndex != len(p.board) {
		c, err := p.turnNextCard()
		if err != nil {
			log.Panic(err)
		}
		p.Output <- "CARD " + string(c.Suit) + c.Rank

		p.Output <- "ATTACK BEGIN"
		p.attackState = true
		p.waitForContinue()
		p.attackState = false
		p.Output <- "ATTACK STOP"
	}
}

func (p *Pyramid) Continue() {
	p.cont = true
}

func (p *Pyramid) dealCards() {
	p.started = true
	p.deck = card.Shuffle(p.deck)
	p.board, p.deck = card.Deal(p.deck, 10)
	for k := range p.players {
		p.players[k].Hand, p.deck = card.Deal(p.deck, 4)
	}
}

func (p *Pyramid) turnNextCard() (*card.Card, error) {
	if !p.started {
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
	p.Output <- attacker.Name + " ATTACKS " + attackee.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
}

func (p *Pyramid) acceptAttack(attacker, attackee *Player, dmg int) {
	p.Output <- attackee.Name + " ACCEPTS ATTACK FROM " + attacker.Name + " FOR " + strconv.Itoa(dmg) + " DAMAGE!"
	attackee.Sips += dmg
}

func (p *Pyramid) inputHandler() {
	//Forslag til struktur: Input = [roomid, acting player, action, message]
	for {
		event := <-p.Input
		s := strings.Split(event, " ")
		if len(s) < 1 {
			fmt.Println("MESSAGE NOT UNDERSTOOD: " + event)
			continue
		}
		switch s[2] {
		//Forslag til struktur: Input = [roomid, acting player, action = "ATTACK", target, dmg]
		case "ATTACK":
			if len(s) < 4 {
				fmt.Println("SHIT ATTACK MESSAGE: " + event)
				continue
			}
			if !p.attackState {
				p.Output <- "ATTACKING FAILED, NOT IN ATTACKING STATE"
			} else {
				var attackingPlayer, targetPlayer *Player
				p.Output <- "ATTACKING " + s[1] + " " + s[2]
				for k := range p.players {
					switch p.players[k].Name {
					case s[1]:
						attackingPlayer = &p.players[k]
					case s[3]:
						targetPlayer = &p.players[k]
					}
				}
				k, _ := strconv.Atoi(s[4])
				p.attack(attackingPlayer, targetPlayer, k)
			}
		//Forslag til struktur: Input = [roomid, acting player, action = "ATTACK", target, dmg]
		case "ACCEPT_ATTACK":
			if len(s) < 3 {
				fmt.Println("SHIT ACCEPT_ATTACK MESSAGE: " + event)
				continue
			}
			if !p.attackState {
				p.Output <- "ACCEPT_ATTACK FAILED, NOT IN ATTACKING STATE"
			} else {
				var attackingPlayer, targetPlayer *Player
				p.Output <- "ACCEPT_ATTACK " + s[1] + " " + s[2]
				for k := range p.players {
					switch p.players[k].Name {
					case s[1]:
						attackingPlayer = &p.players[k]
					case s[3]:
						targetPlayer = &p.players[k]
					}
				}
				k, _ := strconv.Atoi(s[4])
				p.acceptAttack(attackingPlayer, targetPlayer, k)

			}
		case "REJECT_ATTACK":
			if len(s) < 3 {
				fmt.Println("SHIT REJECT_ATTACK MESSAGE: " + event)
				continue
			}
			if !p.attackState {
				p.Output <- "REJECT_ATTACK FAILED, NOT IN ATTACKING STATE"
			} else {
				p.Output <- "REJECT_ATTACK " + s[1] + " " + s[2]
			}
		}
	}
}
