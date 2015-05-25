package mse

import (
	"fmt"
)

type GameState string

const (
	StartState        GameState = "StartOfTurn"
	SystemChoiceState           = "SystemChoice"
	AttackState                 = "Attack"
	CollectState                = "Collect"
	ChooseBuildState = "ChooseBuild"
	DoBuildState = "DoBuild"
	EventState = "Event"
	EndState                    = "End"
)

type stateHandler func(*Game) GameState

var handlers map[GameState]stateHandler

func init() {
	handlers = map[GameState]stateHandler{
		StartState:   handleStart,
		AttackState:  handleAttack,
		CollectState: handleCollect,
		ChooseBuildState: handleChooseBuild,
		DoBuildState: handleDoBuild,
		EventState: handleEvent,
	}
}

type Game struct {
	State                                        GameState
	NearSystemDeck, DistantSystemDeck, EventDeck Deck
	Empire                                       []*SystemCard
	Explored                                     []*SystemCard
	ActiveEvent                                  *EventCard
	Techs                                        map[string]bool
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
		EventDeck:         []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		NearSystemDeck:    []string{"2", "3", "4", "5", "6", "7", "8"},
		DistantSystemDeck: []string{"9", "10", "11"},
		Empire:            []*SystemCard{Systems["1"]},
		Techs:             make(map[string]bool),
		NextStatus:        make(chan *Status),
		NextPrompt:        make(chan *Prompt),
		NextChoice:        make(chan *Choice),
		Ready:             make(chan bool),
	}
	shuffle(g.EventDeck)
	shuffle(g.NearSystemDeck)
	shuffle(g.DistantSystemDeck)

	g.calculateProduction()

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
	if g.mayExploreDistantSystems() {
		return SystemChoiceState
	}

	var id string
	id, g.NearSystemDeck = Draw(g.NearSystemDeck)
	sc := Systems[id]
	g.Logf("Explored %s", sc.Name)
	g.Explored = append(g.Explored, sc)

	choices := make([]*Choice, 0, len(g.Explored)+1)
	for _, sc := range g.Explored {
		choices = append(choices, &Choice{
			Key:  sc.ID,
			Name: fmt.Sprintf("Attack %s", sc.Name),
		})
	}
	choices = append(choices, &Choice{"B", "Bide your time"})

	g.SendPrompt("Select a system to attack, or bide your time.", choices)

	return AttackState
}

func handleAttack(g *Game) GameState {
	c := <-g.NextChoice
	if c.Key == "B" {
		g.Log("Biding time...")
		return CollectState
	}
	w := Systems[c.Key]
	g.Logf("Attacking %s...", w.Name)

	roll := Roll()
	success := roll+g.MilitaryStrength >= w.Resistance
	result := map[bool]string{
		true:  "success",
		false: "failed",
	}[success]

	if success {
		g.exploredToEmpire(w)
	}
	g.Logf("Resistance = %d, military strength = %d, roll = %d...%s!",
		w.Resistance, g.MilitaryStrength, roll, result)
	if !success && g.MilitaryStrength > 0 {
		g.MilitaryStrength -= 1
		g.Logf("Military strength reduced to %d.", g.MilitaryStrength)
	}

	return CollectState
}

func handleCollect(g *Game) GameState {
	g.calculateProduction()
	metal, wealth := g.MetalStorage, g.WealthStorage
	g.collect()
	g.Logf("Collected %d metal and %d wealth.", g.MetalStorage - metal,
			g.WealthStorage - wealth)
	return ChooseBuildState
}

const (
	BuildDone = "Done"
	BuildMilitary = "Military"
)

func handleChooseBuild(g *Game) GameState {
	var choices = []*Choice{{Key: BuildDone, Name: "Done"}}
	if g.WealthStorage > 0 && g.MetalStorage > 0 && g.MilitaryStrength < 6 {
		choices = append(choices, &Choice{
			Key: BuildMilitary, 
			Name: "Increase military strength by 1 (1 wealth, 1 metal)",
		})
	}
	
	g.SendPrompt("Select build:", choices)
	
	return DoBuildState
}

func handleDoBuild(g *Game) GameState {
	c := <-g.NextChoice
	switch c.Key {
	case BuildDone:
	  return EventState
	case BuildMilitary:
	  g.MilitaryStrength += 1
	  g.MetalStorage -= 1
	  g.WealthStorage -= 1
	default:
	  g.Logf("Unknown build key %q", c.Key)
	}
	return ChooseBuildState
}

func handleEvent(g *Game) GameState {
	g.Log("Event state not implemented.")
	return StartState
}

func (g *Game) mayExploreDistantSystems() bool {
	return false
}

func (g *Game) exploredToEmpire(sc *SystemCard) {
	for i := range g.Explored {
		if g.Explored[i] == sc {
			g.Explored = append(g.Explored[:i], g.Explored[i+1:]...)
			break
		}
	}
	g.Empire = append(g.Empire, sc)
}

func (g *Game) empireToExplored(sc *SystemCard) {
	for i := range g.Empire {
		if g.Empire[i] == sc {
			g.Empire = append(g.Empire[:i], g.Empire[i+1:]...)
			break
		}
	}
	g.Explored = append(g.Explored, sc)
}
