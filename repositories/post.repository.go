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
	FindPosts(buffer bytes.Buffer) ([]entities.Post, error)
	FindPostByID(id uuid.UUID) (*entities.Post, error)
	AddPost(id uuid.UUID, body []byte) error
	UpdatePost(id uuid.UUID, body []byte) error
	RemovePost(id uuid.UUID) error
}

type postRepository struct{}

func (*postRepository) FindPosts(buffer bytes.Buffer) ([]entities.Post, error) {
	response, err := data.ElasticSearch.Search(
		data.ElasticSearch.Search.WithContext(context.Background()),
		data.ElasticSearch.Search.WithIndex(index),
		data.ElasticSearch.Search.WithBody(&buffer),
	)

	defer response.Body.Close()
	if err != nil {
		message := fmt.Errorf("error while fetching records from database, %s", err.Error())
		extensions.Error(message.Error())
		return nil, message
	}

	if response.IsError() {
		message := fmt.Errorf("error indexing documents status: %s", response.Status())
		extensions.Error(message.Error())
		return nil, message
	}

	body := make(map[string]interface{})
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		message := fmt.Errorf("error parsing the response body: %s", err.Error())
		extensions.Error(message.Error())
		return nil, message
	}

	posts := []entities.Post{}
	for _, hit := range body["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var post entities.Post
		err = extensions.Decode(hit.(map[string]interface{})["_source"].(map[string]interface{}), &post)

		if err != nil {
			message := fmt.Errorf("error mapping from _source to entity: %s", err.Error())
			extensions.Error(message.Error())
			return nil, err
		}

		posts = append(posts, post)
	}

	extensions.Info("done")
	return posts, nil
}

func (*postRepository) FindPostByID(id uuid.UUID) (*entities.Post, error) {
	channel := make(chan struct {
		*entities.Post
		error
	})
	go func() {
		var entity entities.Post

		request := esapi.GetRequest{
			Index:      index,
			DocumentID: id.String(),
		}

		response, err := request.Do(context.Background(), data.ElasticSearch)
		if err != nil {
			message := fmt.Errorf("error while fetching record from database, %s", err.Error())
			extensions.Error(message.Error())
			channel <- struct {
				*entities.Post
				error
			}{nil, message}
			return
		}

		defer response.Body.Close()
		if response.IsError() {
			message := fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
			extensions.Error(message.Error())
			channel <- struct {
				*entities.Post
				error
			}{nil, message}
			return
		}

		body := make(map[string]interface{})
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			message := fmt.Errorf("error parsing the response body: %s", err.Error())
			extensions.Error(message.Error())
			channel <- struct {
				*entities.Post
				error
			}{nil, message}
			return
		}

		err = extensions.Decode(body["_source"].(map[string]interface{}), &entity)
		if err != nil {
			message := fmt.Errorf("error mapping from _source to entity: %s", err.Error())
			extensions.Error(message.Error())
			channel <- struct {
				*entities.Post
				error
			}{nil, message}
			return
		}

		channel <- struct {
			*entities.Post
			error
		}{&entity, nil}
	}()

	extensions.Info("done")
	response := <-channel
	return response.Post, response.error
}

func (*postRepository) AddPost(id uuid.UUID, body []byte) error {
	channel := make(chan error)
	go func() {
		reader := strings.NewReader(string(body))

		request := esapi.IndexRequest{
			Index:      index,
			DocumentID: id.String(),
			Body:       reader,
		}

		response, err := request.Do(context.Background(), data.ElasticSearch)
		if err != nil {
			message := fmt.Errorf("error while adding to database, %s", err.Error())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		defer response.Body.Close()
		if response.IsError() {
			message := fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		extensions.Info("done")
		channel <- nil
	}()

	return <-channel
}

func (*postRepository) UpdatePost(id uuid.UUID, body []byte) error {
	channel := make(chan error)
	go func() {
		body = []byte(fmt.Sprintf(`{"doc":%s}`, body))
		reader := bytes.NewReader(body)

		request := esapi.UpdateRequest{
			Index:      index,
			DocumentID: id.String(),
			Body:       reader,
		}

		response, err := request.Do(context.Background(), data.ElasticSearch)
		if err != nil {
			message := fmt.Errorf("error while updating to database, %s", err.Error())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		defer response.Body.Close()
		if response.IsError() {
			message := fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		extensions.Info("done")
		channel <- nil
	}()

	return <-channel
}

func (postRepository) RemovePost(id uuid.UUID) error {
	channel := make(chan error)
	go func() {
		request := esapi.DeleteRequest{
			Index:      index,
			DocumentID: id.String(),
		}

		response, err := request.Do(context.Background(), data.ElasticSearch)
		if err != nil {
			message := fmt.Errorf("error while removing from database, %s", err.Error())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		defer response.Body.Close()
		if response.IsError() {
			message := fmt.Errorf("error indexing document with id: %s, status: %s", id, response.Status())
			extensions.Info(message.Error())
			channel <- message
			return
		}

		extensions.Info("done")
		channel <- nil
	}()

	return <-channel
}
