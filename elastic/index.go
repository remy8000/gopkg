package elastic

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func IndexExists(index string) (bool, error) {
	// Check if index exists
	reqExistIndex := esapi.IndicesExistsRequest{
		Index: []string {index},
		Pretty: true,
	}
	resExistIndex, err := reqExistIndex.Do(context.Background(), es); if err != nil {
		return false, err
	} 
	
	// Securely close Body
	if resExistIndex.Body != nil {
		defer resExistIndex.Body.Close()
	}
	
	if resExistIndex.IsError() {
		return false, nil
	}

	return true, nil
}