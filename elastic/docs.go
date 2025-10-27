package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// returns a slice of map containing one doc
// a response in a form of a slice is needed to apply esdto.ToPosts
func GetDocById(index string, id string, source []string, timeOut int) (map[string]interface{}, error) {

	doc := make(map[string]interface{})
	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return doc, err
	}

	if !exists {
		return doc, fmt.Errorf("no index with name '%s'", index)
	}

	// Set up the request object.
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
		Source:     source,
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return doc, fmt.Errorf("getDocById - request timed out")
		}
		return doc, fmt.Errorf("es error getting response: %s", err)
	}

	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return doc, err
	}

	if r["found"].(bool) {
		doc["id"] = r["_id"].(string)
		doc["source"] = r["_source"].(map[string]interface{})
	}

	return doc, nil
}

// returns a list of docs providing a list of ids and a source set
func GetDocsMultiIds(index string, ids []string, source []string, timeOut int) ([]map[string]interface{}, error) {
	var docs []map[string]interface{}

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return docs, err
	}

	if !exists {
		return docs, fmt.Errorf("no index with name '%s'", index)
	}

	// build body
	type docBody struct {
		Id     string   `json:"_id"`
		Source []string `json:"_source"`
	}

	var docsBody []docBody
	for _, id := range ids {
		d := docBody{
			Id:     id,
			Source: source,
		}
		docsBody = append(docsBody, d)
	}

	// Build the request body.
	body, err := json.Marshal(docsBody)
	if err != nil {
		return docs, err
	}

	// Set up the request object.
	req := esapi.MgetRequest{
		Index: index,
		Body:  bytes.NewReader([]byte(fmt.Sprintf(`{"docs":%s}`, body))),
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return docs, fmt.Errorf("GetDocsMultiIds - request timed out")
		}
		return docs, fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return docs, err
	}

	for _, doc := range r["docs"].([]interface{}) {
		source := doc.(map[string]interface{})["_source"]
		id := doc.(map[string]interface{})["_id"].(string)
		m := make(map[string]interface{})
		m["id"] = id
		m["source"] = source
		docs = append(docs, m)
	}

	return docs, nil
}

func SaveDoc(index string, d Doc, timeOut int) (string,error) {

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("no index with name '%s'", index)
	}

	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(body),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("SaveDoc - request timed out")
		}
		return "", fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return "", err
	}

	id, ok := r["_id"].(string)
    if !ok {
        return "", fmt.Errorf("could not get document ID from response")
    }


	return id, nil
}

// return id
// and result ["update not processed",'update processed", "updated", "noop"]
func UpdateDoc(index string, id string, d Doc, timeOut int) (string, string, error) {

	result := "update not processed"

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return id, result, err
	}

	if !exists {
		return id, result, fmt.Errorf("no index with name '%s'", index)
	}

	body, err := json.Marshal(d)
	if err != nil {
		return id, result, err
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return id, result, fmt.Errorf("request timed out")
		}
		return id, result, fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		result = "update processed"
		return id, result, err
	}

	result = r["result"].(string)

	return id, result, nil
}

// based on update_by_query
// returns count of updated docs
func UpdateByQuery(index string, query map[string]interface{}, timeOut int) (int, error) {

	var updatedCount int

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return updatedCount, err
	}

	if !exists {
		return updatedCount, fmt.Errorf("no index with name '%s'", index)
	}

	// Build the request body.
	body, err := json.Marshal(query)
	if err != nil {
		return updatedCount, err
	}

	// Set up the request object.
	req := esapi.UpdateByQueryRequest{
		Index: []string{index},
		Body:  bytes.NewReader(body),
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return updatedCount, fmt.Errorf("UpdateDoc - request timed out")
		}
		return updatedCount, fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return updatedCount, err
	}

	updatedCount = int(r["updated"].(float64))

	return updatedCount, nil
}

func DeleteDoc(index string, id string, timeOut int) error {

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("no index with name '%s'", index)
	}

	// Set up the request object.
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("request timed out")
		}
		return fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return err
	}

	result := r["result"].(string)

	fmt.Printf("Doc id=%v in index=%v result query=%v", id, index, result)

	return nil
}

func DeleteByQuery(index string, query map[string]interface{}, timeOut int) (int, error) {
	var deletedCount int

	// CHECKS
	exists, err := IndexExists(index)
	if err != nil {
		return deletedCount, err
	}

	if !exists {
		return deletedCount, fmt.Errorf("no index with name '%s'", index)
	}

	// Build the request body.
	body, err := json.Marshal(query)
	if err != nil {
		return deletedCount, err
	}

	// Set up the request object.
	req := esapi.DeleteByQueryRequest{
		Index: []string{index},
		Body:  bytes.NewReader(body),
	}

	// Set up a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	// Perform the request with the client.
	res, err := req.Do(ctx, es)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return deletedCount, fmt.Errorf("request timed out")
		}
		return deletedCount, fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	r, err := getResponseMap(res)
	if err != nil {
		return deletedCount, err
	}

	if r["deleted"] != nil {
		deletedCount = int(r["deleted"].(float64))
	}

	return deletedCount, nil
}
