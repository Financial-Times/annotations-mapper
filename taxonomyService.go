package main

import (
	uuidutils "github.com/Financial-Times/uuid-utils-go"
	"strings"
)

// TaxonomyService defines the operations used to process taxonomies
type TaxonomyService interface {
	buildAnnotations(ContentRef) []annotation
}

const (
	conceptMentions       = "mentions"
	conceptMajorMentions  = "majorMentions"
	classification        = "isClassifiedBy"
	primaryClassification = "isPrimarilyClassifiedBy"
	about                 = "about"
	hasAuthor             = "hasAuthor"

	relevanceURI  = "http://api.ft.com/scoringsystem/FT-RELEVANCE-SYSTEM"
	confidenceURI = "http://api.ft.com/scoringsystem/FT-CONFIDENCE-SYSTEM"
)

func transformScore(score int) float32 {
	return float32(score) / float32(100.0)
}

func generateID(cmrTermID string) string {
	return "http://api.ft.com/things/" + uuidutils.NewV3UUID(cmrTermID).String()
}

func extractTags(wantedTagName string, contentRef ContentRef) []tag {
	var wantedTags []tag
	for _, tag := range contentRef.TagHolder.Tags {
		if strings.EqualFold(tag.Term.Taxonomy, wantedTagName) {
			wantedTags = append(wantedTags, tag)
		}
	}
	return wantedTags
}

func buildAnnotation(tag tag, thingType string, predicate string) annotation {
	relevance := score{
		ScoringSystem: relevanceURI,
		Value:         transformScore(tag.TagScore.Relevance),
	}
	confidence := score{
		ScoringSystem: confidenceURI,
		Value:         transformScore(tag.TagScore.Confidence),
	}

	provenances := []provenance{
		provenance{
			Scores: []score{relevance, confidence},
		},
	}
	thing := thing{
		ID:        generateID(tag.Term.ID),
		PrefLabel: tag.Term.CanonicalName,
		Predicate: predicate,
		Types:     []string{thingType},
	}

	return annotation{Thing: thing, Provenance: provenances}
}
