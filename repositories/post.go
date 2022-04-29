package repositories

import (
	"errors"

	"github.com/marius004/phoenix/entities"
	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *internal.Database
}

func (r *PostRepository) GetPosts() ([]*entities.Post, error) {
	var posts []*entities.Post
	result := r.db.Conn.Find(&posts)

	return posts, result.Error
}

func (r *PostRepository) GetPostByTitle(title string) (*entities.Post, error) {
	var post *entities.Post
	result := r.db.Conn.Where("title = ?", title).First(&post)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return post, result.Error
}

func (r *PostRepository) GetPostByID(id uint) (*entities.Post, error) {
	var post *entities.Post
	result := r.db.Conn.Where("id = ?", id).First(&post)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return post, result.Error
}

func (r *PostRepository) CreatePost(post *entities.Post) error {
	result := r.db.Conn.Create(&post)
	return result.Error
}

func (r *PostRepository) UpdatePostByID(id uint, request *models.UpdatePostRequest) error {
	post, err := r.GetPostByID(id)

	if err != nil {
		return err
	} else if post == nil {
		return internal.ErrPostDoesNotExist
	}

	return r.updatePost(post, request)
}

func (r *PostRepository) UpdatePostByTitle(title string, request *models.UpdatePostRequest) error {
	post, err := r.GetPostByTitle(title)

	if err != nil {
		return err
	} else if post == nil {
		return internal.ErrPostDoesNotExist
	}

	return r.updatePost(post, request)
}

func (r *PostRepository) DeletePost(post *entities.Post) error {
	result := r.db.Conn.Unscoped().Delete(&entities.Post{}, "id = ?", post.ID)
	return result.Error
}

func (r *PostRepository) updatePost(post *entities.Post, request *models.UpdatePostRequest) error {
	if len(request.Content) > 0 {
		post.Content = []byte(request.Content)
	}

	if request.Title != "" {
		post.Title = request.Title
	}

	result := r.db.Conn.Save(&post)
	return result.Error
}

func NewPostRepository(db *internal.Database) *PostRepository {
	return &PostRepository{
		db: db,
	}
}
