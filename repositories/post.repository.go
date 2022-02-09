package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"es/data"
	"es/entities"
	"es/extensions"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
)

const index = "posts"

var PostRepository IPostRepository = &postRepository{}

type IPostRepository interface {
	FindPostByID(id uuid.UUID) (*entities.Post, error)
	AddPost(id uuid.UUID, body []byte) error
	UpdatePost(id uuid.UUID, body []byte) error
	RemovePost(id uuid.UUID) error
}

type postRepository struct{}

func (*postRepository) FindPostByID(id uuid.UUID) (*entities.Post, error) {
	var entity entities.Post

	request := esapi.GetRequest{
		Index:      index,
		DocumentID: id.String(),
	}

	response, err := request.Do(context.Background(), data.ElasticSearch)
	if err != nil {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error while fetching record from database, %s", err.Error())
	}

	defer response.Body.Close()
	if response.IsError() {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
	}

	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error parsing the response body: %s", err.Error())
	}

	err = extensions.Decode(body["_source"].(map[string]interface{}), &entity)
	if err != nil {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error mapping from _source to entity: %s", err.Error())
	}

	return &entity, nil
}

func (*postRepository) AddPost(id uuid.UUID, body []byte) error {
	// TODO: MAKE ASYNC
	reader := strings.NewReader(string(body))

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: id.String(),
		Body:       reader,
	}

	response, err := request.Do(context.Background(), data.ElasticSearch)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while adding to database, %s", err.Error())
	}

	defer response.Body.Close()
	if response.IsError() {
		// TODO: Add Logger here
		return fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
	}

	return nil
}

func (*postRepository) UpdatePost(id uuid.UUID, body []byte) error {
	// TODO: MAKE ASYNC
	body = []byte(fmt.Sprintf(`{"doc":%s}`, body))
	reader := bytes.NewReader(body)

	request := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id.String(),
		Body:       reader,
	}

	response, err := request.Do(context.Background(), data.ElasticSearch)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while updating to database, %s", err.Error())
	}

	defer response.Body.Close()
	if response.IsError() {
		// TODO: Add Logger here
		return fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
	}

	return nil
}

func (postRepository) RemovePost(id uuid.UUID) error {
	// TODO: MAKE ASYNC
	request := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id.String(),
	}

	response, err := request.Do(context.Background(), data.ElasticSearch)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while removing from database, %s", err.Error())
	}

	defer response.Body.Close()
	if response.IsError() {
		// TODO: Add Logger here
		return fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
	}

	return nil
}
