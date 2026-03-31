package alternative

type FngResponse struct {
	Name     string    `json:"name"`
	Data     []FngData `json:"data"`
	Metadata Metadata  `json:"metadata"`
}

type FngData struct {
	Value               string `json:"value"`
	ValueClassification string `json:"value_classification"`
	Timestamp           string `json:"timestamp"`
	TimeUntilUpdate     string `json:"time_until_update"`
}

type Metadata struct {
	Error *string `json:"error"`
}
