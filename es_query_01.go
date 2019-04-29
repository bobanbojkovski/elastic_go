// sample to query a search term and filter for last 30 min in Elasticsearch
// reference: https://olivere.github.io/elastic/

package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gopkg.in/olivere/elastic.v<version>"                             //<----- version (v7)
)

// Log represents data structure in Elasticsearch
type Log struct {
	Message string `json:"message"`
}

const (
	hostNamePort = "http://<hostname>:9200"                           //<----- hostname
	indexName    = "filebeat-*"
)

func findTerm(client *elastic.Client) error {
	ctx := context.Background()

	messageField := "message"
	timestampField := "@timestamp"
	term := "<search_term>"                                           //<----- search term

	mins := 30
	lastMins := time.Now().Add(time.Duration(-mins) * time.Minute)

	// compose the query & search
	matchQuery := elastic.NewMatchQuery(messageField, term)
	query := elastic.NewBoolQuery().Must(matchQuery).
		Filter(elastic.NewRangeQuery(timestampField).
			Gte(lastMins))

	res, err := client.Search().
		Index(indexName).
		Query(query).
		Sort(timestampField, false).
		Size(10000).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	// Print message containing the term
	var l Log
	for _, item := range res.Each(reflect.TypeOf(l)) {
		if l, ok := item.(Log); ok {
			fmt.Println(l.Message)
		}
	}

	// Print number of hits containing the term
	fmt.Printf("%s: %d \n", term, res.TotalHits())

	return nil
}

func main() {
	client, err := elastic.NewClient(
		elastic.SetURL(hostNamePort),
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	err = findTerm(client)
	if err != nil {
		panic(err)
	}
}
