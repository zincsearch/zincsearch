package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CatTemplate(c *gin.Context) {
	name := c.Param("name")
	log.Print("CatTemplate:", name)
}

func IngestPipeline(c *gin.Context) {
	name := c.Param("name")
	log.Print("IngestPipeline:", name)
}

func IndexTemplate(c *gin.Context) {
	name := c.Param("name")
	log.Print("IndexTemplate:", name)
}

// Ping is used by filebeat to check a connection. It must be at root of url.
// If filebeat is configured with "localhost:4080/es" then this responds to
// that.
func Ping(c *gin.Context) {
	io.WriteString(c.Writer, `{
		"name" : "NA",
		"cluster_name" : "NA",
		"cluster_uuid" : "NA",
		"version" : {
		  "number" : "0.1.1-zinc",
		  "build_flavor" : "default",
		  "build_type" : "NA",
		  "build_hash" : "NA",
		  "build_date" : "2021-12-12T20:18:09.722761972Z",
		  "build_snapshot" : false,
		  "lucene_version" : "NA",
		  "minimum_wire_compatibility_version" : "NA",
		  "minimum_index_compatibility_version" : "NA"
		},
		"tagline" : "You Know, for Search"
	  }
	  `)
}

func License(c *gin.Context) {
	/*
			{
		  "license" : {
		    "status" : "active",
		    "uid" : "5823bd41-139e-4e3c-93fd-499ce05576a1",
		    "type" : "basic",
		    "issue_date" : "2021-12-12T19:55:48.391Z",
		    "issue_date_in_millis" : 1639338948391,
		    "max_nodes" : 1000,
		    "issued_to" : "docker-cluster",
		    "issuer" : "elasticsearch",
		    "start_date_in_millis" : -1
		  }
		}
	*/
	io.WriteString(c.Writer, `{"license":{"status":"active"}}`)
}

func XPackHandler(c *gin.Context) {
	xpr := XPackResponse{}
	c.JSON(200, xpr)
}

type XPackResponse struct {
	Build struct {
		Date string `json:"date"`
		Hash string `json:"hash"`
	} `json:"build"`
	Features struct {
		AggregateMetric struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"aggregate_metric"`
		Analytics struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"analytics"`
		Ccr struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"ccr"`
		DataStreams struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"data_streams"`
		DataTiers struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"data_tiers"`
		Enrich struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"enrich"`
		Eql struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"eql"`
		FrozenIndices struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"frozen_indices"`
		Graph struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"graph"`
		Ilm struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"ilm"`
		Logstash struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"logstash"`
		Ml struct {
			Available      bool `json:"available"`
			Enabled        bool `json:"enabled"`
			NativeCodeInfo struct {
				BuildHash string `json:"build_hash"`
				Version   string `json:"version"`
			} `json:"native_code_info"`
		} `json:"ml"`
		Monitoring struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"monitoring"`
		Rollup struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"rollup"`
		SearchableSnapshots struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"searchable_snapshots"`
		Security struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"security"`
		Slm struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"slm"`
		Spatial struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"spatial"`
		Sql struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"sql"`
		Transform struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"transform"`
		Vectors struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"vectors"`
		VotingOnly struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"voting_only"`
		Watcher struct {
			Available bool `json:"available"`
			Enabled   bool `json:"enabled"`
		} `json:"watcher"`
	} `json:"features"`
	License struct {
		Mode   string `json:"mode"`
		Status string `json:"status"`
		Type   string `json:"type"`
		UID    string `json:"uid"`
	} `json:"license"`
	Tagline string `json:"tagline"`
}
