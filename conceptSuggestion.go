package main

// ConceptAnnotations models the annotations as it will be written on the queue
type ConceptAnnotations struct {
	UUID        string       `json:"uuid"`
	Annotations []annotation `json:"annotations"`
}

type annotation struct {
	Thing      thing        `json:"thing"`
	Provenance []provenance `json:"provenances,omitempty"`
}

type thing struct {
	ID        string   `json:"id"`
	PrefLabel string   `json:"prefLabel"`
	Predicate string   `json:"predicate"`
	Types     []string `json:"types"`
}

type provenance struct {
	Scores []score `json:"scores"`
}

type score struct {
	ScoringSystem string  `json:"scoringSystem"`
	Value         float32 `json:"value"`
}
