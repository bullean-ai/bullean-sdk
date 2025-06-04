package domain

type IModel interface {
	Predict([]float64) []float64
}

type ITrainer interface {
	Train(interface{}, Examples, Examples, int) (float64, IModel)
}

type IEvaluator interface {
	Evaluate([]float64, []float64) float64
}
