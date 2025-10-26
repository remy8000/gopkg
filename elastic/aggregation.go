package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// this functions makes a query including a aggregation
// the name has to be specified in order to retrieve the contained docs in the buckets
// returns a map[bucket_name] = []docs
func Aggregation(index, aggregationName string, query map[string]interface{}) (map[string][]map[string]interface{}, error) {

	bucketsMap := make(map[string][]map[string]interface{})

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return bucketsMap, fmt.Errorf("index exist err: %s", err)
	}

	if !exists {
		return bucketsMap, fmt.Errorf("index %v doesn't exist or index not included in elastic role for this user", index)
	}

	body, err := json.Marshal(query)
	if err != nil {
		return bucketsMap, err
	}

	// Set up the request object.
	req := esapi.SearchRequest{
		Index:          []string{index},
		Body:           bytes.NewReader(body),
		TrackTotalHits: true,
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return bucketsMap, fmt.Errorf("es error getting response: %s", err)
	}

	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return bucketsMap, err
	}

	// get aggs buckets with safe type assertions
	aggs, ok := r["aggregations"].(map[string]interface{})
	if !ok {
		return bucketsMap, fmt.Errorf("missing 'aggregations' field in response")
	}
	agg, ok := aggs[aggregationName].(map[string]interface{})
	if !ok {
		return bucketsMap, fmt.Errorf("missing aggregation '%s' in response", aggregationName)
	}
	buckets, ok := agg["buckets"].(map[string]interface{})
	if !ok {
		return bucketsMap, fmt.Errorf("missing or invalid 'buckets' field in aggregation '%s'", aggregationName)
	}

	for key, element := range buckets {
		elemMap, ok := element.(map[string]interface{})
		if !ok {
			return bucketsMap, fmt.Errorf("invalid bucket element for key '%s'", key)
		}
		hitsOuter, ok := elemMap["hits"].(map[string]interface{})
		if !ok {
			return bucketsMap, fmt.Errorf("missing 'hits' in bucket '%s'", key)
		}
		hitsInner, ok := hitsOuter["hits"].([]interface{})
		if !ok {
			return bucketsMap, fmt.Errorf("missing or invalid 'hits' array in bucket '%s'", key)
		}
		var docsMap []map[string]interface{}
		for _, h := range hitsInner {
			doc, ok := h.(map[string]interface{})
			if ok {
				docsMap = append(docsMap, doc)
			}
		}
		bucketsMap[key] = docsMap
	}
	return bucketsMap, nil
}
