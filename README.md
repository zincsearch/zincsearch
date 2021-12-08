Note: Zinc and all its APIs are considered to be alpha stage at this time.
# Zinc

Zinc is a search engine that does full text indexing. It is a lightweight alternative to Elasticsearch and runs in less than 100 MB of RAM. It uses [bluge](https://github.com/blugelabs/bluge) as the underlying indexing library.

It is very simple and easy to operate as opposed to Elasticsearch which requires a couple dozen knobs to understand and tune. 

It is a drop-in replacement for Elasticsearch if you are just ingesting data using APIs and searching using kibana (Kibana is not supported with zinc. Zinc provides its own UI).

# Why zinc

  While Elasticsearch is a very good product, it is complex and requires lots of resources and is more than a decade old. I built Zinc so it becomes easier for folks to use full text search indexing without doing a lot of work.

# Features:

1. Provides full text indexing capability
2. Single binary for installation and running. Binaries available under releases for multiple platforms.
3. Web UI for querying data written in Vue
4. Compatibility with Elasticsearch APIs for ingestion of data (single record and bulk API)
5. Out of the box authentication
6. Schema less - No need to define schema upfront and different documents in the same index can have different fields.

# Missing features:
1. Clustering and High Availability


# Screenshots

## Search screen
![Search screen 1](./screenshots/search_screen.jpg)
![Search screen for games](./screenshots/search_screen_paris.jpg)

## User management screen
![Users screen](./screenshots/users_screen.jpg)

# Getting started


## Download / Installation / Run

### Binaries
Binaries can be downloaded from [releases](https://github.com/prabhatsharma/zinc/releases) page for appropriate platform.

Create a data folder that will store the data
> $ mkdir data

> $ FIRST_ADMIN_USER=admin FIRST_ADMIN_PASSWORD=Complexpass#123 zinc 

Now point your browser to http://localhost:4080 and login


### Docker

> $ mkdir data

> $ docker run -v /full/path/of/data:/data -e DATA_PATH="/data" -p 4080:4080 -e FIRST_ADMIN_USER=admin -e FIRST_ADMIN_PASSWORD=Complexpass#123 --name zinc public.ecr.aws/m5j1b6u0/zinc:v0.1.1 

Now point your browser to http://localhost:4080 and login

### Kubernetes

#### Manual Install

> kubectl apply -f kube-deployment.yaml

> kubectl -n zinc port-forward svc/z 4080:4080

Now point your browser to http://localhost:4080 and login

#### Helm

Update Helm values located in [values.yaml](helm/zinc/values.yaml)

Create the namespace:
> kubectl create ns zinc

Install the chart:
> helm install zinc helm/zinc -n zinc

Zinc can be available with an ingress or port-forward:
> kubectl -n zinc port-forward svc/zinc 4080:4080

## Data ingestion

### Single record

python example

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


> curl -L https://github.com/prabhatsharma/zinc/releases/download/v0.1.1/https://github.com/prabhatsharma/zinc/releases/download/v0.1.1/olympics.ndjson.gz -o olympics.ndjson.gz

> gzip -d  olympics.ndjson.gz 

> curl http://localhost:4080/es/olympics/_bulk -i -u admin:Complexpass#123  --data-binary "@olympics.ndjson"


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

output

```json
{"took":0,"timed_out":false,"max_score":7.6978611753656345,"hits":{"total":{"value":3},"hits":[{"_index":"games3","_type":"games3","_id":"bd3e67f0-679b-4aa4-b0f5-81b9dc86a26a","_score":7.6978611753656345,"@timestamp":"2021-10-20T04:56:39.000871Z","_source":{"Athlete":"DEMTSCHENKO, Albert","City":"Turin","Country":"RUS","Discipline":"Luge","Event":"Singles","Gender":"Men","Medal":"Silver","Season":"winter","Sport":"Luge","Year":2006}},{"_index":"games3","_type":"games3","_id":"230349d9-72b3-4225-bac7-a8ab31af046d","_score":7.6978611753656345,"@timestamp":"2021-10-20T04:56:39.215124Z","_source":{"Athlete":"DEMTSCHENKO, Albert","City":"Sochi","Country":"RUS","Discipline":"Luge","Event":"Singles","Gender":"Men","Medal":"Silver","Season":"winter","Sport":"Luge","Year":2014}},{"_index":"games3","_type":"games3","_id":"338fea31-81f2-4b56-a096-b8294fb6cc92","_score":7.671309826309841,"@timestamp":"2021-10-20T04:56:39.215067Z","_source":{"Athlete":"DEMTSCHENKO, Albert","City":"Sochi","Country":"RUS","Discipline":"Luge","Event":"Mixed Relay","Gender":"Men","Medal":"Silver","Season":"winter","Sport":"Luge","Year":2014}}]},"buckets":null,"error":""}
```

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

# Who uses Zinc ?

1. [Quadrantsec](https://quadrantsec.com/)

Please do raise a PR adding your details if you are using Zinc.



