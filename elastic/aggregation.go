package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Aggregation makes a query including an aggregation
// The aggregation name has to be specified in order to retrieve the buckets with doc counts
// Returns a map[bucket_name] = []map with doc_count
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

	// Deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return bucketsMap, err
	}

	// Extract aggregation buckets
	aggs, found := r["aggregations"].(map[string]interface{})[aggregationName].(map[string]interface{})["buckets"]
	if !found {
		return bucketsMap, fmt.Errorf("aggregation %s not found in response", aggregationName)
	}

	buckets := aggs.([]interface{})
	for _, element := range buckets {
		bucket := element.(map[string]interface{})
		key := bucket["key"].(string)
		docCount := bucket["doc_count"].(float64)

		// Add each bucket with its document count
		bucketsMap[key] = []map[string]interface{}{
			{"doc_count": docCount},
		}
	}

	return bucketsMap, nil
}
