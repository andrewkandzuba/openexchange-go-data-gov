# openexchange-go-data-gov

Demonstration of CQRS pattern in Golang

# Problem statement

Service to visualize Federal Commerce News feed.

## Solution

We are going to subscribe to [Commerce.gov News API](https://api.commerce.gov/api/news?api_key=DEMO_KEY) stream and poll from it
periodically in order to collect updates. 

All updates will be streamed into a queue (probably Apache Kafka) then consumed independently from the harvesting process. 
Upon consumption deduplication will be made then all new articles will be parsed out and persisted into a database.
 
The special web service will host REST API to query articles: 

- All with a pagination 

- Facet search by categories 

- Wildcard search      

## Technical requirements 

- Application should be written in Golang.

- Application to be hosted in a public cloud (preferably in Google Cloud Engine)

- Following tools are permitted:
    - https://github.com/labstack/echo
    - https://github.com/go-resty/resty 
    - https://github.com/buger/jsonparser 
    - https://github.com/jinzhu/gorm
 
## Building instructions

### Local 

- Install and configure the latest version of [Golang](https://golang.org/dl/) specific to your development environment.

- Building steps:
    - > git clone https://github.com/andrewkandzuba/openexchange-go-data-gov
    - > cd openexchange-go-data-gov
    - > ./build.sh
 
### Google Cloud Build 

TBD
 
## Components

- [feed/commerce](pkg/feed/commerce.go) - Commerce API News Feed puller and streamer.