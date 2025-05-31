package ffnn

import (
	"github.com/bullean-ai/bullean-sdk/data/domain"
	"github.com/bullean-ai/bullean-sdk/data/neural/ffnn/layer/neuron/synapse"
)

var DefaultConfig = func(input_len int) *domain.Config {
	return &domain.Config{
		Inputs:     input_len,
		Layout:     []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 2},
		Activation: domain.ActivationSigmoid,
		Mode:       domain.ModeMultiClass,
		Weight:     synapse.NewNormal(1e-15, 1e-15),
		Bias:       true,
	}
}
