package main

// OrganisationService extracts and transforms the subject taxonomy into an annotation
type OrganisationService struct {
	HandledTaxonomy string
}

const organisationURI = "http://www.ft.com/ontology/organisation/Organisation"

// BuildAnnotations builds a list of subject annotations from a ContentRef.
// Returns an empty array in case no subject annotations are found
func (organisationService OrganisationService) buildAnnotations(contentRef ContentRef) []annotation {
	subjects := extractTags(organisationService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range subjects {
		annotations = append(annotations, buildAnnotation(value, organisationURI, conceptMajorMentions))
	}

	if contentRef.PrimaryTheme.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimaryTheme.ID),
			PrefLabel: contentRef.PrimaryTheme.CanonicalName,
			Predicate: about,
			Types:     []string{organisationURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
