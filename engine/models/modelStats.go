package models

//ModelStats structure for model statistics
type ModelStats struct {
	TestAccuracy    float64 `json:"testAccuracy"`
	TrainAccuracy   float64 `json:"trainAccuracy"`
	PredictAccuracy float64 `json:"predictAccuracy"`
}
