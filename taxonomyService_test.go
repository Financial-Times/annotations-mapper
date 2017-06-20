package main

import (
	"fmt"
	"testing"

	uuidutils "github.com/Financial-Times/uuid-utils-go"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestSubjectServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := SubjectService{"subjects"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 subject tag",
			buildContentRefWithSubjects(1),
			buildConceptAnnotationsWithSubjects(1),
		},
		{"Build concept annotation from a contentRef with no subject tags",
			buildContentRefWithSubjects(0),
			buildConceptAnnotationsWithSubjects(0),
		},
		{"Build concept annotation from a contentRef with multiple subject tags",
			buildContentRefWithSubjects(2),
			buildConceptAnnotationsWithSubjects(2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestSectionServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := SectionService{"sections"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 section tag",
			buildContentRefWithSections(1),
			buildConceptAnnotationsWithSections(1),
		},
		{"Build concept annotation from a contentRef with no section tags",
			buildContentRefWithSections(0),
			buildConceptAnnotationsWithSections(0),
		},
		{"Build concept annotation from a contentRef with multiple section tags",
			buildContentRefWithSections(2),
			buildConceptAnnotationsWithSections(2),
		},
		{"Build concept annotation from a contentRef with a primary section",
			buildContentRefWithPrimarySection("sections", 2),
			buildConceptAnnotationsWithPrimarySection("sections", 2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestTopicServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := TopicService{"topics"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 topic tag",
			buildContentRefWithTopics(1),
			buildConceptAnnotationsWithTopics(1),
		},
		{"Build concept annotation from a contentRef with no topic tags",
			buildContentRefWithTopics(0),
			buildConceptAnnotationsWithTopics(0),
		},
		{"Build concept annotation from a contentRef with multiple topic tags",
			buildContentRefWithTopics(2),
			buildConceptAnnotationsWithTopics(2),
		},
		{"Build concept annotation from a contentRef with 1 topic tag and a primary theme",
			buildContentRefWithTopicsWithPrimaryTheme(1),
			buildConceptAnnotationsWithTopicsWithPrimaryTheme(1),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestLocationServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := LocationService{"gl"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 location tag",
			buildContentRefWithLocations(1),
			buildConceptAnnotationsWithLocations(1),
		},
		{"Build concept annotation from a contentRef with no location tags",
			buildContentRefWithLocations(0),
			buildConceptAnnotationsWithLocations(0),
		},
		{"Build concept annotation from a contentRef with multiple location tags",
			buildContentRefWithLocations(2),
			buildConceptAnnotationsWithLocations(2),
		},
		{"Build concept annotation from a contentRef with 1 location tag and primamry location",
			buildContentRefWithLocationsWithPrimaryTheme(1),
			buildConceptAnnotationsWithLocationsWithPrimaryTheme(1),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestGenreServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := GenreService{"genres"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 genre tag",
			buildContentRefWithGenres(1),
			buildConceptAnnotationsWithGenres(1),
		},
		{"Build concept annotation from a contentRef with no genre tags",
			buildContentRefWithGenres(0),
			buildConceptAnnotationsWithGenres(0),
		},
		{"Build concept annotation from a contentRef with multiple genre tags",
			buildContentRefWithGenres(2),
			buildConceptAnnotationsWithGenres(2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestSpecialReportServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := SpecialReportService{"specialReports"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 specialReports tag",
			buildContentRefWithSpecialReports(1),
			buildConceptAnnotationsWithSpecialReports(1),
		},
		{"Build concept annotation from a contentRef with no specialReports tags",
			buildContentRefWithSpecialReports(0),
			buildConceptAnnotationsWithSpecialReports(0),
		},
		{"Build concept annotation from a contentRef with multiple specialReports tags",
			buildContentRefWithSpecialReports(2),
			buildConceptAnnotationsWithSpecialReports(2),
		},
		{"Build concept annotation from a contentRef with specialReports as a primary section",
			buildContentRefWithPrimarySection("specialReports", 2),
			buildConceptAnnotationsWithPrimarySection("specialReports", 2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func TestAlphavilleSeriesServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := AlphavilleSeriesService{"alphavilleSeries"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 alphavilleSeries tag",
			buildContentRefWithAlphavilleSeries(1),
			buildConceptAnnotationsWithAlphavilleSeries(1),
		},
		{"Build concept annotation from a contentRef with no alphavilleSeries tags",
			buildContentRefWithAlphavilleSeries(0),
			buildConceptAnnotationsWithAlphavilleSeries(0),
		},
		{"Build concept annotation from a contentRef with multiple alphavilleSeries tags",
			buildContentRefWithAlphavilleSeries(2),
			buildConceptAnnotationsWithAlphavilleSeries(2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect.", test.name))
	}
}

func TestOrganisationsServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := OrganisationService{"on"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 organisation tag",
			buildContentRefWithOrganisations(1),
			buildConceptAnnotationsWithOrganisations(1),
		},
		{"Build concept annotation from a contentRef with no organisation tags",
			buildContentRefWithOrganisations(0),
			buildConceptAnnotationsWithOrganisations(0),
		},
		{"Build concept annotation from a contentRef with multiple organisation tags",
			buildContentRefWithOrganisations(2),
			buildConceptAnnotationsWithOrganisations(2),
		},
		{"Build concept annotation from a contentRef with 1 organisation tag and a primary theme",
			buildContentRefWithOrganisationWithPrimaryTheme(1),
			buildConceptAnnotationsWithOrganisationsWithPrimaryTheme(1),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect: ACTUAL: %s  TEST: %s ", test.name, actualConceptAnnotations, test.annotations))
	}
}

func TestPeopleServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := PeopleService{"PN"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 Person tag",
			buildContentRefWithPeople(1),
			buildConceptAnnotationsWithPeople(1),
		},
		{"Build concept annotation from a contentRef with no Person tags",
			buildContentRefWithPeople(0),
			buildConceptAnnotationsWithPeople(0),
		},
		{"Build concept annotation from a contentRef with multiple Person tags",
			buildContentRefWithPeople(2),
			buildConceptAnnotationsWithPeople(2),
		},
		{"Build concept annotation from a contentRef with a primary theme",
			buildContentRefWithPeopleWithPrimaryTheme(2),
			buildConceptAnnotationsWithPeopleWithPrimaryTheme(2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations,
			actualConceptAnnotations,
			fmt.Sprintf("%s: Actual concept annotations incorrect: ACTUAL: %s  TEST: %s ",
				test.name,
				actualConceptAnnotations,
				test.annotations))
	}
}

func TestAuthorServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := AuthorService{"Authors"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 author tag",
			buildContentRefWithAuthor(1),
			buildConceptAnnotationsWithAuthor(1),
		},
		{"Build concept annotation from a contentRef with no author tags",
			buildContentRefWithAuthor(0),
			buildConceptAnnotationsWithAuthor(0),
		},
		{"Build concept annotation from a contentRef with multiple author tags",
			buildContentRefWithAuthor(2),
			buildConceptAnnotationsWithAuthor(2),
		},
	}

	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect.", test.name))
	}
}

func TestBrandServiceBuildAnnotations(t *testing.T) {
	assert := assert.New(t)
	service := BrandService{HandledTaxonomy: "brands"}
	tests := []struct {
		name        string
		contentRef  ContentRef
		annotations []annotation
	}{
		{"Build concept annotation from a contentRef with 1 brand tag",
			buildContentRefWithBrands(1),
			buildConceptAnnotationsWithBrands(1),
		},
		{"Build concept annotation from a contentRef with no brand tags",
			buildContentRefWithBrands(0),
			buildConceptAnnotationsWithBrands(0),
		},
		{"Build concept annotation from a contentRef with multiple brand tags",
			buildContentRefWithBrands(2),
			buildConceptAnnotationsWithBrands(2),
		},
	}
	for _, test := range tests {
		actualConceptAnnotations := service.buildAnnotations(test.contentRef)
		assert.Equal(test.annotations, actualConceptAnnotations, fmt.Sprintf("%s: Actual concept annotations incorrect", test.name))
	}
}

func buildContentRefWithLocations(locationCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["locations"] = locationCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithLocationsWithPrimaryTheme(locationCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["locations"] = locationCount
	return buildContentRef(taxonomyAndCount, false, true)
}

func buildContentRefWithTopics(topicCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["topics"] = topicCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithTopicsWithPrimaryTheme(topicCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["topics"] = topicCount
	return buildContentRef(taxonomyAndCount, false, true)
}

func buildContentRefWithSubjects(subjectCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["subjects"] = subjectCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithSections(sectionCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["sections"] = sectionCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithPrimarySection(taxonomyName string, sectionCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount[taxonomyName] = sectionCount
	return buildContentRef(taxonomyAndCount, true, false)
}

func buildContentRefWithGenres(genreCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["genres"] = genreCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithBrands(count int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["brands"] = count
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithSpecialReports(specialReportCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["specialReports"] = specialReportCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithAlphavilleSeries(seriesCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["alphavilleSeries"] = seriesCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithOrganisations(organisationCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["organisations"] = organisationCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithOrganisationWithPrimaryTheme(orgCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["organisations"] = orgCount
	return buildContentRef(taxonomyAndCount, false, true)
}

func buildContentRefWithPeople(peopleCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["people"] = peopleCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRefWithPeopleWithPrimaryTheme(peopleCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["people"] = peopleCount
	return buildContentRef(taxonomyAndCount, false, true)
}

func buildContentRefWithAuthor(authorCount int) ContentRef {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["author"] = authorCount
	return buildContentRef(taxonomyAndCount, false, false)
}

func buildContentRef(taxonomyAndCount map[string]int, hasPrimarySection bool, hasPrimaryTheme bool) ContentRef {
	metadataTags := []tag{}
	var primarySection term
	var primaryTheme term
	for key, count := range taxonomyAndCount {
		if strings.EqualFold("subjects", key) {
			for i := 0; i < count; i++ {
				subjectTerm := term{CanonicalName: subjectNames[i], Taxonomy: "Subjects", ID: subjectTMEIDs[i]}
				metadataTags = append(metadataTags, tag{Term: subjectTerm, TagScore: testScore})
			}
		}
		if strings.EqualFold("sections", key) {
			for i := 0; i < count; i++ {
				sectionsTerm := term{CanonicalName: sectionNames[i], Taxonomy: "Sections", ID: sectionTMEIDs[i]}
				sectionsTag := tag{Term: sectionsTerm, TagScore: testScore}
				metadataTags = append(metadataTags, sectionsTag)
			}

			if hasPrimarySection {
				primarySection = term{CanonicalName: sectionNames[0], Taxonomy: "Sections", ID: sectionTMEIDs[0]}
			}
		}
		if strings.EqualFold("topics", key) {
			for i := 0; i < count; i++ {
				topicTerm := term{CanonicalName: topicNames[i], Taxonomy: "Topics", ID: topicTMEIDs[i]}
				topicTag := tag{Term: topicTerm, TagScore: testScore}
				metadataTags = append(metadataTags, topicTag)
			}
			if hasPrimaryTheme {
				primaryTheme = term{CanonicalName: topicNames[0], Taxonomy: "Topics", ID: topicTMEIDs[0]}
			}
		}
		if strings.EqualFold("locations", key) {
			for i := 0; i < count; i++ {
				locationTerm := term{CanonicalName: locationNames[i], Taxonomy: "GL", ID: locationTMEIDs[i]}
				locationTag := tag{Term: locationTerm, TagScore: testScore}
				metadataTags = append(metadataTags, locationTag)
			}
			if hasPrimaryTheme {
				primaryTheme = term{CanonicalName: locationNames[0], Taxonomy: "GL", ID: locationTMEIDs[0]}
			}
		}
		if strings.EqualFold("genres", key) {
			for i := 0; i < count; i++ {
				genreTerm := term{CanonicalName: genreNames[i], Taxonomy: "Genres", ID: genreTMEIDs[i]}
				metadataTags = append(metadataTags, tag{Term: genreTerm, TagScore: testScore})
			}
		}
		if strings.EqualFold("brands", key) {
			for i := 0; i < count; i++ {
				brandTerm := term{CanonicalName: brandNames[i], Taxonomy: "Brands", ID: brandTMEIDs[i]}
				metadataTags = append(metadataTags, tag{Term: brandTerm, TagScore: testScore})
			}
		}
		if strings.EqualFold("specialReports", key) {
			for i := 0; i < count; i++ {
				specialReportsTerm := term{CanonicalName: specialReportNames[i], Taxonomy: "SpecialReports", ID: specialReportTMEIDs[i]}
				sectionsTag := tag{Term: specialReportsTerm, TagScore: testScore}
				metadataTags = append(metadataTags, sectionsTag)
			}

			if hasPrimarySection {
				primarySection = term{CanonicalName: specialReportNames[0], Taxonomy: "SpecialReports", ID: specialReportTMEIDs[0]}
			}
		}
		if strings.EqualFold("alphavilleSeries", key) {
			for i := 0; i < count; i++ {
				alphavilleSeriesTerm := term{CanonicalName: alphavilleSeriesNames[i], Taxonomy: "AlphavilleSeries", ID: alphavilleSeriesTMEIDs[i]}
				alphavilleSeriesTag := tag{Term: alphavilleSeriesTerm, TagScore: testScore}
				metadataTags = append(metadataTags, alphavilleSeriesTag)
			}
		}
		if strings.EqualFold("organisations", key) {
			for i := 0; i < count; i++ {
				organisationTerm := term{CanonicalName: organisationNames[i], Taxonomy: "ON", ID: organisationTMEIDs[i]}
				organisationTag := tag{Term: organisationTerm, TagScore: testScore}
				metadataTags = append(metadataTags, organisationTag)
			}
			if hasPrimaryTheme {
				primaryTheme = term{CanonicalName: organisationNames[0], Taxonomy: "Organisations", ID: organisationTMEIDs[0]}
			}
		}

		if strings.EqualFold("people", key) {
			for i := 0; i < count; i++ {
				peopleTerm := term{CanonicalName: peopleNames[i], Taxonomy: "PN", ID: peopleTMEIDs[i]}
				peopleTag := tag{Term: peopleTerm, TagScore: testScore}
				metadataTags = append(metadataTags, peopleTag)
			}
			if hasPrimaryTheme {
				primaryTheme = term{CanonicalName: peopleNames[0], Taxonomy: "People", ID: peopleTMEIDs[0]}
			}
		}

		if strings.EqualFold("author", key) {
			for i := 0; i < count; i++ {
				authourTerm := term{CanonicalName: authorNames[i], Taxonomy: "Authors", ID: authorTMEIDs[i]}
				authorTag := tag{Term: authourTerm, TagScore: testScore}
				metadataTags = append(metadataTags, authorTag)
			}
		}
	}
	tagHolder := tags{Tags: metadataTags}

	return ContentRef{TagHolder: tagHolder, PrimarySection: primarySection, PrimaryTheme: primaryTheme}
}

func buildConceptAnnotationsWithLocations(locationCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["locations"] = locationCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithLocationsWithPrimaryTheme(locationCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["locations"] = locationCount
	return buildConceptAnnotations(taxonomyAndCount, false, true)
}

func buildConceptAnnotationsWithTopics(topicCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["topics"] = topicCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithTopicsWithPrimaryTheme(topicsCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["topics"] = topicsCount
	return buildConceptAnnotations(taxonomyAndCount, false, true)
}

func buildConceptAnnotationsWithSections(sectionCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["sections"] = sectionCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithPrimarySection(taxonomyName string, sectionCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount[taxonomyName] = sectionCount
	return buildConceptAnnotations(taxonomyAndCount, true, false)
}

func buildConceptAnnotationsWithSubjects(subjectCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["subjects"] = subjectCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithGenres(genreCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["genres"] = genreCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithBrands(count int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["brands"] = count
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithSpecialReports(reportCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["specialReports"] = reportCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithAlphavilleSeries(seriesCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["alphavilleSeries"] = seriesCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithOrganisations(orgsCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["organisations"] = orgsCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithOrganisationsWithPrimaryTheme(orgsCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["organisations"] = orgsCount
	return buildConceptAnnotations(taxonomyAndCount, false, true)
}

func buildConceptAnnotationsWithPeople(peopleCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["people"] = peopleCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotationsWithPeopleWithPrimaryTheme(peopleCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["people"] = peopleCount
	return buildConceptAnnotations(taxonomyAndCount, false, true)
}

func buildConceptAnnotationsWithAuthor(authorCount int) []annotation {
	taxonomyAndCount := make(map[string]int)
	taxonomyAndCount["author"] = authorCount
	return buildConceptAnnotations(taxonomyAndCount, false, false)
}

func buildConceptAnnotations(taxonomyAndCount map[string]int, hasPrimarySection bool, hasPrimaryTheme bool) []annotation {
	annotations := []annotation{}

	relevance := score{ScoringSystem: relevanceURI, Value: 0.65}
	confidence := score{ScoringSystem: confidenceURI, Value: 0.93}
	metadataProvenance := provenance{Scores: []score{relevance, confidence}}
	for key, count := range taxonomyAndCount {
		if strings.EqualFold("subjects", key) {
			for i := 0; i < count; i++ {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(subjectTMEIDs[i]).String(),
					PrefLabel: subjectNames[i],
					Predicate: classification,
					Types:     []string{subjectURI},
				}
				subjectAnnotation := annotation{Thing: thing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, subjectAnnotation)
			}
		}
		if strings.EqualFold("sections", key) {
			for i := 0; i < count; i++ {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(sectionTMEIDs[i]).String(),
					PrefLabel: sectionNames[i],
					Predicate: classification,
					Types:     []string{sectionURI},
				}
				sectionAnnotation := annotation{Thing: thing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, sectionAnnotation)
			}

			if count > 0 && hasPrimarySection {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(sectionTMEIDs[0]).String(),
					PrefLabel: sectionNames[0],
					Predicate: primaryClassification,
					Types:     []string{sectionURI},
				}
				sectionAnnotation := annotation{Thing: thing}
				annotations = append(annotations, sectionAnnotation)
			}
		}
		if strings.EqualFold("topics", key) {
			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(topicTMEIDs[i]).String(),
					PrefLabel: topicNames[i],
					Predicate: conceptMajorMentions,
					Types:     []string{topicURI},
				}
				topicAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}

				annotations = append(annotations, topicAnnotation)
			}
			if count > 0 && hasPrimaryTheme {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(topicTMEIDs[0]).String(),
					PrefLabel: topicNames[0],
					Predicate: about,
					Types:     []string{topicURI},
				}
				topicAnnotation := annotation{Thing: thing}
				annotations = append(annotations, topicAnnotation)
			}
		}
		if strings.EqualFold("locations", key) {
			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(locationTMEIDs[i]).String(),
					PrefLabel: locationNames[i],
					Predicate: conceptMajorMentions,
					Types:     []string{locationURI},
				}
				locationAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}

				annotations = append(annotations, locationAnnotation)
			}
			if count > 0 && hasPrimaryTheme {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(locationTMEIDs[0]).String(),
					PrefLabel: locationNames[0],
					Predicate: about,
					Types:     []string{locationURI},
				}
				locationAnnotation := annotation{Thing: thing}
				annotations = append(annotations, locationAnnotation)
			}
		}
		if strings.EqualFold("genres", key) {
			for i := 0; i < count; i++ {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(genreTMEIDs[i]).String(),
					PrefLabel: genreNames[i],
					Predicate: classification,
					Types:     []string{genreURI},
				}
				genreAnnotation := annotation{Thing: thing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, genreAnnotation)
			}
		}
		if strings.EqualFold("brands", key) {
			for i := 0; i < count; i++ {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(brandTMEIDs[i]).String(),
					PrefLabel: brandNames[i],
					Predicate: classification,
					Types:     []string{brandURI},
				}
				brandAnnotation := annotation{Thing: thing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, brandAnnotation)
			}
		}
		if strings.EqualFold("specialReports", key) {
			for i := 0; i < count; i++ {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(specialReportTMEIDs[i]).String(),
					PrefLabel: specialReportNames[i],
					Predicate: classification,
					Types:     []string{specialReportURI},
				}
				specialReportAnnotation := annotation{Thing: thing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, specialReportAnnotation)
			}

			if count > 0 && hasPrimarySection {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(specialReportTMEIDs[0]).String(),
					PrefLabel: specialReportNames[0],
					Predicate: primaryClassification,
					Types:     []string{specialReportURI},
				}
				specialReportAnnotation := annotation{Thing: thing}
				annotations = append(annotations, specialReportAnnotation)
			}
		}
		if strings.EqualFold("alphavilleSeries", key) {
			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(alphavilleSeriesTMEIDs[i]).String(),
					PrefLabel: alphavilleSeriesNames[i],
					Predicate: classification,
					Types:     []string{alphavilleSeriesURI},
				}
				alphavilleSeriesAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}

				annotations = append(annotations, alphavilleSeriesAnnotation)
			}
		}
		if strings.EqualFold("organisations", key) {

			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(organisationTMEIDs[i]).String(),
					PrefLabel: organisationNames[i],
					Predicate: conceptMajorMentions,
					Types:     []string{organisationURI},
				}
				organisationAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}

				annotations = append(annotations, organisationAnnotation)

			}
			if count > 0 && hasPrimaryTheme {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(organisationTMEIDs[0]).String(),
					PrefLabel: organisationNames[0],
					Predicate: about,
					Types:     []string{organisationURI},
				}
				organisationAnnotation := annotation{Thing: thing}
				annotations = append(annotations, organisationAnnotation)
			}

		}
		if strings.EqualFold("people", key) {

			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(peopleTMEIDs[i]).String(),
					PrefLabel: peopleNames[i],
					Predicate: conceptMajorMentions,
					Types:     []string{personURI},
				}
				peopleAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}
				annotations = append(annotations, peopleAnnotation)

			}

			if count > 0 && hasPrimaryTheme {
				thing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(peopleTMEIDs[0]).String(),
					PrefLabel: peopleNames[0],
					Predicate: about,
					Types:     []string{personURI},
				}
				peopleAnnotation := annotation{Thing: thing}
				annotations = append(annotations, peopleAnnotation)
			}
		}
		if strings.EqualFold("author", key) {
			for i := 0; i < count; i++ {
				oneThing := thing{
					ID:        "http://api.ft.com/things/" + uuidutils.NewV3UUID(authorTMEIDs[i]).String(),
					PrefLabel: authorNames[i],
					Predicate: hasAuthor,
					Types:     []string{authorURI},
				}
				authorAnnotation := annotation{Thing: oneThing, Provenance: []provenance{metadataProvenance}}

				annotations = append(annotations, authorAnnotation)
			}
		}
	}

	return annotations
}

var testScore = tagScore{Confidence: 93, Relevance: 65}
var subjectNames = [...]string{"Mining Industry", "Oil Extraction Subsidies"}
var subjectTMEIDs = [...]string{"Mjk=-U2VjdGlvbnM=", "Nw==-R2VucmVz"}
var sectionNames = [...]string{"Companies", "Emerging Markets"}
var sectionTMEIDs = [...]string{"Nw==-R2Bucm3z", "Nw==-U2VjdGlvbnM="}
var topicNames = [...]string{"Big Data", "BP trial"}
var topicTMEIDs = [...]string{"M2YyN2I0NGEtZGZjMi00MDVjLTlkNjAtNGRlNTNhM2EwYjlm-VG9waWNz", "ZWE3YzNhNmQtNGU4MS00MzE0LWIxZWMtYWQxY2M4Y2ZjZDFk-VG9waWNz"}
var locationNames = [...]string{"New York", "Rio"}
var locationTMEIDs = [...]string{"TmV3IFlvcms=-R0w=", "Umlv-R0w="}
var genreNames = [...]string{"News", "Letter"}
var genreTMEIDs = [...]string{"TmV3cw==-R2VucmVz", "TGV0dGVy-R2VucmVz"}
var brandNames = [...]string{"FT", "Martin Wolf"}
var brandTMEIDs = [...]string{"RlQK-QnJhbmRzCg==", "TWFydGluIFdvbGY=-QnJhbmRzCg=="}
var specialReportNames = [...]string{"Business", "Investment"}
var specialReportTMEIDs = [...]string{"U3BlY2lhbFJlcG9ydHM=-R2Bucm3z", "U3BlY2lhbFJlcG9ydHM=-U2VjdGlvbnM="}
var alphavilleSeriesNames = [...]string{"AV Series 1", "AV Series 2"}
var alphavilleSeriesTMEIDs = [...]string{"series1-AV", "series2-AV"}
var organisationNames = [...]string{"Organisation 1", "Organisation 2"}
var organisationTMEIDs = [...]string{"Organisation-1-TME", "Organisation-2-TME"}
var peopleNames = [...]string{"Person 1", "Person 2"}
var peopleTMEIDs = [...]string{"Person-1-TME", "Person-2-TME"}
var authorNames = [...]string{"Author 1", "Author 2"}
var authorTMEIDs = [...]string{"Author-1-TME", "Author-2-TME"}
