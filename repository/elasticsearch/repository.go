package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/danielhoward314/hexagonal/shortener"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pkg/errors"
)

type esRepository struct {
	client    *elasticsearch.Client
	indexName string
}

type response struct {
	Source json.RawMessage `json:"_source"`
}

func (r *esRepository) CreateIndex(mapping string) error {
	res, err := r.client.Indices.Create(r.indexName, r.client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func NewEsRepository() (shortener.RedirectRepository, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPostgresRepo")
	}
	mapping := `{
		"mappings": {
		  "_doc": {
			"properties": {
			  "code":	{ "type": "keyword" },
			  "url":	{ "type": "keyword" },
			  "created_at":	{ "type": "integer" }
			}
		  }
		}
	}`
	repo := &esRepository{client: client, indexName: "redirects"}
	repo.CreateIndex(mapping)
	return repo, nil
}

func (r *esRepository) Find(code string) (*shortener.Redirect, error) {
	res, _ := esapi.GetRequest{
		Index: r.indexName, DocumentID: code,
	}.Do(context.Background(), r.client)
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	response := &response{}
	e := json.NewDecoder(res.Body)
	e.Decode(response)
	redirect := &shortener.Redirect{}
	if err := json.Unmarshal(response.Source, redirect); err != nil {
		return nil, err
	}
	return redirect, nil
}

func (r *esRepository) Store(redirect *shortener.Redirect) error {
	payload, err := json.Marshal(redirect)
	if err != nil {
		return err
	}
	res, _ := esapi.CreateRequest{
		Index:      r.indexName,
		DocumentID: redirect.Code,
		Body:       bytes.NewReader(payload),
	}.Do(context.Background(), r.client)
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}
