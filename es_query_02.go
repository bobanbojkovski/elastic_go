// sample to query a search term and filter for last 30 min in Elasticsearch
// reference: https://github.com/elastic/go-elasticsearch

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch"
)

func main() {
	log.SetFlags(0)

	hostNamePort := "<hostname>:<port>"
	indexName := "filebeat-*"
	messageField := "message"
	timestampField := "@timestamp"
	term := "<search_term>"
	lastMins := "now-30m"

	var (
		r map[string]interface{}
	)

	cfg := elasticsearch.Config{
		Addresses: []string{
			hostNamePort,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Search for the indexed documents
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithSize(10000),
		es.Search.WithFilterPath("hits", "hits.hits._source"),
		es.Search.WithQuery(messageField+" : "+term+" AND "+timestampField+" :["+lastMins+" TO now]"),
		es.Search.WithSort(timestampField+":desc"),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Print the message containing the search term
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf("%s", hit.(map[string]interface{})["_source"].(map[string]interface{})[messageField])
	}

	// Print the number of results
	log.Printf(
		term+": %d",
		int(r["hits"].(map[string]interface{})["total"].(float64)),
	)

}
