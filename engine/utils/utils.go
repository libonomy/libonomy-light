package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/libonomy/libonomy-gota/dataframe"
	"github.com/libonomy/libonomy-light/engine/models"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

//Sigmoid Activation Function
func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
	// if x < 0 {
	// 	return 0.0
	// }
	// return x
}

//SigmoidPrime Derivative of Sigmoid Function
func SigmoidPrime(x float64) float64 {
	return x * (1.0 - x)
	// if x < 0 {
	// 	return 0.0
	// }

	// return 1.0
}

//NewNetwork for obtaining new network
func NewNetwork(Config models.NeuralNetConfig) *models.NeuralNet {
	return &models.NeuralNet{Config: Config}
}

//SumAlongAxis function for sum
func SumAlongAxis(axis int, m *mat.Dense) (*mat.Dense, error) {
	numRows, numCols := m.Dims()
	var output *mat.Dense
	switch axis {
	case 0:
		data := make([]float64, numCols)
		for i := 0; i < numCols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		output = mat.NewDense(1, numCols, data)
	case 1:
		data := make([]float64, numRows)
		for i := 0; i < numRows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		output = mat.NewDense(numRows, 1, data)
	default:
		return nil, errors.New("invalid axis, must be 0 or 1")
	}
	return output, nil
}

//SplitData function for spliting data set
func SplitData(data dataframe.DataFrame) {

	trainingNum := (data.Nrow() * 4) / 5
	testNum := data.Nrow() / 5

	fmt.Println("Training ", trainingNum, "Testing", testNum)
	trainingIndex := make([]int, trainingNum)
	testingIndex := make([]int, testNum)

	for i := 0; i < trainingNum; i++ {
		trainingIndex[i] = i
	}
	for i := 0; i < testNum; i++ {
		testingIndex[i] = trainingNum + i
	}

	trainingDF := data.Subset(trainingIndex)
	testDF := data.Subset(testingIndex)

	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	for i, nameDF := range []string{"iris_train.csv", "iris_test.csv"} {
		f, err := os.Create(nameDF)
		if err != nil {
			log.Fatal(err)
		}

		// Create a buffered writer.
		w := bufio.NewWriter(f)

		setMap[i].WriteCSV(w)
	}

}

//ShuffleRawCSVdata function for shuffling dataset
func ShuffleRawCSVdata(data [][]string) [][]string {

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	return data
}

//Train functions
func Train(x, y *mat.Dense, nn *models.NeuralNet) (*mat.Dense, error) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	wHiddenRaw := make([]float64, nn.Config.HiddenNeurons*nn.Config.InputNeurons)
	bHiddenRaw := make([]float64, nn.Config.HiddenNeurons)
	wOutRaw := make([]float64, nn.Config.HiddenNeurons*nn.Config.OutputNeurons)
	bOutRaw := make([]float64, nn.Config.OutputNeurons)

	for _, param := range [][]float64{wHiddenRaw, bHiddenRaw, wOutRaw, bOutRaw} {
		for i := range param {
			param[i] = randGen.Float64()
			//fmt.Println(param[i])
		}
	}

	wHidden := mat.NewDense(nn.Config.InputNeurons, nn.Config.HiddenNeurons, wHiddenRaw)
	bHidden := mat.NewDense(1, nn.Config.HiddenNeurons, bHiddenRaw)
	wOut := mat.NewDense(nn.Config.HiddenNeurons, nn.Config.OutputNeurons, wOutRaw)
	bOut := mat.NewDense(1, nn.Config.OutputNeurons, bOutRaw)
	// fmt.Println("Code Reached Here")
	var output mat.Dense

	for i := 0; i < nn.Config.NumEpochs; i++ {
		// Feed Forward Process
		var hiddenLayerInput mat.Dense
		//hiddenLayerInput := mat.NewDense(0, 0, nil)
		hiddenLayerInput.Mul(x, wHidden)
		addBHidden := func(_, col int, v float64) float64 { return v + bHidden.At(0, col) }
		hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)
		//fmt.Println("Code Reached Here", i)
		// hiddenLayerActivations := mat.NewDense(0, 0, nil)
		var hiddenLayerActivations mat.Dense
		applySigmoid := func(_, _ int, v float64) float64 { return Sigmoid(v) }
		hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

		//outputLayerInput := mat.NewDense(0, 0, nil)
		var outputLayerInput mat.Dense
		outputLayerInput.Mul(&hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 { return v + bOut.At(0, col) }
		outputLayerInput.Apply(addBOut, &outputLayerInput)
		output.Apply(applySigmoid, &outputLayerInput)

		// Back Propogation Process
		// networkError := mat.NewDense(0, 0, nil)
		var networkError mat.Dense
		networkError.Sub(y, &output)

		// slopeOutputLayer := mat.NewDense(0, 0, nil)
		var slopeOutputLayer mat.Dense
		applySigmoidPrime := func(_, _ int, v float64) float64 { return SigmoidPrime(v) }
		slopeOutputLayer.Apply(applySigmoidPrime, &output)
		// slopeHiddenLayer := mat.NewDense(0, 0, nil)
		var slopeHiddenLayer mat.Dense
		slopeHiddenLayer.Apply(applySigmoidPrime, &hiddenLayerActivations)

		// dOutput := mat.NewDense(0, 0, nil)
		var dOutput mat.Dense
		dOutput.MulElem(&networkError, &slopeOutputLayer)
		// errorAtHiddenLayer := mat.NewDense(0, 0, nil)
		var errorAtHiddenLayer mat.Dense
		errorAtHiddenLayer.Mul(&dOutput, wOut.T())

		// dHiddenLayer := mat.NewDense(0, 0, nil)
		var dHiddenLayer mat.Dense
		dHiddenLayer.MulElem(&errorAtHiddenLayer, &slopeHiddenLayer)

		// Adjust The Parameters
		// wOutAdj := mat.NewDense(0, 0, nil)
		var wOutAdj mat.Dense
		wOutAdj.Mul(hiddenLayerActivations.T(), &dOutput)
		wOutAdj.Scale(nn.Config.LearningRate, &wOutAdj)
		wOut.Add(wOut, &wOutAdj)

		bOutAdj, err := SumAlongAxis(0, &dOutput)
		if err != nil {
			return nil, err
		}
		bOutAdj.Scale(nn.Config.LearningRate, bOutAdj)
		bOut.Add(bOut, bOutAdj)

		// wHiddenAdj := mat.NewDense(0, 0, nil)
		var wHiddenAdj mat.Dense
		wHiddenAdj.Mul(x.T(), &dHiddenLayer)
		wHiddenAdj.Scale(nn.Config.LearningRate, &wHiddenAdj)
		wHidden.Add(wHidden, &wHiddenAdj)

		bHiddenAdj, err := SumAlongAxis(0, &dHiddenLayer)
		if err != nil {
			return nil, err
		}
		bHiddenAdj.Scale(nn.Config.LearningRate, bHiddenAdj)
		bHidden.Add(bHidden, bHiddenAdj)

	}

	nn.WHidden = wHidden
	nn.BHidden = bHidden
	nn.WOut = wOut
	nn.BOut = bOut

	return &output, nil
}

//Predict for prediction
func Predict(x *mat.Dense, nn *models.NeuralNet) (*mat.Dense, error) {

	// Check to make sure that our neuralNet value
	// represents a trained model.
	if nn.WHidden == nil || nn.WOut == nil || nn.BHidden == nil || nn.BOut == nil {
		return nil, errors.New("the supplied neurnal net weights and biases are empty")
	}

	// Define the output of the neural network.
	var output mat.Dense

	// Complete the feed forward process.
	var hiddenLayerInput mat.Dense
	hiddenLayerInput.Mul(x, nn.WHidden)
	addBHidden := func(_, col int, v float64) float64 { return v + nn.BHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)

	var hiddenLayerActivations mat.Dense
	applySigmoid := func(_, _ int, v float64) float64 { return Sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

	var outputLayerInput mat.Dense
	outputLayerInput.Mul(&hiddenLayerActivations, nn.WOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.BOut.At(0, col) }
	outputLayerInput.Apply(addBOut, &outputLayerInput)
	output.Apply(applySigmoid, &outputLayerInput)

	return &output, nil
}
