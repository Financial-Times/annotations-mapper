package main

// AuthorService extracts and transforms the author taxonomy into an annotation
type AuthorService struct {
	HandledTaxonomy string
}

const authorURI = "http://www.ft.com/ontology/person/Person"

// BuildAnnotations builds a list of author annotations from a ContentRef.
// Returns an empty array in case no author annotations are found
func (authorService AuthorService) buildAnnotations(contentRef ContentRef) []annotation {
	authors := extractTags(authorService.HandledTaxonomy, contentRef)
	annotations := []annotation{}

	for _, value := range authors {
		annotations = append(annotations, buildAnnotation(value, authorURI, hasAuthor))
	}

	return annotations
}
