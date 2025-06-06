package q_star

import (
	"fmt"
	"github.com/bullean-ai/bullean-go/neural_nets/domain"
)

type QStarTrainer struct {
	Mapping *map[float64]float64 `json:"mapping"`
}

func NewQStarTrainer(mapping map[float64]float64) *QStarTrainer {
	return &QStarTrainer{Mapping: &mapping}
}

func (t *QStarTrainer) Train(agent interface{}, train_data domain.Examples, validate_data domain.Examples, iterations int) (float64, domain.IModel) {

	var neural *QStar
	switch agent.(type) {
	case *QStar:
		neural = agent.(*QStar)
	default:
		fmt.Println("QStarTrainer can only train QSTAR")
		return 0, nil
	}
	var (
		wins     = 0
		lastWins = 0
		count    = 0
		// Our agent has a learning rate of 0.7 and discount of 1.0.
	)

	progressAt = len(train_data)
	progress := func() {
		// Print our progress every 1000 rows.
		if count > 0 && count%progressAt == 0 {
			rate := float32(wins-lastWins) / float32(progressAt) * 100.0
			lastWins = wins
			fmt.Printf("%d games played: %d WINS %d LOSSES %.0f%% WIN RATE\n", count, wins, count-wins, rate)
		}
	}

	for i := 0; i < iterations; i++ {
		for count = 0; count < len(train_data); count++ {
			signal := train_data[count].Response[0]
			// Get a new word and game for each iteration...
			game := NewGame(train_data[count].Input, debug, *t.Mapping, iterations) //TODO: tüm olası pozisyonlar ve hash tablosu entegre edilecek
			game.Log("Game created")

			// While the game is still active, we'll continue to update
			// our agent and learn from its choices.
			// Pick the next move, which is going to be a letter choice.
			action := Next(neural, game)

			// Whatever that choice is, let's update our model for its
			// impact. If the character chosen is in the game's word,
			// then this action will be positive. Otherwise, it will be
			// negative.
			neural.Learn(action, game, signal)

			// Reward doesn't change state so we can check what the
			// reward would be for this action, and report how the
			// game changed.
			if game.Reward(action, signal) > 0.0 {
				game.Log("%s was selected", action.Action.Float())
			}
			game.Log("%s was incorrect", action.Action.Float())

			progress()
		}
	}

	progress()

	fmt.Printf("\nAgent performance: %d games played, %d WINS %d LOSSES %.0f%% WIN RATE\n", count, wins, count-wins, float32(wins)/float32(count)*100.0)

	return 0, neural
}
