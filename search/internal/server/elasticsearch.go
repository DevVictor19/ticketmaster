package server

import (
	"log"

	"github.com/elastic/go-elasticsearch/v9"
)

func getEsClient(cfg *Elasticsearch) *elasticsearch.Client {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Addr},
	})
	if err != nil {
		log.Fatalf("failed to create elasticsearch client: %v\n", err)
	}
	return esClient
}
