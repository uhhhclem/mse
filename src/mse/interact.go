package mse

import (
	"fmt"
	"strings"
)

type Prompt struct {
	State   GameState
	Message string
	Choices []*Choice
}

type Choice struct {
	Key  string
	Name string
}

type Status struct {
	Message string
}

func (g *Game) SendPrompt(message string, choices []*Choice) {
	g.Prompt = &Prompt{
		State:   g.State,
		Message: message,
		Choices: choices,
	}
	go func() { g.NextPrompt <- g.Prompt }()
}

func (g *Game) Log(m string) {
	go func() { g.NextStatus <- &Status{m} }()
}

// Logf sends a formatted Status message.
func (g *Game) Logf(f string, args ...interface{}) {
	g.Log(fmt.Sprintf(f, args...))
}

func (g *Game) MakeChoice(key string) error {
	for _, c := range g.Prompt.Choices {
		if strings.ToLower(key) == strings.ToLower(c.Key) {
			go func() { g.NextChoice <- c }()
			return nil
		}
	}
	return fmt.Errorf("%q is not a valid choice.", key)
}
