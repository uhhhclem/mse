package mse

import (
	"fmt"
)

func (g *Game) doAsteroidEvent() string {
	wealth := g.addWealth(1)
	return fmt.Sprintf("Added %d wealth", wealth)
}

func (g *Game) doDerelictShipEvent() string {
	metal := g.addMetal(1)
	return fmt.Sprintf("Added %d metal", metal)
}

func (g *Game) doPeaceAndQuietEvent() string {
	return "No effect"
}

func (g *Game) doStrikeEvent() string {
	if g.Techs[RobotWorkers] {
		return "Production halved next turn."
	}
	return "No production next turn."
}

func handleRevolt(g *Game) GameState {
	g.Log("Revolt not implemented.")
	return StartState
}

func handleSmallInvasionForce(g *Game) GameState {
	g.Log("Small Invasion Force not implemented.")
	return StartState
}

func handleLargeInvasionForce(g *Game) GameState {
	g.Log("Large Invasion Force not implemented.")
	return StartState
}

func (g *Game) isStrikeActive() bool {
	return g.ActiveEvent != nil && g.ActiveEvent.Name == Strike
}
