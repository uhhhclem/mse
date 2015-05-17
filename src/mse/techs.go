package mse

const (
	CapitalShips          = iota + 1
	RobotWorkers          
	HyperTelevision       
	InterspeciesCommerce  
	ForwardStarbases      
	PlanetaryDefenses     
	InterstellarDiplomacy 
	InterstellarBanking   
)

type Tech struct {
	ID        int
	Name      string
	Ability   string
	Cost      int
	DependsOn int
	Enables   int
}

var Techs map[int]Tech

func init() {
	Techs = map[int]Tech{
		CapitalShips: {
			Name:    "Capital Ships",
			Ability: "Advance beyond military strength 3",
			Cost:    3,
			Enables: ForwardStarbases,
		},
		RobotWorkers: {
			Name:    "Robot Workers",
			Cost:    2,
			Ability: "Receive 1/2 production during strike",
			Enables: PlanetaryDefenses,
		},
		HyperTelevision: {
			Name:    "Hyper Television",
			Cost:    3,
			Ability: "+1 to resistance during revolt",
			Enables: InterstellarDiplomacy,
		},
		InterspeciesCommerce: {
			Name:    "Interspecies Commerce",
			Cost:    2,
			Ability: "Exchange 2 of one resource for 1 of the other",
			Enables: InterstellarBanking,
		},
		ForwardStarbases: {
			Name:      "Forward Starbases",
			Cost:      4,
			Ability:   "Required to explore distant systems",
			DependsOn: CapitalShips,
		},
		PlanetaryDefenses: {
			Name:      "Planetary Defenses",
			Cost:      4,
			Ability:   "+1 to resistance during invasion",
			DependsOn: RobotWorkers,
		},
		InterstellarDiplomacy: {
			Name:      "Interstellar Diplomacy",
			Cost:      5,
			Ability:   "Next planet is conquered for free",
			DependsOn: HyperTelevision,
		},
		InterstellarBanking: {
			Name:      "Interstellar Banking",
			Cost:      3,
			Ability:   "Advance beyond storage value 3",
			DependsOn: InterspeciesCommerce,
		},
	}
}
