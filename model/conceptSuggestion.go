package model

// ConceptSuggestion models the suggestions as it will be written on the queue
type ConceptSuggestion struct {
	UUID        string       `json:"uuid"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Thing      Thing        `json:"thing"`
	Provenance []Provenance `json:"provenances"`
}

type Thing struct {
	ID        string   `json:"id"`
	PrefLabel string   `json:"prefLabel"`
	Types     []string `json:"types"`
}

type Provenance struct {
	Scores []Score `json:"scores"`
}

type Score struct {
	ScoringSystem string  `json:"scoringSystem"`
	Value         float32 `json:"value"`
}
