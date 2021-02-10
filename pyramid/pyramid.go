package pyramid

import (
	"bytes"
	"fmt"
	"gitlab.com/dentych/demic/card"
	"log"
	"sort"
	"strconv"
	"sync"
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
	Input       chan Action  `json:"-"`
	PlayerJoin  chan *Player `json:"-"`
	PlayerLeave chan *Player `json:"-"`

	RoomId         string   `json:"room_id"`
	Started        bool     `json:"started"`
	Players        []Player `json:"players"`
	Attacks        Attacks  `json:"attacks"`
	boardCardIndex int
	board          []card.Card
	deck           []card.Card
	cont           bool
	attackState    bool
}

type Attack struct {
	Attacker Player
	Target   Player
}

func (a *Attack) EqualTo(other Attack) bool {
	if a.Attacker.Name == "" || a.Target.Name == "" {
		return false
	}
	return a.Attacker.Name == other.Attacker.Name && a.Target.Name == other.Target.Name
}

type Attacks struct {
	Attacks []Attack `json:"attacks"`
	mutex   sync.Mutex
}

func (a *Attacks) Add(attack Attack) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Attacks = append(a.Attacks, attack)
}

func (a *Attacks) Remove(attack Attack) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.Attacks) == 0 {
		return
	}

	attackIndex := 0
	for k, v := range a.Attacks {
		if v.EqualTo(attack) {
			attackIndex = k
		}
	}

	a.Attacks[attackIndex], a.Attacks[len(a.Attacks)-1] = a.Attacks[len(a.Attacks)-1], a.Attacks[attackIndex]
	a.Attacks = a.Attacks[:len(a.Attacks)-1]
}

func (a *Attacks) Len() int {
	return len(a.Attacks)
}

func init() {
	Rooms = make(map[string]*Pyramid)
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
	go p.inputHandler()
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
		//p.setAttackState(true)
		p.waitForContinue()
		//p.setAttackState(false)
	}
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

	p.Players[1].Output <- Action{
		ActionType: ActionAttackState,
		Origin:     p.Players[0].Name,
		Target:     strconv.FormatBool(p.attackState),
	}

}

func (p *Pyramid) continueGame() {
	p.cont = true
}

func (p *Pyramid) output(action Action) {
	for _, player := range p.Players {
		player.Output <- action
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

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) rejectAttack(event Action) error {
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

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionRejectAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) acceptAttack(event Action) error {
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

	p.Players[targetIdx].Output <- Action{
		ActionType: ActionAcceptAttack,
		Origin:     p.Players[originIdx].Name,
		Target:     p.Players[targetIdx].Name,
	}
	return nil
}

func (p *Pyramid) pickCard(event Action) error {
	var originIdx, handIdx int

	for k, player := range p.Players {
		switch player.Name {
		case event.Origin:
			originIdx = k
		}
	}

	chosenCard := p.Players[originIdx].Hand[handIdx].String()

	p.Players[0].Output <- Action{
		ActionType: ActionPickCard,
		Origin:     p.Players[originIdx].Name,
		Target:     chosenCard,
	}
	newCard := card.Deal(&p.deck, 1)[0]
	p.Players[originIdx].Hand[handIdx] = newCard
	p.Players[originIdx].Output <- Action{ActionType: ActionDealHand, Target: newCard.String()}
	return nil
}

func (p *Pyramid) newCard(event Action) {
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

	p.Players[originIdx].Output <- Action{
		ActionType: ActionDealHand,
		Origin:     p.Players[originIdx].Name,
		Target:     cardStr,
	}
}

func (p *Pyramid) inputHandler() {
	for {
		event := <-p.Input
		switch event.ActionType {
		case ActionStartGame:
			p.Started = true
		case ActionAttack:
			p.attack(event)
		case ActionAcceptAttack:
			p.acceptAttack(event)
		case ActionRejectAttack:
			p.rejectAttack(event)
		case ActionPickCard:
			p.pickCard(event)
		case ActionNewCard:
			p.newCard(event)
		case ActionContinue:
			p.continueGame()
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
