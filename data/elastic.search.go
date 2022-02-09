package data

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ElasticSearch *elasticsearch.Client

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	config := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	var err error
	ElasticSearch, err = elasticsearch.NewClient(config)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := ElasticSearch.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	log.Println(Green, res, Reset)
}

func GetInfo() {
	log.Printf("Client: %s", elasticsearch.Version)
}
