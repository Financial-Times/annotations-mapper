[![Circle CI](https://circleci.com/gh/Financial-Times/annotations-mapper.svg?style=shield)](https://circleci.com/gh/Financial-Times/annotations-mapper)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/annotations-mapper)](https://goreportcard.com/report/github.com/Financial-Times/annotations-mapper) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/annotations-mapper/badge.svg)](https://coveralls.io/github/Financial-Times/annotations-mapper)

# Annotations Mapper

Processes metadata about content that comes from QMI system - aka V1 metadata.  

* Reads V1 metadata for an article from the  kafka source topic _NativeCmsMetadataPublicationEvents_
* Filters and transforms it to UP standard json representation
* Puts the result onto the kafka destination topic _ConceptAnnotations_

## Prerequisites

* Go Dep
* Kafka + Zookeeper - either locally or installed in an external cluster

## Installation

```
go get -u github.com/Financial-Times/annotations-mapper
cd $GOPATH/src/github.com/Financial-Times/annotations-mapper
dep ensure -vendor-only
go build
```

## Startup parameters
You can find the necessary startup parameters by running:
```bash
./annotations-mapper --help
```

## Run locally

Set the required environment variables:
```
export|set ZOOKEEPR_ADDRESS=http://kafkahost:9092
export|set CONSUMER_GROUP=FooGroup
export|set CONSUMER_TOPIC=FooBarEvents
export|set BROKER_ADDRESS=http://kafkahost:9092
export|set PRODUCER_TOPIC=DestTopic
```

And run the binary.
```
./annotations-mapper[.exe]
```

## Build in Docker
````
git config remote.origin.url https://github.com/Financial-Times/annotations-mapper.git
docker build -t coco/annotations-mapper:$DOCKER_APP_VERSION .
git config remote.origin.url git@github.com:Financial-Times/annotations-mapper.git
````

## Run in Docker
````
docker run --name annotations-mapper -p 8080 \
	--env "ZOOKEEPER_ADDRESS=http://zookeper:9092" \
	--env "CONSUMER_GROUP=annotations-mapper" \
	--env "CONSUMER_TOPIC=NativeCmsMetadataPublicationEvents" \
	--env "BROKER_ADDRESS=http://kafka:9092" \
	--env "PRODUCER_TOPIC=ConceptAnnotations" \
	coco/annotations-mapper:$DOCKER_APP_VERSION
````

## Admin Endpoints
|Endpoint     | Explanation |
|---|---|
| /__health      | checks that annotations-mapper can communicate to kafka|
|/__ping         | _response status_: **200**  _body_:**"pong"** |
|/ping           | the same as above for compatibility with Dropwizard java apps |
|/__gtg          | _response status_: **200** when "good to go" or **503** when not "good to go"|
|/__build-info   | consisting of _**version** (release tag), git **repository** url, **revision** (git commit-id), deployment **datetime**, **builder** (go or java or ...)_
|/build-info     | the same as above for compatibility with Dropwizard java apps |


## Example Message-In
````
FTMSG/1.0  
Content-Type: application/json  
Message-Id: 266c7604-b582-47a3-9b7e-c8aad93f1ec9  
Message-Timestamp: 2016-12-29T14:54:10.160Z  
Message-Type: cms-content-published  
Origin-System-Id: http://cmdb.ft.com/systems/methode-web-pub
X-Request-Id: tid_9rvfuynl4b  
{"value":"<base64 encoded message body>"}  
````

**Decoded Message-In body**
````
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>  
<ns5:contentRef ns5:created="2016-12-29T14:54:10.000Z" ns5:id="3505101"
	xmlns:ns14="http://metadata.internal.ft.com/metadata/xsd/metadata_concept_v1.0.xsd"
	xmlns:ns9="http://metadata.internal.ft.com/metadata/xsd/metadata_taxonomy_v1.0.xsd"
	xmlns:ns5="http://metadata.internal.ft.com/metadata/xsd/metadata_content_reference_v1.0.xsd"
	xmlns:ns12="http://metadata.internal.ft.com/metadata/xsd/metadata_notification_v1.0.xsd"
	xmlns:ns13="http://metadata.internal.ft.com/metadata/xsd/metadata_search_v1.0.xsd"
	xmlns:ns6="http://metadata.internal.ft.com/metadata/xsd/metadata_tag_v1.0.xsd"
	xmlns:ns7="http://metadata.internal.ft.com/metadata/xsd/metadata_binding_v1.0.xsd"
	xmlns:ns10="http://metadata.internal.ft.com/metadata/xsd/metadata_suggestion_v1.0.xsd"
	xmlns:ns8="http://metadata.internal.ft.com/metadata/xsd/metadata_property_v1.0.xsd"
	xmlns:ns11="http://metadata.internal.ft.com/metadata/xsd/metadata_count_response_v1.0.xsd"
	xmlns:ns2="http://metadata.internal.ft.com/metadata/xsd/metadata_party_v1.0.xsd"
	xmlns:ns1="http://metadata.internal.ft.com/metadata/xsd/metadata_base_v1.0.xsd"
	xmlns:ns4="http://metadata.internal.ft.com/metadata/xsd/metadata_term_v1.0.xsd"
	xmlns:ns3="http://metadata.internal.ft.com/metadata/xsd/metadata_lifecycle_v1.0.xsd">  
	<ns5:primarySection ns4:status="ACTIVE" ns4:externalTermId="116" ns4:taxonomy="Sections" ns1:id="MTE2-U2VjdGlvbnM=">  
	<ns4:canonicalName>  
	Comment</ns4:canonicalName>  
</ns5:primarySection>  
<ns5:primaryTheme ns4:status="ACTIVE" ns4:externalTermId="a8e4a619-3c38-41fd-9e20-8ac64ed06447" ns4:taxonomy="Topics" ns1:id="YThlNGE2MTktM2MzOC00MWZkLTllMjAtOGFjNjRlZDA2NDQ3-VG9waWNz">  
	<ns4:canonicalName>  
	Global politics</ns4:canonicalName>  
</ns5:primaryTheme>  
<ns5:tags>  
	<ns6:tag>  
	<ns6:meta ns1:provenance="USER"/>  
<ns6:term ns4:status="ACTIVE" ns4:externalTermId="a8e4a619-3c38-41fd-9e20-8ac64ed06447" ns4:taxonomy="Topics" ns1:id="YThlNGE2MTktM2MzOC00MWZkLTllMjAtOGFjNjRlZDA2NDQ3-VG9waWNz">  
	<ns4:canonicalName>  
	Global politics</ns4:canonicalName>  
</ns6:term>  
<ns6:score ns6:relevance="100" ns6:confidence="100"/>  
</ns6:tag>  
<ns6:tag>  
	<ns6:meta ns1:provenance="USER"/>  
<ns6:term ns4:status="ACTIVE" ns4:externalTermId="8" ns4:taxonomy="Genres" ns1:id="OA==-R2VucmVz">  
	<ns4:canonicalName>  
	Comment</ns4:canonicalName>  
</ns6:term>  
<ns6:score ns6:relevance="100" ns6:confidence="100"/>  
</ns6:tag>  
<ns6:tag>  
	<ns6:meta ns1:provenance="USER"/>  
<ns6:term ns4:status="ACTIVE" ns4:externalTermId="116" ns4:taxonomy="Sections" ns1:id="MTE2-U2VjdGlvbnM=">  
	<ns4:canonicalName>  
	Comment</ns4:canonicalName>  
</ns6:term>  
<ns6:score ns6:relevance="100" ns6:confidence="100"/>  
</ns6:tag>  
<ns6:tag>  
	<ns6:meta ns1:provenance="PREPROCESSOR"/>  
<ns6:term ns4:status="ACTIVE" ns4:externalTermId="f30ca667-0056-4e98-b41e-f99196e324ef" ns4:taxonomy="MediaTypes" ns1:id="ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==">  
	<ns4:canonicalName>  
	Text</ns4:canonicalName>  
</ns6:term>  
<ns6:score ns6:relevance="100" ns6:confidence="100"/>  
</ns6:tag>  
</ns5:tags>  
<ns5:externalReferences>  
	<ns7:reference ns1:cmrId="1227570" ns1:externalId="980913e6-cdd6-11e6-864f-20dcb35cede2" ns1:externalSource="METHODE"/>  
</ns5:externalReferences>  
</ns5:contentRef>  
````
