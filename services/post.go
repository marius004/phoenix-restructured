package services

import (
	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
)

type PostService struct {
	postRepository internal.PostRepository
}

func (s *PostService) GetPosts() ([]*entities.Post, error) {
	return s.postRepository.GetPosts()
}

func (s *PostService) GetPostByTitle(title string) (*entities.Post, error) {
	return s.postRepository.GetPostByTitle(title)
}

func (s *PostService) GetPostByID(id uint) (*entities.Post, error) {
	return s.postRepository.GetPostByID(id)
}

func (s *PostService) CreatePost(post *entities.Post) error {
	return s.postRepository.CreatePost(post)
}

func (s *PostService) UpdatePostByID(id uint, request *models.UpdatePostRequest) error {
	return s.postRepository.UpdatePostByID(id, request)
}

func (s *PostService) UpdatePostByTitle(title string, request *models.UpdatePostRequest) error {
	return s.postRepository.UpdatePostByTitle(title, request)
}

func (s *PostService) DeletePost(post *entities.Post) error {
	return s.postRepository.DeletePost(post)
}
