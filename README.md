[![Circle CI](https://circleci.com/gh/Financial-Times/annotations-mapper.svg?style=shield)](https://circleci.com/gh/Financial-Times/annotations-mapper)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/annotations-mapper)](https://goreportcard.com/report/github.com/Financial-Times/annotations-mapper) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/annotations-mapper/badge.svg)](https://coveralls.io/github/Financial-Times/annotations-mapper)

# annotations-mapper
Processes metadata about content that comes from QMI system - aka V1 annotations.  

* Reads V1 metadata for an article from the  kafka source topic _NativeCmsMetadataPublicationEvents_
* Filters and transforms it to UP standard json representation
* Puts the result onto the kafka destination topic _V1ConceptAnnotations_

annotations-mapper service communicates with kafka via http-rest-proxy. It polls kafka-rest-proxy for messages and POSTs transformed messages to kafka-rest-proxy.  
This service is deployed in the Delivery clusters.
## Installation

```
go get -u github.com/kardianos/govendor
go get -u github.com/Financial-Times/v1-suggestor
cd $GOPATH/src/github.com/Financial-Times/v1-suggestor
govendor sync
go build .
```

## Startup parameters

| **Parameter** | **Value in prod** | **Explained** |
|---|---|---|
| **SRC_ADDR** |_http://localhost:8080_ | Url of the _http-rest-proxy_ host to connect to in order to **receive** messages from kafka. |
| **SRC_GROUP** | _v1Suggestor_ | The consumer group for receiving messages from kafka. |
| **SRC_TOPIC** | _NativeCmsMetadataPublicationEvents_ | kafka topic to consume messages from. |
| **SRC_QUEUE** | _kafka_ |  Used by _Vulcan_ to route http requests based on _Host_ header. In docker cluster all hosts are at _http://localhost:8080_. This http header is supplied to distinguish one service from another.  Host header _kafka_ points to _http-rest-proxy_. |
| **SRC_CONCURRENT_PROCESSING** | _false_ | Should the consumer process messages concurrently or sequentially. |
| **DEST_ADDRESS** | _http://localhost:8080_| Url of the _http-rest-proxy_ host to connect to in order to **send** messages to kafka. In prod env this is typically the same address as the SRC_ADDR. |
| **DEST_TOPIC** | _ConceptSuggestions_ | kafka topic to **send** messages to.  |
| **DEST_QUEUE** | _kafka_ |  Used by _Vulcan_ to route http requests based on _Host_ header. In prod docker cluster it is the same as SRC_QUEUE. |


## Prerequisites
In order to run annotations-mapper you would need at least kafka/zookeeper and kafka-rest-proxy to be accessible somewhere
and you would need to provide the host and the port to connect to them as startup parameters.

## Run locally
````
   export|set SRC_ADDR=http://kafkahost:8080
   export|set SRC_GROUP=FooGroup
   export|set SRC_TOPIC=FooBarEvents
   export|set SRC_QUEUE=kafka
   export|set SRC_CONCURRENT_PROCESSING=true
   export|set DEST_ADDRESS=http://kafkahost:8080
   export|set DEST_TOPIC=DestTopic
   export|set DEST_QUEUE=kafka
   export|set ENVIRONMENT=coco-semantic
   export|set DOCKER_APP_VERSION=latest
````

````
./annotations-mapper[.exe]
````

## Build in Docker
````
git config remote.origin.url https://github.com/Financial-Times/annotations-mapper.git
docker build -t coco/annotations-mapper:$DOCKER_APP_VERSION .
git config remote.origin.url git@github.com:Financial-Times/annotations-mapper.git
````

# #Run in Docker
````

docker run --name annotations-mapper -p 8080 \
--env "SRC_ADDR=http://kafka:8080" \
	--env "SRC_GROUP=annotations-mapper" \
	--env "SRC_TOPIC=NativeCmsMetadataPublicationEvents" \
	--env "SRC_QUEUE=kafka" \
	--env "SRC_CONCURRENT_PROCESSING=false" \
	--env "DEST_ADDRESS=http://kafka:8080" \
	--env "DEST_TOPIC=ConceptSuggestions" \
	--env "DEST_QUEUE=kafka" \
	--env "ENVIRONMENT=coco-$ENVIRONMENT_TAG" \
	coco/annotations-mapper:$DOCKER_APP_VERSION
````
 
## Admin Endpoints

|===Endpoint ===    | Explained |
|---|---|
| /__health      | checks that annotations-mapper can communicate to kafka via http-rest-proxy|
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
Origin-System-Id: http://cmdb.ft.com/systems/binding-service  
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
