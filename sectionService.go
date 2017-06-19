package main

// SectionService extracts and transforms the section taxonomy into an annotation
type SectionService struct {
	HandledTaxonomy string
}

const sectionURI = "http://www.ft.com/ontology/Section"

// BuildAnnotations builds a list of section annotations from a ContentRef.
// Returns an empty array in case no section annotations are found
func (sectionService SectionService) buildAnnotations(contentRef ContentRef) []annotation {
	sections := extractTags(sectionService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range sections {
		annotations = append(annotations, buildAnnotation(value, sectionURI, classification))
	}

	if contentRef.PrimarySection.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimarySection.ID),
			PrefLabel: contentRef.PrimarySection.CanonicalName,
			Predicate: primaryClassification,
			Types:     []string{sectionURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
