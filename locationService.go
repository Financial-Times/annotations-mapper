package main

// LocationService extracts and transforms the location taxonomy into an annotation
type LocationService struct {
	HandledTaxonomy string
}

const locationURI = "http://www.ft.com/ontology/Location"

// BuildAnnotations builds a list of location annotations from a ContentRef.
// Returns an empty array in case no location annotations are found
func (locationService LocationService) buildAnnotations(contentRef ContentRef) []annotation {
	locations := extractTags(locationService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range locations {
		annotations = append(annotations, buildAnnotation(value, locationURI, conceptMajorMentions))
	}

	if contentRef.PrimaryTheme.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimaryTheme.ID),
			PrefLabel: contentRef.PrimaryTheme.CanonicalName,
			Predicate: about,
			Types:     []string{locationURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
