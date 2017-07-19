package main

// PeopleService extracts and transforms the subject taxonomy into an annotation
type PeopleService struct {
	HandledTaxonomy string
}

const personURI = "http://www.ft.com/ontology/person/Person"

// BuildAnnotations builds a list of subject annotations from a ContentRef.
// Returns an empty array in case no subject annotations are found
func (peopleService PeopleService) buildAnnotations(contentRef ContentRef) []annotation {
	people := extractTags(peopleService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range people {
		annotations = append(annotations, buildAnnotation(value, personURI, conceptMajorMentions))
	}

	if contentRef.PrimaryTheme.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimaryTheme.ID),
			PrefLabel: contentRef.PrimaryTheme.CanonicalName,
			Predicate: about,
			Types:     []string{personURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
