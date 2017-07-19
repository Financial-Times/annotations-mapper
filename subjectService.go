package main

// SubjectService extracts and transforms the subject taxonomy into an annotation
type SubjectService struct {
	HandledTaxonomy string
}

const subjectURI = "http://www.ft.com/ontology/Subject"

// BuildAnnotations builds a list of subject annotations from a ContentRef.
// Returns an empty array in case no subject annotations are found
func (subjectService SubjectService) buildAnnotations(contentRef ContentRef) []annotation {
	subjects := extractTags(subjectService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range subjects {
		annotations = append(annotations, buildAnnotation(value, subjectURI, classification))
	}

	return annotations
}
