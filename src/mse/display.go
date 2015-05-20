package mse

type Board struct {
	MetalProduction  int
	WealthProduction int
	MetalStorage     int
	WealthStorage    int
	MilitaryStrength int
	Empire           []*SystemCard
	Explored         []*SystemCard
	Gen1Techs        []TechDisplay
	Gen2Techs        []TechDisplay
}

type TechDisplay struct {
	ID      int
	Cost    int
	Name    string
	Ability string
	Owned   bool
}

func (g *Game) GetBoard() *Board {
	b := &Board{
		MetalProduction:  g.MetalProduction,
		WealthProduction: g.WealthProduction,
		MetalStorage:     g.MetalStorage,
		WealthStorage:    g.WealthStorage,
		MilitaryStrength: g.MilitaryStrength,
		Empire:           g.Empire,
		Explored:         g.Explored,
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

func (g *Game) getTechDisplay(id int) TechDisplay {
	t := Techs[id]
	return TechDisplay{
		Name:    t.Name,
		Ability: t.Ability,
		Cost:    t.Cost,
		Owned:   g.Techs[id],
	}
}