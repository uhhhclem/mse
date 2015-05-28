package mse

type Board struct {
	State                   string
	Year                    int
	MetalProduction         int
	WealthProduction        int
	MetalStorage            int
	WealthStorage           int
	MilitaryStrength        int
	Empire                  []*SystemCard
	Explored                []*SystemCard
	Gen1Techs               []TechDisplay
	Gen2Techs               []TechDisplay
	ActiveEvent             *EventCard
	EventsRemaining         int
	NearSystemsRemaining    int
	DistantSystemsRemaining int
}

type TechDisplay struct {
	ID      string
	Cost    int
	Name    string
	Ability string
	Owned   bool
}

type StatusResponse struct {
	Status
	End bool
}

type PromptResponse struct {
	Prompt
	End bool
}

func (g *Game) GetBoard() *Board {
	b := &Board{
		State:                   string(g.State),
		Year:                    g.Year,
		MetalProduction:         g.MetalProduction,
		WealthProduction:        g.WealthProduction,
		MetalStorage:            g.MetalStorage,
		WealthStorage:           g.WealthStorage,
		MilitaryStrength:        g.MilitaryStrength,
		Empire:                  g.Empire,
		Explored:                g.Explored,
		ActiveEvent:             g.ActiveEvent,
		EventsRemaining:         len(g.EventDeck),
		NearSystemsRemaining:    len(g.NearSystemDeck),
		DistantSystemsRemaining: len(g.DistantSystemDeck),
	}
	b.Gen1Techs = []TechDisplay{
		g.getTechDisplay(CapitalShips),
		g.getTechDisplay(RobotWorkers),
		g.getTechDisplay(HyperTelevision),
		g.getTechDisplay(InterspeciesCommerce),
	}
	b.Gen2Techs = []TechDisplay{
		g.getTechDisplay(ForwardStarbases),
		g.getTechDisplay(PlanetaryDefenses),
		g.getTechDisplay(InterstellarDiplomacy),
		g.getTechDisplay(InterstellarBanking),
	}
	return b
}

func (g *Game) getTechDisplay(id string) TechDisplay {
	t := Techs[id]
	return TechDisplay{
		Name:    t.Name,
		Ability: t.Ability,
		Cost:    t.Cost,
		Owned:   g.Techs[id],
	}
}
