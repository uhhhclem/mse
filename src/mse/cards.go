package mse

import (
    "math/rand"
    "time"
)

type CardName string

const (
    Asteroid = "Asteroid"
    DerelictShip CardName = "Derelict Ship"
    LargeInvasionForce = "Large Invasion Force"
    PeaceAndQuiet = "Peace & Quiet"
    Revolt CardName = "Revolt"
    SmallInvasionForce CardName = "Small Invasion Force"
    Strike CardName = "Strike"
)

type EventModifier string

const (
    InvasionMod EventModifier = "Add +1 Resistance With Planetary Defenses"
    RevoltMod = "+1 System Resistance With Hyper Television"
    StrikeMod = "With Robot Workerrs, gain 1/2 instead of zero (round up)"
)

type EventCard struct {
    ID int
    Name CardName
    Year1Effect string
    Year2Effect string
    Modifier EventModifier
}

var events = []EventCard {
    {1, Asteroid, "+1 Wealth", "+1 Wealth", ""},
    {2, DerelictShip, "Gain 1 Metal", "Gain 1 Metal", ""},
    {3, LargeInvasionForce, "Force +2", "Force +3", InvasionMod},
    {4, PeaceAndQuiet, "No event", "No event", ""},
    {5, Revolt, "Force +1", "Force +2", RevoltMod},
    {6, Revolt, "Force +1", "Force +3", RevoltMod},
    {7, SmallInvasionForce, "Force +1", "Force +2", InvasionMod},
    {8, Strike, "No resources next turn", "No resources next turn", StrikeMod},
}

type SystemType string

const (
    StartingSystem SystemType = "Starting System"
    NearSystem = "Near System"
    DistantSystem = "DistantSystem"
)

type SystemCard struct {
    ID int
    Name CardName
    Type SystemType
    Resistance int
    Metal int
    Wealth int
    VPs int
}

var systems = []SystemCard {
    {ID: 1, Name: "Home World", Metal: 1, Wealth: 1},
    {ID: 2, Name: "Cygnus", Resistance: 5, Wealth: 1, VPs: 1},
    {ID: 3, Name: "Epsilon Eridani", Resistance: 8, VPs: 1},
    {ID: 4, Name: "Procyon", Resistance: 7, Wealth: 1, VPs: 1},
    {ID: 5, Name: "Proxima", Resistance: 6, Metal: 1, VPs: 1},
    {ID: 6, Name: "Sirius", Resistance: 6, VPs: 1},
    {ID: 7, Name: "Wolf 359", Resistance: 5, Metal: 1, VPs: 1},
    {ID: 8, Name: "Tau Ceti", Resistance: 4, VPs: 1},
    {ID: 9, Name: "Canopus", Resistance: 9, Wealth: 1, VPs: 2},
    {ID: 10, Name: "Galaxy's Edge", Resistance: 10, VPs: 3},
    {ID: 11, Name: "Polaris", Resistance: 9, VPs: 2},
}

var (
    Systems map[int]*SystemCard
     Events map[int]*EventCard
    EventDeck []int
    NearSystemDeck []int
    DistantSystemDeck []int
)

func init() {
    rand.Seed(time.Now().UnixNano())
    
    for _, c := range Systems {
        switch {
        case c.ID == 1:
            c.Type = StartingSystem
        case c.ID > 1 && c.ID < 9:
            c.Type = NearSystem
        case c.ID > 9:
            c.Type = DistantSystem
        }
    }
    
    Events = make(map[int]*EventCard)
    for i := range events {
        Events[events[i].ID] = &events[i]
    }
    
    Systems = make(map[int]*SystemCard)
    for i := range systems {
        Systems[systems[i].ID] = &systems[i]
    }
    
    EventDeck = []int{1, 2, 3, 4, 5, 6, 7, 8}
    NearSystemDeck = []int{2, 3, 4, 5, 6, 7, 8}
    DistantSystemDeck = []int{9, 10, 11}
    
    shuffle(EventDeck)
    shuffle(NearSystemDeck)
    shuffle(DistantSystemDeck)
}

func shuffle(deck []int) {
    for i := range deck {
        n := len(deck) - i
        k := rand.Intn(n)
        deck[i], deck[i + k] = deck[i + k], deck[i]
    }
}

func Draw(deck []int) (int, []int) {
    card := deck[0]
    return card, deck[1:]
}