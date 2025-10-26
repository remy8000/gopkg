package elastic

import (
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

var es *elasticsearch.Client

func Setup(cfg elasticsearch.Config) error {
	var err error
	es, err = elasticsearch.NewClient(cfg); if err != nil {
		return err
	}
	return nil
}

type Doc interface {
	IsDoc()
}

func ES() *elasticsearch.Client {
	return es
}

// get cluster info return client and server version
func ClusterInfo() (string, string, error) {
	res, err := es.Info()
	if err != nil {
		return "", "", err
	}
	
	// Securely close Body
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Check response status
	if res.IsError() {
		return "", "", fmt.Errorf("es res error: %s", res.String())
	}


	// Deserialize the response into a map.
	var r  map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", "", fmt.Errorf("error parsing the response body: %s", err)
	}
	
	return elasticsearch.Version, fmt.Sprintf("%v",r["version"].(map[string]interface{})["number"]), nil
}