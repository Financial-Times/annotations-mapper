package main

// SpecialReportService extracts and transforms the special report taxonomy into an annotation
type SpecialReportService struct {
	HandledTaxonomy string
}

const specialReportURI = "http://www.ft.com/ontology/SpecialReport"

// BuildAnnotations builds a list of specialReport annotations from a ContentRef.
// Returns an empty array in case no specialReport annotations are found
func (specialReportService SpecialReportService) buildAnnotations(contentRef ContentRef) []annotation {
	specialReports := extractTags(specialReportService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range specialReports {
		annotations = append(annotations, buildAnnotation(value, specialReportURI, classification))
	}

	if contentRef.PrimarySection.CanonicalName != "" {
		thing := thing{
			ID:        generateID(contentRef.PrimarySection.ID),
			PrefLabel: contentRef.PrimarySection.CanonicalName,
			Predicate: primaryClassification,
			Types:     []string{specialReportURI},
		}
		annotations = append(annotations, annotation{Thing: thing})
	}

	return annotations
}
