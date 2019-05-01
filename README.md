# elastic_go
elasticsearch &amp; go query


Log JSON structure snippet:
```json
 "hits" : {
    "total" : <total>,
    "max_score" : null,
    "hits" : [
      {
        "_index" : "filebeat-<version>",
        "_type" : "doc",
        "_id" : "<id>",
        "_score" : null,
        "_source" : {
          "@timestamp" : "<timestamp>",
          "input" : {
            "type" : "log"
          },
          "prospector" : {
            "type" : "log"
          },
          "beat" : {
            "name" : "<name>",
            "hostname" : "<hostname>",
            "version" : "<version>"
          },
          "host" : {
            "name" : "<name>",
            "id" : "<id>",
            "containerized" : true,
            "architecture" : "x86_64",
            "os" : {
              "version" : "<version>",
              "family" : "",
              "codename" : "<codename>",
              "platform" : "<platform>"
            }
          },
          "message" : "<message>",
          "source" : "<source>",
          "offset" : <offset>,
          "log" : {
            "flags" : [
              "multiline"
            ]
          }
        },
        "sort" : [
          <sort>
        ]
      },
```
