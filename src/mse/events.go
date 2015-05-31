package mse

import (
	"fmt"
	"math/rand"

	"interact"
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

func handleRevolt(g *Game) interact.GameState {
	if len(g.Empire) == 1 {
		if g.Year == 1 {
			g.Log("The Home World won't revolt in year 1.")
			return StartState
		}
		g.Log("The Home World has revolted.")
		return LoseState
	}

	w := g.lowestResistanceWorld()

	forceMod := g.ActiveEvent.Year1Effect
	if g.Year == 2 {
		forceMod = g.ActiveEvent.Year2Effect
	}

	r := w.Resistance
	resistanceMod := ""
	if g.Techs[HyperTelevision] {
		r += 1
		resistanceMod = " (+1 for Hyper Television)"
	}
	roll := Roll()
	result := "failed"
	f := map[string]int{
		"Force +1": 1,
		"Force +2": 2,
		"Force +3": 3,
	}[forceMod]
	if roll+f >= r {
		result = "succeeded"
	}

	g.Logf(
		"Revolt on %s: %s, Resistance of %d%s, rolled %d...revolt %s!",
		w.Name, forceMod, w.Resistance, resistanceMod, roll, result)

	if result == "succeeded" {
		w.Revolted = true
		g.empireToExplored(w)
	}

	return EndOfTurnState
}

func (g *Game) lowestResistanceWorld() *SystemCard {
	minR := 0
	nonHome := g.Empire[1:]
	for _, w := range nonHome {
		if minR == 0 || w.Resistance < minR {
			minR = w.Resistance
		}
	}
	var worlds []*SystemCard
	for _, w := range nonHome {
		if w.Resistance == minR {
			worlds = append(worlds, w)
		}
	}
	return worlds[rand.Intn(len(worlds))]
}

func handleInvasion(g *Game) interact.GameState {
	if len(g.Empire) == 1 {
		if g.Year == 1 {
			g.Log("Invasion force won't attack the Home World in year 1.")
			return StartState
		}
		g.Log("The Home World has been invaded.")
		return LoseState
	}

	w := g.Empire[len(g.Empire)-1]

	forceMod := g.ActiveEvent.Year1Effect
	if g.Year == 2 {
		forceMod = g.ActiveEvent.Year2Effect
	}

	r := w.Resistance
	resistanceMod := ""
	if g.Techs[PlanetaryDefenses] {
		r += 1
		resistanceMod = " (+1 for Planetary Defenses)"
	}
	roll := Roll()
	result := "failed"
	f := map[string]int{
		"Force +1": 1,
		"Force +2": 2,
		"Force +3": 3,
	}[forceMod]
	if roll+f >= r {
		result = "succeeded"
	}

	g.Logf(
		"Invasion on %s: %s, Resistance of %d%s, rolled %d...invasion %s!",
		w.Name, forceMod, w.Resistance, resistanceMod, roll, result)

	if result == "succeeded" {
		w.Invaded = true
		g.empireToExplored(w)
	}

	return EndOfTurnState
}

func (g *Game) isStrikeActive() bool {
	return g.ActiveEvent != nil && g.ActiveEvent.Name == Strike
}
