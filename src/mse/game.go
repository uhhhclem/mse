package mse

import (
	"fmt"
	
	"code.google.com/p/go-uuid/uuid"
)

type Game struct {
	ID										     string
	State                                        GameState
	Year                                         int
	NearSystemDeck, DistantSystemDeck, EventDeck Deck
	Empire                                       []*SystemCard
	Explored                                     []*SystemCard
	ActiveEvent                                  *EventCard
	Techs                                        map[string]bool
	UsedTech                                     map[string]bool
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

type GameState string

const (
	StartState              GameState = "StartOfTurn"
	PhaseIState                       = "PhaseI"
	CollectState                      = "Collect"
	ChooseBuildState                  = "ChooseBuild"
	DoBuildState                      = "DoBuild"
	EventState                        = "Event"
	RevoltState                       = "Revolt"
	SmallInvasionForceState           = "SmallInvasionForce"
	LargeInvasionForceState           = "LargeInvasionForce"
	EndOfTurnState                    = "EndOfTurn"
	WinState                          = "Win"
	LoseState                         = "Lose"
	EndState                          = "End"
)

type stateHandler func(*Game) GameState

var handlers map[GameState]stateHandler

const (
	BuildDone            = "Done"
	BuildMilitary        = "Military"
	BuildWealthFromMetal = "Wealth"
	BuildMetalFromWealth = "Metal"
)

var buildChoices map[string]string

func init() {
	handlers = map[GameState]stateHandler{
		StartState:              handleStart,
		PhaseIState:             handlePhaseI,
		CollectState:            handleCollect,
		ChooseBuildState:        handleChooseBuild,
		DoBuildState:            handleDoBuild,
		EventState:              handleEvent,
		RevoltState:             handleRevolt,
		SmallInvasionForceState: handleInvasion,
		LargeInvasionForceState: handleInvasion,
		EndOfTurnState:          handleEndOfTurn,
		WinState:                handleWin,
		LoseState:               handleLose,
	}

	buildChoices = map[string]string{
		BuildDone:            "Done building",
		BuildMilitary:        "Increase military strength (cost: 1 wealth, 1 metal)",
		BuildWealthFromMetal: "Exchange 2 metal for 1 wealth",
		BuildMetalFromWealth: "Exchange 2 wealth for 1 metal",
	}

}

func NewGame() *Game {
	g := &Game{
		ID:                uuid.New(),
		EventDeck:         []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		NearSystemDeck:    []string{"2", "3", "4", "5", "6", "7", "8"},
		DistantSystemDeck: []string{"9", "10", "11"},
		Year:              1,
		Empire:            []*SystemCard{Systems["1"]},
		Techs:             make(map[string]bool),
		UsedTech:          make(map[string]bool),
		NextStatus:        make(chan *Status),
		NextPrompt:        make(chan *Prompt),
		NextChoice:        make(chan *Choice),
		Ready:             make(chan bool),
	}
	shuffle(g.EventDeck)
	_, g.EventDeck = Draw(g.EventDeck)

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

	if !g.isStrikeActive() {
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

func (g *Game) addMetal(add int) int {
	orig := g.MetalStorage
	max := g.maxStorage()
	g.MetalStorage += add
	if g.MetalStorage > max {
		g.MetalStorage = max
	}
	return g.MetalStorage - orig
}

func (g *Game) addWealth(add int) int {
	orig := g.WealthStorage
	max := g.maxStorage()
	g.WealthStorage += add
	if g.WealthStorage > max {
		g.WealthStorage = max
	}
	return g.WealthStorage - orig
}

func (g *Game) maxStorage() int {
	if g.Techs[InterstellarBanking] {
		return 5
	}
	return 3
}

func handleStart(g *Game) GameState {
	// Interspecies Commerce is usable once a turn
	g.UsedTech[InterspeciesCommerce] = false

	choices := make([]*Choice, 0, len(g.Explored)+2)
	if g.mayExploreAndAttack() {
		choices = append(choices, &Choice{"X", "Explore and attack"})
	}

	for _, sc := range g.Explored {
		choices = append(choices, &Choice{
			Key:  sc.ID,
			Name: fmt.Sprintf("Conquer %s", sc.Name),
		})
	}
	choices = append(choices, &Choice{"B", "Bide your time"})

	g.SendPrompt("Select a system to attack, or bide your time.", choices)

	return PhaseIState
}

func (g *Game) mayExploreAndAttack() bool {
	if len(g.NearSystemDeck) > 0 {
		return true
	}
	if g.Techs[ForwardStarbases] && len(g.DistantSystemDeck) > 0 {
		return true
	}
	return false
}

func handlePhaseI(g *Game) GameState {
	c := <-g.NextChoice

	if c.Key == "B" {
		g.Log("Biding time...")
		return CollectState
	}

	var w *SystemCard
	if c.Key == "X" {
		w = g.exploreWorld()
	} else {
		w = Systems[c.Key]
	}

	if g.mayMakeFreeAttack() {
		g.Logf("%s conquered through interstellar diplomacy.", w.Name)
		g.exploredToEmpire(w)
		g.UsedTech[InterstellarDiplomacy] = true
		return CollectState
	}

	g.Logf("Attacking %s...", w.Name)

	roll := Roll()
	result := "failed"
	modifier := ""

	r := w.Resistance
	if w.Revolted && g.Techs[HyperTelevision] {
		r += 1
		modifier = " (+1 for previous revolt)"
	}
	if w.Invaded && g.Techs[PlanetaryDefenses] {
		r += 1
		modifier = " (+1 for invaded world's defenses)"
	}

	success := roll+g.MilitaryStrength >= r

	if success {
		result = "success"
		g.exploredToEmpire(w)
		w.Revolted = false
		w.Invaded = false
	}
	g.Logf("Resistance = %d%s, military strength = %d, roll = %d...%s!",
		w.Resistance, modifier, g.MilitaryStrength, roll, result)
	if !success && g.MilitaryStrength > 0 {
		g.MilitaryStrength -= 1
		g.Logf("Military strength reduced to %d.", g.MilitaryStrength)
	}

	return CollectState
}

func (g *Game) exploreWorld() *SystemCard {
	var id string
	if len(g.NearSystemDeck) > 0 {
		id, g.NearSystemDeck = Draw(g.NearSystemDeck)
	} else {
		id, g.DistantSystemDeck = Draw(g.DistantSystemDeck)
	}
	w := Systems[id]
	g.Explored = append(g.Explored, w)
	g.Logf("Explored %s.", w.Name)
	return w
}

func handleCollect(g *Game) GameState {
	g.calculateProduction()
	metal := g.addMetal(g.MetalProduction)
	wealth := g.addWealth(g.WealthProduction)
	g.Logf("Collected %d metal and %d wealth.", metal, wealth)
	return ChooseBuildState
}

func handleChooseBuild(g *Game) GameState {
	var choices []*Choice
	choices = addChoice(choices, BuildDone)

	maxMilitary := 3
	if g.mayIncreaseMilitaryAbove3() {
		maxMilitary = 5
	}
	if g.WealthStorage > 0 && g.MetalStorage > 0 && g.MilitaryStrength < maxMilitary {
		choices = addChoice(choices, BuildMilitary)
	}

	if g.mayExchangeGoods() {
		if g.MetalStorage > 1 && g.WealthStorage < g.maxStorage() {
			choices = addChoice(choices, BuildWealthFromMetal)
		}
		if g.WealthStorage > 1 && g.MetalStorage < g.maxStorage() {
			choices = addChoice(choices, BuildMetalFromWealth)
		}
	}

	for k, t := range Techs {
		if g.Techs[k] {
			continue
		}
		if t.DependsOn != "" && !g.Techs[t.DependsOn] {
			continue
		}
		if t.Cost > g.WealthStorage {
			continue
		}
		choices = addChoice(choices, k)
	}

	g.SendPrompt("Select build:", choices)

	return DoBuildState
}

func addChoice(choices []*Choice, key string) []*Choice {
	if name, ok := buildChoices[key]; ok {
		c := &Choice{Key: key, Name: name}
		choices = append(choices, c)
	}
	if t, ok := Techs[key]; ok {
		c := &Choice{Key: key, Name: fmt.Sprintf("Build %s", t.Name)}
		choices = append(choices, c)
	}
	return choices
}

func handleDoBuild(g *Game) GameState {
	c := <-g.NextChoice

	if t, ok := Techs[c.Key]; ok {
		g.Techs[c.Key] = true
		g.Logf("Bought %s.", t.Name)
		g.WealthStorage -= t.Cost
		if c.Key == InterstellarDiplomacy {
			g.Log("If you attack next turn, it will automatically succeed.")
		}
		return ChooseBuildState
	}

	switch c.Key {
	case BuildDone:
		return EventState
	case BuildMilitary:
		g.MilitaryStrength += 1
		g.MetalStorage -= 1
		g.WealthStorage -= 1
		g.Log("Built 1 military.")
	case BuildWealthFromMetal:
		g.MetalStorage -= 2
		g.WealthStorage += 1
		g.UsedTech[InterspeciesCommerce] = true
		g.Log("Swapped 2 metal for 1 wealth.")
	case BuildMetalFromWealth:
		g.WealthStorage -= 2
		g.MetalStorage += 1
		g.UsedTech[InterspeciesCommerce] = true
		g.Log("Swapped 2 wealth for 1 metal.")
	default:
		g.Logf("Unknown build key %q", c.Key)
	}
	return ChooseBuildState
}

func handleEvent(g *Game) GameState {
	var id string
	id, g.EventDeck = Draw(g.EventDeck)
	e := Events[id]
	g.ActiveEvent = e
	g.Logf("Drew event: %s", e.Name)

	var result string
	switch g.ActiveEvent.Name {
	case Asteroid:
		result = g.doAsteroidEvent()
	case PeaceAndQuiet:
		result = g.doPeaceAndQuietEvent()
	case DerelictShip:
		result = g.doDerelictShipEvent()
	case Strike:
		result = g.doStrikeEvent()
	case Revolt:
		return RevoltState
	case SmallInvasionForce:
		return SmallInvasionForceState
	case LargeInvasionForce:
		return LargeInvasionForceState
	default:
		result = fmt.Sprintf("No handler defined for event %s", g.ActiveEvent.Name)
	}

	g.Log(result)
	return EndOfTurnState
}

func handleEndOfTurn(g *Game) GameState {
	if len(g.EventDeck) == 0 {
		g.Logf("End of Year %d.", g.Year)
		if g.Year == 2 {
			return WinState
		}
		g.Year += 1
		g.EventDeck = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
		shuffle(g.EventDeck)
		_, g.EventDeck = Draw(g.EventDeck)
		_, g.EventDeck = Draw(g.EventDeck)
	}
	return StartState
}

func handleWin(g *Game) GameState {
	empireVPs, techVPs := 0, 0
	for _, sc := range g.Empire {
		empireVPs += sc.VPs
	}
	for id := range g.Techs {
		if g.Techs[id] {
			techVPs += 1
		}
	}

	vps := 0
	g.Logf("%d VPs from your empire.", empireVPs)
	vps += empireVPs
	g.Logf("%d VPs from discovered technologies.", techVPs)
	vps += techVPs
	if len(g.DistantSystemDeck) == 0 {
		g.Logf("Exploration Bonus (1VP) for exploring all systems.")
		vps += 1
	}
	if techVPs == 8 {
		g.Logf("Scientific Bonus (1VP) for researching all technologies.")
		vps += 1
	}
	if len(g.DistantSystemDeck) == 0 && len(g.Explored) == 0 {
		g.Logf("Warlord Bonus (3VP) for conquering all systems.")
		vps += 3
	}
	g.Logf("Final score: %d VPs.", vps)

	return EndState
}

func handleLose(g *Game) GameState {
	g.Log("You lose.")
	return EndState
}

func (g *Game) mayExchangeGoods() bool {
	return g.Techs[InterspeciesCommerce] && !g.UsedTech[InterspeciesCommerce]
}

func (g *Game) mayMakeFreeAttack() bool {
	return g.Techs[InterstellarDiplomacy] && !g.UsedTech[InterstellarDiplomacy]
}

func (g *Game) mayIncreaseMilitaryAbove3() bool {
	return g.Techs[CapitalShips]
}

func (g *Game) mayExploreDistantSystems() bool {
	return g.Techs[ForwardStarbases]
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
