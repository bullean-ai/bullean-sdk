package q_star

import (
	"fmt"
)

const (
	startingLives = 6

	Lost   = 3
	Active = 1
	Won    = 2

	Buy  = 1
	Sell = -1
)

var (
	debug      bool = false
	progressAt int  = 0
)

// Game represents the state of any given game of Hangman. It implements
// Agent, Rewarder, and State.
type Game struct {
	State     []float64
	Lives     int
	Attempted map[float64]bool
	Positions map[float64]float64
	debug     bool
}

// NewGame creates a new Hangman game for the given word. If debug
// is true, Game.Log messages will print to stdout.
func NewGame(state []float64, debug bool, positions map[float64]float64, lives int) *Game {
	game := &Game{debug: debug}
	game.New(state, positions, lives)

	return game
}

// New resets the current game to a new game for the given word.
func (game *Game) New(states []float64, positions map[float64]float64, lives int) {
	game.State = states
	game.Lives = lives
	game.Attempted = make(map[float64]bool)
	game.Positions = positions
}

// Returns Lost, Active, or Won based on the game's current state.
func (game *Game) IsComplete() int {
	/*
		if game.Lives > 0 {
			return Active
		} else if game.Profit > 0 {
			return Won
		} else if game.Profit <= 0 {
			return Lost
		}
	*/
	return 0
}

// Choose applies a character attempt in the current game, returning
// true if char is present in Game.Word.
//
// Choose updates the s's state.
func (game *Game) Choose(position float64) bool {
	game.Attempted[position] = true

	hit := false

	for _, key := range game.State {
		if key == position {
			game.Lives -= 1
			hit = true
		}
	}

	if !hit {
		return false
	}

	return true
}

// Reward returns a score for a given StateAction. Reward is a
// member of the Rewarder interface. If the choice is found in
// the game's word, a positive score is returned. Otherwise, a static
// -1000 is returned.
func (game *Game) Reward(action *StateAction, signal float64) float64 {
	choice := action.Action.Float()

	if signal == choice {
		return 10
	} else {
		return -10
	}

}

// Next creates a new slice of Action instances. A possible
// action is created for each character that has not been attempted in
// in the game.
func (game *Game) Next() []Action {
	actions := make([]Action, 0, len(game.State))

	for check, _ := range game.Positions {
		attempted := game.Attempted[check]
		if !attempted {
			actions = append(actions, &Choice{Position: check})
		}
	}

	return actions
}

// Log is a wrapper of fmt.Printf. If Game.debug is true, Log will print
// to stdout.
func (game *Game) Log(msg string, args ...interface{}) {
	if game.debug {
		logMsg := fmt.Sprintf("[GAME %v] (%d moves, %d lives) %s\n", game.State, len(game.Attempted), game.Lives, msg)
		fmt.Printf(logMsg, args...)
	}
}

// String returns a consistent hash for the current game state to be
// used in a Agent.
func (game *Game) String() string {
	return fmt.Sprintf("%v", game.State)
}

// Choice implements Action for a character choice in a game
// of Hangman.
type Choice struct {
	Position float64
}

// String returns the character for the current action.
func (choice *Choice) Float() float64 {
	return choice.Position
}

// String returns the character for the current action.
func (choice *Choice) String() string {
	return fmt.Sprintf("%v", choice.Position)
}

// Apply updates the state of the game for a given character choice.
func (choice *Choice) Apply(state State) State {
	game := state.(*Game)
	game.Choose(choice.Position)

	return game
}
