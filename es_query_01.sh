#!/bin/bash

# sample using query DSL to search for a term and filter the last 30 min in Elasticsearch
# reference: https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html

hostNamePort="<hostname>:<port>"      #<--- hostname & port
term="<search_term>"                  #<--- search term
timestampField="@timestamp"
mins=30

query=$(curl -s -X GET 'http://'$hostNamePort'/filebeat-*/_search?pretty' -H 'Content-Type: application/json' -d '
{
    "from": 0, "size": 10000,
    "query": {
        "bool": {
            "must": {
                "match": {
                    "message": {
                        "query": "'$term'"
                    }
                }
            },
            "filter": {
                "range": {
                    "'$timestampField'": {
                        "gte": "now-'$mins'm"
                    }
                }
            }
        }
    },
    "sort": [
      {
        "'$timestampField'": {
          "order": "desc"
        }
      }
    ]
}' 
)

echo $query | jq -r '.hits.hits[]._source.message'

echo -n "$term: "; echo $query | jq '.hits.total'
