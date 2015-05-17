package mse

type Game struct {
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
}

func NewGame() *Game {
	g := &Game{
		EventDeck:         []int{1, 2, 3, 4, 5, 6, 7, 8},
		NearSystemDeck:    []int{2, 3, 4, 5, 6, 7, 8},
		DistantSystemDeck: []int{9, 10, 11},
		Empire:            []*SystemCard{Systems[1]},
		Techs:             make(map[int]bool),
	}
	shuffle(g.EventDeck)
	shuffle(g.NearSystemDeck)
	shuffle(g.DistantSystemDeck)

	g.calculateProduction()

	return g
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
