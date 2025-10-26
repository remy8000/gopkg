package elastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// returns the json response deserialize into a map[string]interface{}. 
// deserialize in error if errors included in elastic response
func getResponseMap(res *esapi.Response) (map[string]interface{}, error) {

	var r map[string]interface{}
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return r, fmt.Errorf("doRequest(): :error parsing the response body: %s", err)
		} else {
			fmt.Println(e)

			// error handling
			var typ,reason interface{}
			if e["error"] != nil {
				typ = e["error"].(map[string]interface{})["type"]
				reason = e["error"].(map[string]interface{})["reason"]
			}
			
			// failures handling (update by query)
			
			// Print the response status and error information.
			return r, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				typ,
				reason,
			)
		}
		
	}

	// Deserialize the json response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, fmt.Errorf("doRequest(): decoding the response body: %s", err)
	}

	return r, nil
}