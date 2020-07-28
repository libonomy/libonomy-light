package models

import "gonum.org/v1/gonum/mat"

//NeuralNet structure
type NeuralNet struct {
	Config  NeuralNetConfig
	WHidden *mat.Dense
	BHidden *mat.Dense
	WOut    *mat.Dense
	BOut    *mat.Dense
}

//NeuralNetConfig structure
type NeuralNetConfig struct {
	InputNeurons  int
	OutputNeurons int
	HiddenNeurons int
	NumEpochs     int
	LearningRate  float64
}
