package q_star

import (
	"fmt"
	"math/rand"
	"time"
)

// State is an interface wrapping the current state of the model.
type State interface {

	// String returns a string representation of the given state.
	// Implementers should take care to insure that this is a consistent
	// hash for a given state.
	String() string

	// Next provides a slice of possible Actions that could be applied to
	// a state.
	Next() []Action
}

// Action is an interface wrapping an action that can be applied to the
// model's current state.
//
// BUG (ecooper): A state should apply an action, not the other way
// around.
type Action interface {
	Float() float64
	Apply(State) State
}

// Rewarder is an interface wrapping the ability to provide a reward
// for the execution of an action in a given state.
type Rewarder interface {
	// Reward calculates the reward value for a given action in a given
	// state.
	Reward(*StateAction, float64) float64
}

// Agent is an interface for a model's agent and is able to learn
// from actions and return the current Q-value of an action at a given state.
type Agent interface {
	// Learn updates the model for a given state and action, using the
	// provided Rewarder implementation.
	Learn(*StateAction, Rewarder, float64)

	// Value returns the current Q-value for a State and Action.
	Value(State, Action) float64

	// Return a string representation of the Agent.
	String() string
}

// StateAction is a struct grouping an action to a given State. Additionally,
// a Value can be associated to StateAction, which is typically the Q-value.
type StateAction struct {
	State  State
	Action Action
	Value  float64
}

// NewStateAction creates a new StateAction for a State and Action.
func NewStateAction(state State, action Action, val float64) *StateAction {
	return &StateAction{
		State:  state,
		Action: action,
		Value:  val,
	}
}

// Next uses an Agent and State to find the highest scored Action.
//
// In the case of Q-value ties for a set of actions, a random
// value is selected.
func Next(agent Agent, state State) *StateAction {
	best := make([]*StateAction, 0)
	bestVal := 0.0

	for _, action := range state.Next() {
		val := agent.Value(state, action)

		if bestVal == 0.0 {
			best = append(best, NewStateAction(state, action, val))
			bestVal = val
		} else {
			if val > bestVal {
				best = []*StateAction{NewStateAction(state, action, val)}
				bestVal = val
			} else if val == bestVal {
				best = append(best, NewStateAction(state, action, val))
			}
		}
	}

	return best[rand.Intn(len(best))]
}

// QStar is an Agent implementation that stores Q-values in a
// map of maps.
type QStar struct {
	q       map[string]map[float64]float64 `json:"q"`
	lr      float64                        `json:"lr"`
	d       float64                        `json:"d"`
	mapping map[float64]float64            `json:"mapping"`
	lives   int                            `json:"lives"`
}

// NewQStar creates a QStar with the provided learning rate
// and discount factor.
func NewQStar(lr, d float64, mapping map[float64]float64, lives int) *QStar {
	return &QStar{
		q:       make(map[string]map[float64]float64),
		d:       d,
		lr:      lr,
		mapping: mapping,
		lives:   lives,
	}
}

// getActions returns the current Q-values for a given state.
func (agent *QStar) getActions(state string) map[float64]float64 {
	if _, ok := agent.q[state]; !ok {
		agent.q[state] = make(map[float64]float64)
	}

	return agent.q[state]
}

// Learn updates the existing Q-value for the given State and Action
// using the Rewarder.
//
// See https://en.wikipedia.org/wiki/Q-learning#Algorithm
func (agent *QStar) Learn(action *StateAction, reward Rewarder, signal float64) {
	current := action.State.String()
	next := action.Action.Apply(action.State).String()

	actions := agent.getActions(current)

	maxNextVal := 0.0
	for _, v := range agent.getActions(next) {
		if v > maxNextVal {
			maxNextVal = v
		}
	}

	currentVal := actions[action.Action.Float()]
	actions[action.Action.Float()] = currentVal + agent.lr*(reward.Reward(action, signal)+agent.d*maxNextVal-currentVal)
}

// Value gets the current Q-value for a State and Action.
func (agent *QStar) Value(state State, action Action) float64 {
	return agent.getActions(state.String())[action.Float()]
}

func (agent *QStar) Predict(state []float64) (prediction []float64) {
	game := NewGame(state, false, agent.mapping, agent.lives)
	action := Next(agent, game)

	return []float64{action.Action.Float()}
}

// String returns the current Q-value map as a printed string.
//
// BUG (ecooper): This is useless.
func (agent *QStar) String() string {
	return fmt.Sprintf("%v", agent.q)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
