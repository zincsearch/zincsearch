Note: Zinc and all its APIs are considered to be alpha stage at this time. Expect breaking changes in API contracts and data format at this stage.
# Zinc Search Engine

Zinc is a search engine that does full text indexing. It is a lightweight alternative to Elasticsearch and runs using a fraction of the resources. It uses [bluge](https://github.com/blugelabs/bluge) as the underlying indexing library.

It is very simple and easy to operate as opposed to Elasticsearch which requires a couple dozen knobs to understand and tune. 

It is a drop-in replacement for Elasticsearch if you are just ingesting data using APIs and searching using kibana (Kibana is not supported with zinc. Zinc provides its own UI).

Check the below video for a quick demo of Zinc.

[![Zinc Youtube](./screenshots/zinc-youtube.jpg)](https://www.youtube.com/watch?v=aZXtuVjt1ow)

Join slack channel

[![Slack](./screenshots/slack.png)](https://join.slack.com/t/zinc-nvh4832/shared_invite/zt-11r96hv2b-UwxUILuSJ1duzl_6mhJwVg)

# Why zinc

  While Elasticsearch is a very good product, it is complex and requires lots of resources and is more than a decade old. I built Zinc so it becomes easier for folks to use full text search indexing without doing a lot of work.

# Features:

1. Provides full text indexing capability
2. Single binary for installation and running. Binaries available under releases for multiple platforms.
3. Web UI for querying data written in Vue
4. Compatibility with Elasticsearch APIs for ingestion of data (single record and bulk API)
5. Out of the box authentication
6. Schema less - No need to define schema upfront and different documents in the same index can have different fields.
7. Index storage in s3 (experimental)
8. aggregation support

# Roadmap items:
1. High Availability
1. Distributed reads and writes
1. Geosptial search
1. Raise an issue if you are looking for something.

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

------------------------
**Optional - Only if you have AWS CLI installed.**

If you have AWS CLI installed amd get login error then run below command:

> aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

------------------------

Docker images are available at https://gallery.ecr.aws/prabhat/zinc

> $ mkdir data

> $ docker run -v /full/path/of/data:/data -e DATA_PATH="/data" -p 4080:4080 -e FIRST_ADMIN_USER=admin -e FIRST_ADMIN_PASSWORD=Complexpass#123 --name zinc public.ecr.aws/prabhat/zinc:latest



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
Data ingestion can be done using APIs and log forwarders like fluent-bit and syslog-ng. Check [docs](https://docs.zincsearch.io/ingestion/bulk-ingestion/#bulk-ingestion) for details.

## API Reference

Check [docs](https://docs.zincsearch.io/API%20Reference/)


# Who uses Zinc (Known users)?

1. [Quadrantsec](https://quadrantsec.com/)

Please do raise a PR adding your details if you are using Zinc.



