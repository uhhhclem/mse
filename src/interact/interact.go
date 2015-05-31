// Package interact provides services for interacting with the game.
package interact

import (
	"fmt"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

type GameState string

// Game contains the common structures needed to support game interaction.
type Game struct {
	// ID uniquely identifies the game object.
	ID string
	// State identifies the games' current state.
	State GameState
	// Prompt contains the current prompt while it's under construction
	// (via the NewPrompt and AddChoice methods), and is also used to validate
	// the choice retrieved from NextChoice.
	Prompt *Prompt
	// Games use Ready to signal to the client that the game has
	// been updated.
	Ready chan bool
	// NextPrompt contains the next prompt (including valid choices)
	// available to the user.
	NextPrompt chan *Prompt
	// NextStatus contains the next status message available for the client
	// to retrieve.
	NextStatus chan *Status
	// NextChoice contains the next choice made by the player in response to
	// a prompt.
	NextChoice chan *Choice
}

// NewGame returns a new Game object with all channels initialized.
func NewGame() *Game {
	return &Game{
		ID:         uuid.New(),
		NextStatus: make(chan *Status),
		NextPrompt: make(chan *Prompt),
		NextChoice: make(chan *Choice),
		Ready:      make(chan bool),
	}
}

// Prompt represents a multiple-choice prompt to the player.
type Prompt struct {
	State   GameState
	Message string
	Choices []*Choice
}

// NewPrompt initialized the current prompt for the game; follow this with
// calls to AddChoice to populate the prompt's choices.
func (g *Game) NewPrompt(msg string) {
	g.Prompt = &Prompt{
		State:   g.State,
		Message: msg,
		Choices: make([]*Choice, 0),
	}
}

// AddChoice adds a new choice to the current prompt.
func (g *Game) AddChoice(key, name string) {
	g.Prompt.Choices = append(g.Prompt.Choices, &Choice{key, name})
}

// Choice represents a choice available to the player at the current prompt.
type Choice struct {
	Key  string
	Name string
}

// Status contains messages logged to the player via *game.Log() and .Logf().
type Status struct {
	Message string
}

// SendPrompt makes the current prompt available to the client.
func (g *Game) SendPrompt() {
	func() { g.NextPrompt <- g.Prompt }()
}

// Log sends a Status message to the player.
func (g *Game) Log(m string) {
	go func() { g.NextStatus <- &Status{m} }()
}

// Logf sends a formatted Status message to the player.
func (g *Game) Logf(f string, args ...interface{}) {
	g.Log(fmt.Sprintf(f, args...))
}

// MakeChoice validates the player's choice and puts it in the NextChoice
// channel.
func (g *Game) MakeChoice(key string) error {
	for _, c := range g.Prompt.Choices {
		if strings.ToLower(key) == strings.ToLower(c.Key) {
			go func() { g.NextChoice <- c }()
			return nil
		}
	}
	return fmt.Errorf("%q is not a valid choice.", key)
}
