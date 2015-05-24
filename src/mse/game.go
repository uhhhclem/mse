package mse

import (
	"fmt"
)

type GameState string

const (
	StartState        GameState = "StartOfTurn"
	SystemChoiceState           = "SystemChoice"
	AttackState                 = "Attack"
	EndState                    = "End"
)

type stateHandler func(*Game) GameState

var handlers map[GameState]stateHandler

func init() {
	handlers = map[GameState]stateHandler{
		StartState:  handleStart,
		AttackState: handleAttack,
	}
}

type Game struct {
	State                                        GameState
	NearSystemDeck, DistantSystemDeck, EventDeck Deck
	Empire                                       []*SystemCard
	Explored                                     []*SystemCard
	ActiveEvent                                  *EventCard
	Techs                                        map[int]bool
	UsedInterstellarDiplomacy                    bool
	MetalStorage                                 int
	WealthStorage                                int
	MilitaryStrength                             int
	MetalProduction                              int
	WealthProduction                             int
	Prompt                                       *Prompt
	NextPrompt                                   chan *Prompt
	NextStatus                                   chan *Status
	NextChoice                                   chan *Choice
	Ready                                        chan bool
}

func NewGame() *Game {
	g := &Game{
		EventDeck:         []int{1, 2, 3, 4, 5, 6, 7, 8},
		NearSystemDeck:    []int{2, 3, 4, 5, 6, 7, 8},
		DistantSystemDeck: []int{9, 10, 11},
		Empire:            []*SystemCard{Systems[1]},
		Techs:             make(map[int]bool),
		NextStatus:        make(chan *Status),
		NextPrompt:        make(chan *Prompt),
		NextChoice:        make(chan *Choice),
		Ready:             make(chan bool),
	}
	shuffle(g.EventDeck)
	shuffle(g.NearSystemDeck)
	shuffle(g.DistantSystemDeck)

	g.State = StartState

	return g
}

func (g *Game) Run() {
	for {
		s := g.State
		if s == EndState {
			g.NextStatus <- nil
			g.NextPrompt <- nil
			return
		}
		g.State = handlers[s](g)
		g.Ready <- true
	}
}

func (g *Game) calculateProduction() {
	g.MetalProduction, g.WealthProduction = 0, 0
	for _, sc := range g.Empire {
		g.MetalProduction += sc.Metal
		g.WealthProduction += sc.Wealth
	}

	if g.ActiveEvent == nil || g.ActiveEvent.Name == Strike {
		return
	}

	if !g.Techs[RobotWorkers] {
		g.MetalProduction = 0
		g.WealthProduction = 0
	} else {
		g.MetalProduction = g.MetalProduction/2 + g.MetalProduction%2
		g.WealthProduction = g.WealthProduction/2 + g.WealthProduction%2
	}
}

func (g *Game) collect() {
	g.MetalStorage += g.MetalProduction
	if g.MetalStorage > 3 && !g.Techs[InterstellarDiplomacy] {
		g.MetalStorage = 3
	}
	if g.MetalStorage > 5 {
		g.MetalStorage = 5
	}
	g.WealthStorage += g.WealthProduction
	if g.WealthStorage > 3 && !g.Techs[InterstellarDiplomacy] {
		g.WealthStorage = 3
	}
	if g.WealthStorage > 5 {
		g.WealthStorage = 5
	}
}

func handleStart(g *Game) GameState {
	g.calculateProduction()
	if g.mayExploreDistantSystems() {
		return SystemChoiceState
	}

	var id int
	id, g.NearSystemDeck = Draw(g.NearSystemDeck)
	sc := Systems[id]
	g.Logf("Explored %s", sc.Name)
	g.Explored = append(g.Explored, sc)

	choices := make([]Choice, len(g.Explored)+1)
	for i, sc := range g.Explored {
		choices[i] = Choice{
			Key:  fmt.Sprintf("%s", sc.ID),
			Name: fmt.Sprintf("Attack %s", sc.Name),
		}
	}
	choices = append(choices, Choice{"B", "Bide your time"})

	g.Prompt = &Prompt{
		State:   g.State,
		Message: "Select a system to attack, or bide your time.",
		Choices: choices,
	}

	return AttackState
}

func handleAttack(g *Game) GameState {
	g.Log("Ending game...")
	return EndState
}

func (g *Game) mayExploreDistantSystems() bool {
	return false
}
