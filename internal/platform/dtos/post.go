package dtos

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
