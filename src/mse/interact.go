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

// SendPrompt blocks until the front-end picks up the Prompt.
func (g *Game) SendPrompt(p *Prompt) {
	p.State = g.State
	g.Prompt = p
	g.NextPrompt <- p
}

// Log blocks until the front-end picks up the Status.
func (g *Game) Log(m string) {
	g.NextStatus <- &Status{m}
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
