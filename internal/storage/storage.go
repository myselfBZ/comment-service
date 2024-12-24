package storage

import (
	"context"
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	BlogID    string    `json:"blog_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentStorage interface {
	Create(*Comment, context.Context) error
	GetByID(string, context.Context) (*Comment, error)
}
