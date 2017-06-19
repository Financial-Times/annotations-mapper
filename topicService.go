package main

// TopicService extracts and transforms the topic taxonomy into an annotation
type TopicService struct {
	HandledTaxonomy string
}

const topicURI = "http://www.ft.com/ontology/Topic"

// BuildAnnotations builds a list of topic annotations from a ContentRef.
// Returns an empty array in case no topic annotations are found
func (topicService TopicService) buildAnnotations(contentRef ContentRef) []annotation {
	topics := extractTags(topicService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range topics {
		annotations = append(annotations, buildAnnotation(value, topicURI, conceptMajorMentions))
	}

	if contentRef.PrimaryTheme.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimaryTheme.ID),
			PrefLabel: contentRef.PrimaryTheme.CanonicalName,
			Predicate: about,
			Types:     []string{topicURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
