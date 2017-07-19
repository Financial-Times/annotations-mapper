package main

// AlphavilleSeriesService extracts and transforms the series taxonomy into an annotation
type AlphavilleSeriesService struct {
	HandledTaxonomy string
}

const alphavilleSeriesURI = "http://www.ft.com/ontology/AlphavilleSeries"

// BuildAnnotations builds a list of topic annotations from a ContentRef.
// Returns an empty array in case no topic annotations are found
func (alphavilleSeriesService AlphavilleSeriesService) buildAnnotations(contentRef ContentRef) []annotation {
	series := extractTags(alphavilleSeriesService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range series {
		annotations = append(annotations, buildAnnotation(value, alphavilleSeriesURI, classification))
	}

	return annotations
}
