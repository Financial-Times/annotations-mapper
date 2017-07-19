package main

const brandURI = "http://www.ft.com/ontology/Brand"

type BrandService struct {
	HandledTaxonomy string
}

func (brandService BrandService) buildAnnotations(contentRef ContentRef) []annotation {
	authors := extractTags(brandService.HandledTaxonomy, contentRef)
	annotations := []annotation{}
	for _, value := range authors {
		annotations = append(annotations, buildAnnotation(value, brandURI, classification))
	}
	return annotations
}
