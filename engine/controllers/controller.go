package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"

	"github.com/libonomy/libonomy-gota/dataframe"
	"github.com/libonomy/libonomy-light/engine/dto"
	"github.com/libonomy/libonomy-light/engine/models"
	"github.com/libonomy/libonomy-light/engine/utils"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

type bodyVariables struct {
	ComputerPower float64 `json:"computerPower"`
	DownloadSpeed float64 `json:"downSpeed"`
	Ylabels       string  `json:"yLabels"`
}

type check1 struct {
	Name string
}

type check2 struct {
	Name string
}

type check3 struct {
	Name string
}

//Testing Function To Test its working
func Testing(w http.ResponseWriter, r *http.Request) {
	s := "This is for testing function only"

	dto.SendResponse(w, r, http.StatusOK, "Success", map[string]interface{}{"testing": s})
}

//GenerateCSV function to generate csv file from json data.
func GenerateCSV(w http.ResponseWriter, r *http.Request) {
	variables := []bodyVariables{}
	fileVariables := []bodyVariables{}
	filename := "./datasets/testing/dummyDataset.json"
	csvFilename := "./datasets/testing/dummyDataset.csv"
	err := json.NewDecoder(r.Body).Decode(&variables)
	if err != nil {
		fmt.Println("There is some Error in Decoding Body Request", err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Bad", map[string]interface{}{"Error": err.Error()})
		return
	}

	//fmt.Println(variables)
	_, err = os.Stat(filename)

	if err != nil {
		//checking if file does not exists
		fmt.Println("Error is os.stat", err.Error())
		//fmt.Println("File Info ", fileInfo)

		jsonFile, _ := os.Create(filename)
		csvFile, _ := os.Create(csvFilename)

		defer jsonFile.Close()
		file, _ := json.MarshalIndent(variables, "", " ")
		_ = ioutil.WriteFile(filename, file, 0644)

		fileOpened, _ := os.Open(filename)
		dataFrame := dataframe.ReadJSON(fileOpened)
		csvWriter := csv.NewWriter(csvFile)
		csvWriter.WriteAll(dataFrame.Records())
		// csvData := dataframe.ReadJSON(jsonFile)
		// csvData.WriteCSV(csvFile)
		dto.SendResponse(w, r, http.StatusOK, "Successfully Created CSV File From JSON", nil)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Cannot create file", nil)
	}
	fileBytes, _ := ioutil.ReadAll(file)
	fmt.Println("Code is Here 0")
	json.Unmarshal(fileBytes, &fileVariables)
	fmt.Println("Code is Here 1")

	fileVariables = append(fileVariables, variables[0])
	writeF, _ := json.MarshalIndent(fileVariables, "", " ")
	_ = ioutil.WriteFile(filename, writeF, 0644)

	csvFile, _ := os.Create(csvFilename)
	fileOpened, _ := os.Open(filename)
	dataFrame := dataframe.ReadJSON(fileOpened)
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.WriteAll(dataFrame.Records())

	dto.SendResponse(w, r, http.StatusOK, "Successfully Created CSV File From JSON", nil)
}

//CleanData function to clean data and convert data to specific format
func CleanData(w http.ResponseWriter, r *http.Request) {
	filename := "./datasets/testing/dummyDataset.csv"
	fmt.Println("code Here")
	f, err := os.Open(filename)
	newFile := f
	defer newFile.Close()
	if err != nil {
		fmt.Println("Error in opening file", f.Name(), err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in opening file "+filename+" Error is \n"+err.Error(), nil)
		return
	}
	fmt.Println("code Here")
	reader := csv.NewReader(f)
	rawData, err := reader.ReadAll()
	if err != nil {
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in opening file 2"+err.Error(), nil)
		return
	}

	data := dataframe.LoadRecords(rawData)
	// fmt.Println(data)

	// check := []int{1, 1, 1, 1, 1, 1}

	labelExtraction := []string{}

	labels := data.Col("yLabels").Records()

	var index int
	dataRecords := data.Records()
	for i, record := range dataRecords[0] {
		// fmt.Println("Index ", i, "Value", record)
		if record == "yLabels" {
			index = i
		}

	}

	for _, label := range labels {
		// fmt.Println("Index ", i, "Label", label)
		sort.Strings(labelExtraction)
		res := sort.SearchStrings(labelExtraction, label)

		if len(labelExtraction) == res {
			labelExtraction = append(labelExtraction, label)
		}
		// fmt.Println("Printing Result ", res)
	}
	// fmt.Println(rawData)

	var newRecords [][]string
	// fmt.Println(labelExtraction)
	length := len(rawData[0])
	for i, record := range rawData {

		// fmt.Println("Length", length)
		if i == 0 {
			for _, label := range labelExtraction {
				record = append(record, label)
			}
			newRecords = append(newRecords, record)

			continue
		}
		for _, label := range labelExtraction {
			label = label
			// fmt.Println(label)
			record = append(record, "0.0")

		}

		for index, label := range labelExtraction {
			if record[length-1] == label {
				record[index+length] = "1.0"
				// fmt.Println("Yes in here")
			}

			// fmt.Println("Record", record[len(record)-(length)])
		}
		newRecords = append(newRecords, record)

	}

	// fmt.Println(newRecords)
	var finalRecords [][]string
	for _, record := range newRecords {
		// fmt.Println(record[:length-1])
		// fmt.Println(record[length:])

		modified := append(record[:length-1], record[length:]...)
		finalRecords = append(finalRecords, modified)
		// finalRecords = append(finalRecords, record[length:])
	}
	fmt.Println(finalRecords[0])

	writer, _ := os.Create("./datasets/dataset.csv")
	wr := csv.NewWriter(writer)
	wr.WriteAll(finalRecords)
	wr.Flush()
	dto.SendResponse(w, r, http.StatusOK, "Success", map[string]interface{}{"raw csv data": finalRecords, "Index ": index, "Length": len(newRecords)})
}

//NormalizeData function for data normalization from 0-1
func NormalizeData(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./datasets/dataset.csv")

	if err != nil {
		fmt.Println("Code Is Here")

		dto.SendResponse(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	fmt.Println("Code Is Here")

	csvReader := csv.NewReader(f)
	rawCSVdata, _ := csvReader.ReadAll()
	dataFrame := dataframe.LoadRecords(rawCSVdata)
	normalize := [][]string{}
	// fmt.Println("Data Name ", dataFrame.Col(string(rawCSVdata[0][0])))
	for i, record := range rawCSVdata {
		if i == 0 {
			normalize = append(normalize, record)
			continue
		}
		noVal := []string{}
		for x, values := range record {
			val, _ := strconv.ParseFloat(values, 64)
			nVal := (val - dataFrame.Col(rawCSVdata[0][x]).Min()) / (dataFrame.Col(rawCSVdata[0][x]).Max() - dataFrame.Col(rawCSVdata[0][x]).Min())
			sVal := strconv.FormatFloat(nVal, 'f', 6, 64)
			noVal = append(noVal, sVal)
			//fmt.Println(nVal)
		}
		normalize = append(normalize, noVal)
	}

	headers := rawCSVdata[0]
	wr, _ := os.Create("./datasets/normalized.csv")
	csvWriter := csv.NewWriter(wr)
	csvWriter.WriteAll(normalize)
	dto.SendResponse(w, r, http.StatusOK, "Everthings Fine", map[string]interface{}{"data": rawCSVdata, "summary": dataFrame.Describe().Records(),
		"Headers": headers, "Normalixe": normalize, "Original": rawCSVdata})
}

//SplitAndShuffle to split and shuffle data
func SplitAndShuffle(w http.ResponseWriter, r *http.Request) {
	percentage := r.FormValue("trainPercentage")
	percent, _ := strconv.ParseFloat(percentage, 64)

	trainPercent := percent / 100
	testPercent := ((100 - percent) / 100) / 2
	var validatePercent int
	var even bool

	f, _ := os.Open("./datasets/normalized.csv")
	reader := csv.NewReader(f)
	rawCSV, _ := reader.ReadAll()
	fmt.Println(len(rawCSV))
	if int(float64(len(rawCSV))*testPercent)%2 == 0 {
		fmt.Println("Yes its even")
		even = true
	} else {
		fmt.Println("No ")
		even = false
	}

	fmt.Println(testPercent, even, validatePercent, trainPercent)
	dto.SendResponse(w, r, http.StatusOK, "Success", map[string]interface{}{})
}

//Train funciton
func Train(w http.ResponseWriter, r *http.Request) {
	rateString := r.URL.Query().Get("rate")
	epochsString := r.URL.Query().Get("epochs")
	hiddenString := r.URL.Query().Get("hidden")

	rate, _ := strconv.ParseFloat(rateString, 64)
	epochs, _ := strconv.ParseFloat(epochsString, 64)
	hidden, _ := strconv.ParseFloat(hiddenString, 64)

	fmt.Println("Rate is :\t", rate)
	fmt.Println("Type of Rate is :\t", reflect.TypeOf(rate))

	fmt.Println("Epochs is :\t", epochs)
	fmt.Println("hidden", hidden)

	f, err := os.Open("./datasets/iris_train.csv")

	if err != nil {
		// fmt.Println("Error in Reading File ", err.Error())
		return
	}
	reader := csv.NewReader(f)
	rawCSVdata, err := reader.ReadAll()
	inputsData := make([]float64, 4*len(rawCSVdata))
	labelsData := make([]float64, 3*len(rawCSVdata))
	// fmt.Println(len(rawCSVdata))
	// fmt.Println(len(inputsData))
	// fmt.Println(len(labelsData))

	var inputIdx int
	var labelIdx int

	for idx, record := range rawCSVdata {
		if idx == 0 {
			continue
		}

		for i, val := range record {
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				// fmt.Println("Error in Parsing Float Value", err.Error())
				return
			}

			if i == 4 || i == 5 || i == 6 {
				labelsData[labelIdx] = parsedVal
				labelIdx++
			} else {
				inputsData[inputIdx] = parsedVal
				inputIdx++
			}
		}
	}

	inputs := mat.NewDense(len(rawCSVdata), 4, inputsData)
	labels := mat.NewDense(len(rawCSVdata), 3, labelsData)

	config := models.NeuralNetConfig{
		InputNeurons:  4,
		OutputNeurons: 3,
		HiddenNeurons: 5,
		NumEpochs:     int(epochs),
		LearningRate:  rate,
	}

	network := utils.NewNetwork(config)
	// fmt.Println("Printing Network Before Training", network)

	trainOutput, err := utils.Train(inputs, labels, network)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Printing Network After Training", network)

	// fmt.Println(trainOutput)
	// fmt.Println(network.BHidden.RawMatrix())
	// fmt.Print("Printing B hidden \t")
	// fmt.Println(network.BHidden.Dims())
	rowsBHidden, colsBHidden := network.BHidden.Dims()

	// fmt.Print("Printing W hidden \t")
	// fmt.Println(network.WHidden.Dims())
	rowsWHidden, colsWHidden := network.WHidden.Dims()

	// fmt.Print("Printing W out  \t")
	// fmt.Println(network.WOut.Dims())
	rowsWOut, colsWOut := network.WOut.Dims()

	// fmt.Print("Printing B out \t")
	// fmt.Println(network.BOut.Dims())
	rowsBOut, colsBOut := network.BOut.Dims()

	configuration := models.ModelConfig{}

	for i := 0; i < rowsBHidden; i++ {
		configuration.BHidden = append(configuration.BHidden, network.BHidden.RawRowView(i))
	}
	for i := 0; i < rowsWHidden; i++ {
		configuration.WHidden = append(configuration.WHidden, network.WHidden.RawRowView(i))
	}
	for i := 0; i < rowsWOut; i++ {
		configuration.WOut = append(configuration.WOut, network.WOut.RawRowView(i))
	}
	for i := 0; i < rowsBOut; i++ {
		configuration.BOut = append(configuration.BOut, network.BOut.RawRowView(i))
	}

	configuration.BHiddenDims = append(configuration.BHiddenDims, rowsBHidden)
	configuration.BHiddenDims = append(configuration.BHiddenDims, colsBHidden)
	configuration.WHiddenDims = append(configuration.WHiddenDims, rowsWHidden)
	configuration.WHiddenDims = append(configuration.WHiddenDims, colsWHidden)
	configuration.WOutDims = append(configuration.WOutDims, rowsWOut)
	configuration.WOutDims = append(configuration.WOutDims, colsWOut)
	configuration.BOutDims = append(configuration.BOutDims, rowsBOut)
	configuration.BOutDims = append(configuration.BOutDims, colsBOut)
	configuration.InputNeurons = network.Config.InputNeurons
	configuration.HiddenNeurons = network.Config.HiddenNeurons
	configuration.LearningRate = network.Config.LearningRate
	configuration.NumEpochs = network.Config.NumEpochs
	configuration.OutputNeurons = network.Config.OutputNeurons

	// fmt.Println("Configuration Values ", configuration)

	var truePosNeg int
	numPreds, _ := trainOutput.Dims()
	for i := 0; i < numPreds; i++ {
		// fmt.Println("Prediction Index ", i, "\t", trainOutput.RowView(i))
		// Get the label.
		labelRow := mat.Row(nil, i, labels)
		var species int
		for idx, label := range labelRow {
			if label == 1.0 {
				species = idx
				break
			}
		}

		// Accumulate the true positive/negative count.
		if trainOutput.At(i, species) == floats.Max(mat.Row(nil, i, trainOutput)) {
			truePosNeg++
		}
	}

	// Calculate the accuracy (subset accuracy).
	accuracy := float64(truePosNeg) / float64(numPreds)

	// Output the Accuracy value to standard out.
	fmt.Printf("\nAccuracy of Testing = %0.2f %%\n", accuracy*100)
	stats := models.ModelStats{}
	stats.TrainAccuracy = accuracy * 100

	file, _ := json.MarshalIndent(configuration, "", " ")

	_ = ioutil.WriteFile("./models/test.json", file, 0644)

	file, _ = json.MarshalIndent(stats, "", " ")

	_ = ioutil.WriteFile("./models/stats.json", file, 0644)

	dto.SendResponse(w, r, http.StatusOK, "Training Result ", map[string]interface{}{"output": configuration, "Accuracy": accuracy * 100})

}

//Predict Function for prediction
func Predict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	stringf1 := r.FormValue("f1")
	stringf2 := r.FormValue("f2")
	stringf3 := r.FormValue("f3")
	stringf4 := r.FormValue("f4")

	f1, _ := strconv.ParseFloat(stringf1, 64)
	f2, _ := strconv.ParseFloat(stringf2, 64)
	f3, _ := strconv.ParseFloat(stringf3, 64)
	f4, _ := strconv.ParseFloat(stringf4, 64)

	input := []float64{f1, f2, f3, f4}
	configuration := models.ModelConfig{}
	file, _ := ioutil.ReadFile("./models/test.json")
	_ = json.Unmarshal([]byte(file), &configuration)

	// fmt.Println("Printing Configurations", configuration)
	config := models.NeuralNetConfig{
		InputNeurons:  configuration.InputNeurons,
		OutputNeurons: configuration.OutputNeurons,
		HiddenNeurons: configuration.HiddenNeurons,
		NumEpochs:     configuration.NumEpochs,
		LearningRate:  configuration.LearningRate,
	}

	network := utils.NewNetwork(config)

	network.BHidden = mat.NewDense(configuration.BHiddenDims[0], configuration.BHiddenDims[1], nil)
	for i := 0; i < configuration.BHiddenDims[0]; i++ {
		fmt.Println(i)
		network.BHidden.SetRow(i, configuration.BHidden[i])
	}

	network.WHidden = mat.NewDense(configuration.WHiddenDims[0], configuration.WHiddenDims[1], nil)
	for i := 0; i < configuration.WHiddenDims[0]; i++ {
		network.WHidden.SetRow(i, configuration.WHidden[i])
	}

	network.BOut = mat.NewDense(configuration.BOutDims[0], configuration.BOutDims[1], nil)
	for i := 0; i < configuration.BOutDims[0]; i++ {
		network.BOut.SetRow(i, configuration.BOut[i])
	}

	network.WOut = mat.NewDense(configuration.WOutDims[0], configuration.WOutDims[1], nil)
	for i := 0; i < configuration.WOutDims[0]; i++ {
		network.WOut.SetRow(i, configuration.WOut[i])
	}

	// fmt.Println("Printing Configurations", configuration)
	features := mat.NewDense(1, 4, input)
	predictions, err := utils.Predict(features, network)

	if err != nil {
		fmt.Println("Error in something", err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in Prediction", map[string]interface{}{"Error": err.Error()})
	}

	fmt.Println(predictions)
	fmt.Println(floats.MaxIdx(mat.Row(nil, 0, predictions)))

	dto.SendResponse(w, r, http.StatusOK, "Success ", map[string]interface{}{"Class Name": models.ClassNames[floats.MaxIdx(mat.Row(nil, 0, predictions))], "Prediction": predictions.RawMatrix(), "Max Index": floats.MaxIdx(mat.Row(nil, 0, predictions)) + 1})

}
