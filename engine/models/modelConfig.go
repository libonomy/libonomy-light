package models

//ModelConfig structure for saving model details
type ModelConfig struct {
	InputNeurons  int         `json:"inputNeurons"`
	OutputNeurons int         `json:"outputNeurons"`
	HiddenNeurons int         `json:"hiddenNeurons"`
	NumEpochs     int         `json:"numEpochs"`
	LearningRate  float64     `json:"learningRate"`
	BHidden       [][]float64 `json:"bHidden"`
	BHiddenDims   []int       `json:"bHiddenDims"`
	WHidden       [][]float64 `json:"wHidden"`
	WHiddenDims   []int       `json:"wHiddenDims"`
	BOut          [][]float64 `json:"bOut"`
	BOutDims      []int       `json:"bOutDims"`
	WOut          [][]float64 `json:"wOut"`
	WOutDims      []int       `json:"wOutDims"`
}
