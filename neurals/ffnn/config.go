package ffnn

import (
	"github.com/bullean-ai/bullean-sdk/neurals/domain"
	"github.com/bullean-ai/bullean-sdk/neurals/ffnn/layer/neuron/synapse"
)

var DefaultConfig = func(input_len int) *domain.Config {
	return &domain.Config{
		Inputs:     input_len,
		Layout:     []int{30, 70, 70, 70, 70, 70, 70, 30, 2},
		Activation: domain.ActivationSigmoid,
		Mode:       domain.ModeMultiClass,
		Weight:     synapse.NewNormal(1e-20, 1e-20),
		Bias:       true,
	}
}
