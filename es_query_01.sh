#!/usr/bin/env bash

# sample using query DSL to search for a term and filter the last 30 min in Elasticsearch
# reference: https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html

set -o errexit
set -o pipefail
set -o nounset


es_query () {

    curl --silent --request GET 'http://'$1:$2'/_search?pretty' -H 'Content-Type: application/json' -d '
    {
        "from": 0, "size": 10000,
        "query":{  
           "bool":{  
              "must":[  
                 {  
                    "query_string":{  
                       "query":"'$3'",
                       "analyze_wildcard":true,
                       "default_field":"*"
                    }
                 },
                 {  
                    "range":{  
                       "'$4'":{  
                          "gte":"now-'$5'm"
                       }
                    }
                 }
              ],
              "filter":[],
              "should":[],
              "must_not":[]
           }
        }
    }'

}


main () {

    local STATUS_SUCCESS="0"
    local STATUS_WARNING="1"
    local STATUS_CRITICAL="2"
    local STATUS_UNKNOWN="3"

    # UPDATE THE DEFAULT VALUES
    local hostname="localhost"
    local port="9200"
    local term="<DEFAULT_TERM>"
    local timestampField="@timestamp"
    local warn="2"
    local crit="3"
    local mins="30"

    while test -n "${1-}"; do
        case "$1" in
            --hostname|-H)
                hostname="$2"
                shift
                ;;
            --port|-P)
                port="$2"
                shift
                ;;
            --term|-t)
                term="$2"
                shift
                ;;
            -w)
                warn="$2"
                shift
                ;;
            -c)
                crit="$2"
                shift
                ;;
            *)
                echo "Unknown argument: "$1""
                echo "Syntax: ./check_es_term.sh -H <hostname> -P <port> -t <term> -w <waring_threshold> -c <critical treshold>"
                exit "${STATUS_UNKNOWN}"
                ;;
        esac
        shift
    done

    query=$(es_query "$hostname" "$port" "$term" "$timestampField" "$mins")
    total_hits=$(echo $query | jq '.hits.total')

    if [[ -z "${total_hits}" ]]; then
        echo "UNKNOWN - Number of total hits of "${term}" is not known :) Check es accessibility, liveness check"
        exit "${STATUS_UNKNOWN}"
    fi

    if [[ "${total_hits}" -lt "${warn}" ]]; then
        echo "OK - Number of total hits of "${term}" is "${total_hits}""
        exit "${STATUS_SUCCESS}"
    fi


    if [[ "${total_hits}" -gt "${warn}" ]] && [[ "${total_hits}" -lt "${crit}" ]]; then
        echo "WARNING - Number of total hits of "${term}" is "${total_hits}""
        exit "${STATUS_WARNING}"
    fi


    if [[ "${total_hits}" -ge "${crit}" ]]; then
        echo "CRITICAL - Number of total hits of "${term}" is "${total_hits}""
        exit "${STATUS_CRITICAL}"
    fi

}

main "${@:-noargs}"
