Note: Zinc and all its APIs are considered to be alpha stage at this time. Expect breaking changes in API contracts and data format before at this stage.
# Zinc Search Engine

Zinc is a search engine that does full text indexing. It is a lightweight alternative to Elasticsearch and runs using a fraction of the resources. It uses [bluge](https://github.com/blugelabs/bluge) as the underlying indexing library.

It is very simple and easy to operate as opposed to Elasticsearch which requires a couple dozen knobs to understand and tune. 

It is a drop-in replacement for Elasticsearch if you are just ingesting data using APIs and searching using kibana (Kibana is not supported with zinc. Zinc provides its own UI).

Check the below video for a quick demo of Zinc.

[![Zinc Youtube](./screenshots/zinc-youtube.jpg)](https://www.youtube.com/watch?v=aZXtuVjt1ow)

Join slack channel

[![Slack](./screenshots/slack.png)](https://join.slack.com/t/zinc-nvh4832/shared_invite/zt-10a4jb2nl-tQUWwVQgylFEImicA7Fw6A)

# Why zinc

  While Elasticsearch is a very good product, it is complex and requires lots of resources and is more than a decade old. I built Zinc so it becomes easier for folks to use full text search indexing without doing a lot of work.

# Features:

1. Provides full text indexing capability
2. Single binary for installation and running. Binaries available under releases for multiple platforms.
3. Web UI for querying data written in Vue
4. Compatibility with Elasticsearch APIs for ingestion of data (single record and bulk API)
5. Out of the box authentication
6. Schema less - No need to define schema upfront and different documents in the same index can have different fields.

# Roadmap items:
1. Index storage in s3
1. High Availability
1. Distributed reads and writes
1. Geosptial search

# Screenshots

## Search screen
![Search screen 1](./screenshots/search_screen.jpg)
![Search screen for games](./screenshots/search_screen_paris.jpg)

## User management screen
![Users screen](./screenshots/users_screen.jpg)

# Getting started


## Download / Installation / Run

### Windows 

Binaries can be downloaded from [releases](https://github.com/prabhatsharma/zinc/releases) page for appropriate platform.

```shell
C:\> set FIRST_ADMIN_USER=admin
C:\> set FIRST_ADMIN_PASSWORD=Complexpass#123
C:\> mkdir data
C:\> zinc.exe
```
### MacOS - Homebrew 

> $ brew tap prabhatsharma/tap

> $ brew install prabhatsharma/tap/zinc

> $ mkdir data

> $ FIRST_ADMIN_USER=admin FIRST_ADMIN_PASSWORD=Complexpass#123 zinc 

Now point your browser to http://localhost:4080 and login

### MacOS/Linux Binaries
Binaries can be downloaded from [releases](https://github.com/prabhatsharma/zinc/releases) page for appropriate platform.

Create a data folder that will store the data
> $ mkdir data

> $ FIRST_ADMIN_USER=admin FIRST_ADMIN_PASSWORD=Complexpass#123 ./zinc 

Now point your browser to http://localhost:4080 and login

### Docker

> $ mkdir data

> $ docker run -v /full/path/of/data:/data -e DATA_PATH="/data" -p 4080:4080 -e FIRST_ADMIN_USER=admin -e FIRST_ADMIN_PASSWORD=Complexpass#123 --name zinc hiprabhat/zinc:0.1.3

Now point your browser to http://localhost:4080 and login

### Kubernetes

#### Manual Install

Create the namespace:
> $ kubectl create ns zinc

> $ kubectl apply -f k8s/kube-deployment.yaml

> $ kubectl -n zinc port-forward svc/z 4080:4080

Now point your browser to http://localhost:4080 and login

#### Helm

Update Helm values located in [values.yaml](helm/zinc/values.yaml)

Create the namespace:
> $ kubectl create ns zinc

Install the chart:
> $ helm install zinc helm/zinc -n zinc

Zinc can be available with an ingress or port-forward:
> $ kubectl -n zinc port-forward svc/zinc 4080:4080

## Data ingestion

### Single record

curl example

> $ curl \
  -u admin:Complexpass#123 \
  -XPUT \
  -d '{"author":"Prabhat Sharma"}' \
  http://localhost:4080/api/myshinynewindex/document

Python example

```py
import base64, json
import requests

user = "admin"
password = "Complexpass#123"
bas64encoded_creds = base64.b64encode(bytes(user + ":" + password, "utf-8")).decode("utf-8")


data = {
    "Athlete": "DEMTSCHENKO, Albert",
    "City": "Turin",
    "Country": "RUS",
    "Discipline": "Luge",
    "Event": "Singles",
    "Gender": "Men",
    "Medal": "Silver",
    "Season": "winter",
    "Sport": "Luge",
    "Year": 2006
  }

headers = {"Content-type": "application/json", "Authorization": "Basic " + bas64encoded_creds}
index = "games3"
zinc_host = "http://localhost:4080"
zinc_url = zinc_host + "/api/" + index + "/document"

res = requests.put(zinc_url, headers=headers, data=json.dumps(data))

```

### Bulk ingestion

Bulk ingestion API follows same interface as Elasticsearch API defined in [documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html).


> curl -L https://github.com/prabhatsharma/zinc/releases/download/v0.1.1/olympics.ndjson.gz -o olympics.ndjson.gz

> gzip -d  olympics.ndjson.gz 

> curl http://localhost:4080/api/_bulk -i -u admin:Complexpass#123  --data-binary "@olympics.ndjson"


## Search

Python example

```py
import base64
import json
import requests

user = "admin"
password = "Complexpass#123"
bas64encoded_creds = base64.b64encode(
    bytes(user + ":" + password, "utf-8")).decode("utf-8")


params = {
    "search_type": "match",
    "query":
    {
        "term": "DEMTSCHENKO",
        "start_time": "2021-06-02T14:28:31.894Z",
        "end_time": "2021-12-02T15:28:31.894Z"
    },
    "from": 40, # use together with max_results for paginated results.
    "max_results": 20,
    "fields": ["_all"]
}

# params = {
#     "search_type": "querystring",
#     "query":
#     {
#         "term": "+City:Turin +Silver",
#         "start_time": "2021-06-02T14:28:31.894Z",
#         "end_time": "2021-12-02T15:28:31.894Z"
#     },
#     "fields": ["_all"]
# }

headers = {"Content-type": "application/json",
           "Authorization": "Basic " + bas64encoded_creds}
index = "games3"
zinc_host = "http://localhost:4080"
zinc_url = zinc_host + "/api/" + index + "/_search"

res = requests.post(zinc_url, headers=headers, data=json.dumps(params))

print(res.text)

```

Output

```json
{
  "took": 0,
  "timed_out": false,
  "max_score": 7.6978611753656345,
  "hits": {
    "total": {
      "value": 3
    },
    "hits": [
      {
        "_index": "games3",
        "_type": "games3",
        "_id": "bd3e67f0-679b-4aa4-b0f5-81b9dc86a26a",
        "_score": 7.6978611753656345,
        "@timestamp": "2021-10-20T04:56:39.000871Z",
        "_source": {
          "Athlete": "DEMTSCHENKO, Albert",
          "City": "Turin",
          "Country": "RUS",
          "Discipline": "Luge",
          "Event": "Singles",
          "Gender": "Men",
          "Medal": "Silver",
          "Season": "winter",
          "Sport": "Luge",
          "Year": 2006
        }
      },
      {
        "_index": "games3",
        "_type": "games3",
        "_id": "230349d9-72b3-4225-bac7-a8ab31af046d",
        "_score": 7.6978611753656345,
        "@timestamp": "2021-10-20T04:56:39.215124Z",
        "_source": {
          "Athlete": "DEMTSCHENKO, Albert",
          "City": "Sochi",
          "Country": "RUS",
          "Discipline": "Luge",
          "Event": "Singles",
          "Gender": "Men",
          "Medal": "Silver",
          "Season": "winter",
          "Sport": "Luge",
          "Year": 2014
        }
      },
      {
        "_index": "games3",
        "_type": "games3",
        "_id": "338fea31-81f2-4b56-a096-b8294fb6cc92",
        "_score": 7.671309826309841,
        "@timestamp": "2021-10-20T04:56:39.215067Z",
        "_source": {
          "Athlete": "DEMTSCHENKO, Albert",
          "City": "Sochi",
          "Country": "RUS",
          "Discipline": "Luge",
          "Event": "Mixed Relay",
          "Gender": "Men",
          "Medal": "Silver",
          "Season": "winter",
          "Sport": "Luge",
          "Year": 2014
        }
      }
    ]
  },
  "buckets": null,
  "error": ""
}
```


# API reference

These APIs can be used to programatically interact with Zinc.

All APIs must have an authorization header

e.g. Header:

Authorization: Basic YWRtaW46Q29tcGxleHBhc3MjMTIz

## CreateIndex - Create a new index
Endpoint - PUT /api/index 

While you do not need to create indexes manually as they are created automatically when you start ingesting the data, you could create them in advance using this API. S3 backed indexes must be created before they can be used.

e.g. 
PUT http://localhost:4080/api/index

Payload: 

{ "name": "myshinynewindex", "storage_type": "s3" }

OR

{ "name": "myshinynewindex", "storage_type": "disk" }

Default storage_type is disk

## DeleteIndex - Delete an index
Endpoint - DELETE /api/index/:indexName

This will delete the index and its associated metadata. Be careful using this as data is deleted unrecoverably.

e.g. 
DELETE http://localhost:4080/api/index/indextodelete

## ListIndexes - List existing indexes
Endpoint - GET /api/index

Get the list of existing indexes

e.g. 
GET http://localhost:4080/api/index

## UpdateDocument - Create/Update a document and index it for searches
Endpoint - PUT /api/:target/document

Create/Update a document and index it for searches

e.g. 
PUT http://localhost:4080/api/myindex/document

Payload: { "name": "Prabhat Sharma" }

## UpdateDocumentWithId - Create/Update a document and index it for searches. Provide a doc Id
Endpoint - PUT /api/:target/_doc/:id 

Create/Update a document and index it for searches

e.g. 
PUT http://localhost:4080/api/myindex/_doc/1

Payload: { "name": "Prabhat Sharma is meeting friends in San Francisco" }

## DeleteDocument - Delete a document
Endpoint - DELETE /api/:target/_doc/:id

This will delete the document from the index.

e.g. 
DELETE http://localhost:4080/api/myindex/_doc/1

## Search - search for documents
Endpoint - POST /api/:target/_search

Search for documents

e.g. 
POST http://localhost:4080/api/stackoverflow-6/_search

Payload: 
```json
{
    "search_type": "matchphrase",
    "query": {
        "term": "shell window",
        "start_time": "2021-12-25T15:08:48.777Z",
        "end_time": "2021-12-28T16:08:48.777Z"
    },
    "sort_fields": ["-@timestamp"],
    "from": 1,
    "max_results": 20,
    "fields": [
        "*"
    ]
}
```

combine "from" and "max_results" to allow pagination.

sort_fields: list of fields to sort the results. Put a minus "-" before the field to change to descending order.

search_type can have following values:

1. alldocuments
2. wildcard
3. fuzzy
4. term
5. daterange
6. matchall
7. match
8. matchphrase
9. multiphrase
10. querystring


## BulkUpdate - Upload bulk data
Endpoint - POST /api/_bulk

This will upload multiple documents in batch. It is preferred over PUT /api/:target/_doc/:id when you have multiple documents to be inserted as it is magnitude of times faster than uploading individual documents.

You may want to experiment with actual number of documents that you send in a single request. An ideal number may range from 500 to 5000 documents. Log forwaders like fluentbit use this API.

e.g. 
POST /api/_bulk

Payload - ndjson (newline delimited json) content

e.g. ndjson contents

```json
{ "index" : { "_index" : "olympics" } } 
{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
{ "index" : { "_index" : "olympics" } } 
{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}
{ "index" : { "_index" : "olympics" } } 
{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "DRIVAS, Dimitrios", "Country": "GRE", "Gender": "Men", "Event": "100M Freestyle For Sailors", "Medal": "Bronze", "Season": "summer"}
{ "index" : { "_index" : "olympics" } } 
{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "MALOKINIS, Ioannis", "Country": "GRE", "Gender": "Men", "Event": "100M Freestyle For Sailors", "Medal": "Gold", "Season": "summer"}
{ "index" : { "_index" : "olympics" } } 
{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "CHASAPIS, Spiridon", "Country": "GRE", "Gender": "Men", "Event": "100M Freestyle For Sailors", "Medal": "Silver", "Season": "summer"}
```

# Who uses Zinc (Known users)?

1. [Quadrantsec](https://quadrantsec.com/)

Please do raise a PR adding your details if you are using Zinc.



