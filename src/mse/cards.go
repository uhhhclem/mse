package mse

import (
	"math/rand"
	"time"
)

type EventName string

const (
	Asteroid           EventName = "Asteroid"
	DerelictShip                 = "Derelict Ship"
	LargeInvasionForce           = "Large Invasion Force"
	PeaceAndQuiet                = "Peace & Quiet"
	Revolt                       = "Revolt"
	SmallInvasionForce           = "Small Invasion Force"
	Strike                       = "Strike"
)

type EventModifier string

const (
	InvasionMod EventModifier = "Add +1 Resistance With Planetary Defenses"
	RevoltMod                 = "+1 System Resistance With Hyper Television"
	StrikeMod                 = "With Robot Workers, gain 1/2 instead of zero (round up)"
)

type EventCard struct {
	ID          string
	Name        EventName
	Year1Effect string
	Year2Effect string
	Modifier    EventModifier
}

var events = []EventCard{
	{"1", Asteroid, "+1 Wealth", "+1 Wealth", ""},
	{"2", DerelictShip, "Gain 1 Metal", "Gain 1 Metal", ""},
	{"3", LargeInvasionForce, "Force +2", "Force +3", InvasionMod},
	{"4", PeaceAndQuiet, "No event", "No event", ""},
	{"5", Revolt, "Force +1", "Force +2", RevoltMod},
	{"6", Revolt, "Force +1", "Force +3", RevoltMod},
	{"7", SmallInvasionForce, "Force +1", "Force +2", InvasionMod},
	{"8", Strike, "No resources next turn", "No resources next turn", StrikeMod},
}

type SystemType string

const (
	StartingSystem SystemType = "Starting System"
	NearSystem                = "Near System"
	DistantSystem             = "DistantSystem"
)

type SystemName string

type SystemCard struct {
	ID         string
	Name       SystemName
	Type       SystemType
	Resistance int
	Metal      int
	Wealth     int
	VPs        int
	Revolted   bool
}

var systems = []SystemCard{
	{ID: "1", Name: "Home World", Metal: 1, Wealth: 1},
	{ID: "2", Name: "Cygnus", Resistance: 5, Wealth: 1, VPs: 1},
	{ID: "3", Name: "Epsilon Eridani", Resistance: 8, VPs: 1},
	{ID: "4", Name: "Procyon", Resistance: 7, Wealth: 1, VPs: 1},
	{ID: "5", Name: "Proxima", Resistance: 6, Metal: 1, VPs: 1},
	{ID: "6", Name: "Sirius", Resistance: 6, VPs: 1},
	{ID: "7", Name: "Wolf 359", Resistance: 5, Metal: 1, VPs: 1},
	{ID: "8", Name: "Tau Ceti", Resistance: 4, VPs: 1},
	{ID: "9", Name: "Canopus", Resistance: 9, Wealth: 1, VPs: 2},
	{ID: "10", Name: "Galaxy's Edge", Resistance: 10, VPs: 3},
	{ID: "11", Name: "Polaris", Resistance: 9, VPs: 2},
}

var (
	Systems           map[string]*SystemCard
	Events            map[string]*EventCard
	EventDeck         []string
	NearSystemDeck    []string
	DistantSystemDeck []string
)

type Deck []string

func init() {
	rand.Seed(time.Now().UnixNano())

	Systems = make(map[string]*SystemCard)
	for i := range systems {
		c := &systems[i]
		switch {
		case c.ID == "1":
			c.Type = StartingSystem
		case c.ID == "10" || c.ID == "11" || c.ID == "12":
			c.Type = DistantSystem
		default:
			c.Type = NearSystem
		}
		Systems[c.ID] = c
	}

	Events = make(map[string]*EventCard)
	for i := range events {
		Events[events[i].ID] = &events[i]
	}

}

func shuffle(deck []string) {
	for i := range deck {
		n := len(deck) - i
		k := rand.Intn(n)
		deck[i], deck[i+k] = deck[i+k], deck[i]
	}
}

func Draw(deck []string) (string, []string) {
	card := deck[0]
	return card, deck[1:]
}

func Roll() int {
	return rand.Intn(6) + 1
}
