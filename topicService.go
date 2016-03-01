package main

// TopicService extracts and transforms the topic taxonomy into a suggestion
type TopicService struct {
	HandledTaxonomy string
}

const topicURI = "http://www.ft.com/ontology/thing/Topic"

// BuildSuggestions builds a list of topic suggestions from a ContentRef.
// Returns an empty array in case no topic annotations are found
func (topicService TopicService) buildSuggestions(contentRef ContentRef) []suggestion {
	topics := extractTags(topicService.HandledTaxonomy, contentRef)
	suggestions := []suggestion{}

	for _, value := range topics {
		suggestions = append(suggestions, buildSuggestion(value, topicURI, conceptMentions))
		suggestions = append(suggestions, buildSuggestion(value, topicURI, conceptAbout))
	}

	return suggestions
}
