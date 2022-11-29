[![Go Report Card](https://goreportcard.com/badge/github.com/zinclabs/zinc)](https://goreportcard.com/report/github.com/zinclabs/zinc)
[![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/zincsearch/shared_invite/zt-11r96hv2b-UwxUILuSJ1duzl_6mhJwVg) [![Docs](https://img.shields.io/badge/Docs-Docs-green)](https://docs.zincsearch.com/) [![codecov](https://codecov.io/github/zinclabs/zinc/branch/main/graph/badge.svg)](https://codecov.io/github/zinclabs/zinc)

# Zinc Search Engine

Zinc is a search engine that does full text indexing. It is a lightweight alternative to Elasticsearch and runs using a fraction of the resources. It uses [bluge](https://github.com/blugelabs/bluge) as the underlying indexing library.

It is very simple and easy to operate as opposed to Elasticsearch which requires a couple dozen knobs to understand and tune which you can get up and running in 2 minutes

It is a drop-in replacement for Elasticsearch if you are just ingesting data using APIs and searching using kibana (Kibana is not supported with zinc. Zinc provides its own UI).

Check the below video for a quick demo of Zinc.

[![Zinc Youtube](./screenshots/zinc-youtube.jpg)](https://www.youtube.com/watch?v=aZXtuVjt1ow)

# Playground Server

You could try ZincSearch without installing using below details: 

|          |                                        |
-----------|-----------------------------------------
| Server   | https://playground.dev.zincsearch.com  |
| User ID  | admin                                  |
| Password | Complexpass#123                        |

Note: Do not store sensitive data on this server as its available to everyone on internet. Data will also be cleaned on this server regularly.

# Why zinc

  While Elasticsearch is a very good product, it is complex and requires lots of resources and is more than a decade old. I built Zinc so it becomes easier for folks to use full text search indexing without doing a lot of work.

# Features:

1. Provides full text indexing capability
2. Single binary for installation and running. Binaries available under releases for multiple platforms.
3. Web UI for querying data written in Vue
4. Compatibility with Elasticsearch APIs for ingestion of data (single record and bulk API)
5. Out of the box authentication
6. Schema less - No need to define schema upfront and different documents in the same index can have different fields.
7. Index storage in disk (default), s3 or minio (experimental)
8. aggregation support

# Roadmap items:

Public roadmap is available at https://github.com/orgs/zinclabs/projects/3/views/1

Please create an issue if you would like something to be added to the roadmap.

# Screenshots

## Search screen
![Search screen](./screenshots/search_screen.jpg)

## User management screen
![Users screen](./screenshots/users_screen.jpg)

# Getting started


## Quickstart

Check [Quickstart](https://docs.zincsearch.com/quickstart/) 


# Releases

ZincSearch currently has most of its API contracts frozen. It's data format may still experience changes as we improve things. Currently ZincSearch is in beta. Data format should become highly stable when we move to GA (version 1).

# How to develop and contribute to Zinc

Check the [contributing guide](./CONTRIBUTING.md) . Also check the [roadmap items](https://github.com/orgs/zinclabs/projects/3)


# Editions

| Feature             | Zinc      |   Zinc Cloud                      |
----------------------|-----------|-----------------------------------|
| Ideal use case      | App search| Logs and Events (Immutable Data)  | 
| Storage             | Disk      |  Object (S3), GCS, Azure blob coming soon   |
| Preferred Use case  | App search | Log / event search |
| Max  data supported | 100s of GBs | Petabyte scale |
| High availability   | Will be available soon | Yes |
| Open source         | Yes | Willl soon be available |
| ES API compatibility| Search and Ingestion | Ingestion only | 
| GUI                 | Basic     | Advanced for log search |
| Cost                | Free (self hosting may cost money based on size)| Generous free tier. 1 TB ingest / month free.| 
| Get started         | [Quick start](https://docs.zincsearch.com/quickstart/) | [![Sign up](./screenshots/get-started-for-free.png)](https://cloud.zincsearch.com) |



