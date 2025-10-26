package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func DataStreamSaveDoc(alias string, d Doc) error {

	body, err := json.Marshal(d); if err != nil {
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   alias,
		Body:    bytes.NewReader(body),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("es error getting response: %s", err)
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	//  deserialize response and possible errors
	_, err = getResponseMap(res); if err != nil {
		return err
	}

	return nil
}