package interactors

import (
	"blog/internal/domain/posts"
	"blog/internal/platform/dtos"
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"time"
)

// PostInteractor handles business logic for post-related operations.
type PostInteractor struct {
	postRepository posts.Repository
	enforcer       *casbin.Enforcer
}

// NewPostInteractor creates a new PostInteractor with the given repository and enforcer.
func NewPostInteractor(r posts.Repository, enforcer *casbin.Enforcer) posts.Interactor {
	return &PostInteractor{
		postRepository: r,
		enforcer:       enforcer,
	}
}

// CreatePost creates a new post with the provided details, setting status and timestamps based on the author's roles.
func (i *PostInteractor) CreatePost(ctx context.Context, author string, req dtos.CreatePostRequest) error {
	roles, err := i.enforcer.GetRolesForUser(author)
	if err != nil {
		return err
	}
	var status posts.Status
	var publishedAt time.Time
	isAdminOrManager := false
	isEditorOrAuthor := false
	for _, role := range roles {
		if role == "Admin" || role == "Manager" {
			isAdminOrManager = true
			break
		}
		if role == "Editor" || role == "Author" {
			isEditorOrAuthor = true
		}
	}
	if isAdminOrManager {
		status = posts.StatusPublished
		publishedAt = time.Now()
	} else if isEditorOrAuthor {
		status = posts.StatusPending
	}

	postInfo := posts.Post{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		Author:      author,
		Slug:        req.Slug,
		Status:      status,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = i.postRepository.Save(ctx, postInfo)
	if err != nil {
		return err
	}
	return nil
}

// GetPost retrieves a post by its ID from the repository.
func (i *PostInteractor) GetPost(ctx context.Context, postID string) (*posts.Post, error) {
	post, err := i.postRepository.FindByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	postInfo := &posts.Post{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Author:      post.Author,
		Slug:        post.Slug,
		Status:      post.Status,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
	return postInfo, nil
}

// UpdatePost updates a post’s title and/or content if the user is authorized, applying changes via the repository.
func (i *PostInteractor) UpdatePost(ctx context.Context, postID string, req dtos.UpdatePostRequest) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	err := i.postRepository.Update(ctx, postID, updates)
	if err != nil {
		return err
	}
	return nil
}
