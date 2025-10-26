package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// search api
// return hits, total, err
func Search(indices []string, query map[string]interface{}, timeOut int) ([]map[string]interface{}, int, error) {

	var hits []map[string]interface{}
	var total int

	// CHECKS
	for _, index := range indices {
		exists, err := IndexExists(index)
		if err != nil {
			return hits, total, fmt.Errorf("index exist err: %s", err)
		}

		if !exists {
			return hits, total, fmt.Errorf("index %v doesn't exist or index not included in elastic role for this user", index)
		}
	}

	body, err := json.Marshal(query)
	if err != nil {
		return hits, total, err
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Set up the request object.
	req := esapi.SearchRequest{
		Index:          indices,
		Body:           bytes.NewReader(body),
		TrackTotalHits: true,
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		return hits, total, fmt.Errorf("es error getting response: %s", err)
	}

	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return hits, total, err
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		id := hit.(map[string]interface{})["_id"]
		highlight := hit.(map[string]interface{})["highlight"]
		index := hit.(map[string]interface{})["_index"]
		m := make(map[string]interface{})
		m["id"] = id
		m["index"] = index
		m["source"] = source
		m["highlight"] = highlight
		hits = append(hits, m)
	}
	// Print the response status, number of results, and request duration.
	total = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	return hits, total, nil
}
