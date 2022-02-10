package services

import (
	"bytes"
	"encoding/json"
	"es/entities"
	"es/inputs"
	"es/payloads"
	"es/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var PostService IPostService = &postService{}

type IPostService interface {
	FindPosts() ([]payloads.Post, error)
	FindPostByID(id uuid.UUID) (*payloads.Post, error)
	AddPost(input inputs.Post) error
	UpdatePost(id uuid.UUID, input inputs.Post) error
	RemovePost(id uuid.UUID) error
}

type postService struct{}

func (*postService) FindPosts() ([]payloads.Post, error) {
	payloads := []payloads.Post{}
	buffer := bytes.Buffer{}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	if err := json.NewEncoder(&buffer).Encode(query); err != nil {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error while encoding to buffer from json, %s", err.Error())
	}

	entities, err := repositories.PostRepository.FindPosts(buffer)
	if err != nil {
		// TODO: Add Logger here
		return nil, err
	}

	copier.Copy(&payloads, &entities)

	return payloads, nil
}

func (*postService) FindPostByID(id uuid.UUID) (*payloads.Post, error) {
	payload := payloads.Post{}

	entity, err := repositories.PostRepository.FindPostByID(id)
	if err != nil {
		// TODO: Add Logger here
		return nil, err
	}

	err = copier.Copy(&payload, &entity)
	if err != nil {
		// TODO: Add Logger here
		return nil, fmt.Errorf("error while copying to payload from entity, %s", err.Error())
	}

	return &payload, nil
}

func (*postService) AddPost(input inputs.Post) error {
	entity := &entities.Post{
		Id:      uuid.New(),
		Created: time.Now(),
	}

	err := copier.Copy(&entity, &input)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while copying to entity from input, %s", err.Error())
	}

	body, err := json.Marshal(entity)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while marshal post to json, %s", err.Error())
	}

	err = repositories.PostRepository.AddPost(entity.Id, body)
	if err != nil {
		// TODO: Add Logger here
		return err
	}

	return nil
}

func (service *postService) UpdatePost(id uuid.UUID, input inputs.Post) error {
	updated := time.Now()

	entity, err := service.FindPostByID(id)
	if err != nil {
		return err
	}

	if entity == nil {
		// TODO: Add Logger here
		return fmt.Errorf("record does not exist")
	}
	entity.Updated = &updated

	err = copier.Copy(&entity, &input)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while copying to entity from input, %s", err.Error())
	}

	body, err := json.Marshal(entity)
	if err != nil {
		// TODO: Add Logger here
		return fmt.Errorf("error while marshal post to json, %s", err.Error())
	}

	err = repositories.PostRepository.UpdatePost(id, body)
	if err != nil {
		// TODO: Add Logger here
		return err
	}

	return nil
}

func (*postService) RemovePost(id uuid.UUID) error {
	err := repositories.PostRepository.RemovePost(id)
	if err != nil {
		// TODO: Add Logger here
		return err
	}

	return nil
}
