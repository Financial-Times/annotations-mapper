package main

// GenreService extracts and transforms the genre taxonomy into an annotation
type GenreService struct {
	HandledTaxonomy string
}

const genreURI = "http://www.ft.com/ontology/Genre"

// BuildAnnotations builds a list of genre annotations from a ContentRef.
// Returns an empty array in case no genre annotations are found
func (genreService GenreService) buildAnnotations(contentRef ContentRef) []annotation {
	genres := extractTags(genreService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range genres {
		annotations = append(annotations, buildAnnotation(value, genreURI, classification))
	}

	return annotations
}
