package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/myselfBZ/comment-service/internal/storage"
)

type createCommentPayload struct {
	Content string `json:"content"`
	BlogID  string `json:"blog_id"`
	UserID  string `json:"user_id"`
}

func (a *app) createComment(w http.ResponseWriter, r *http.Request) {
	var comment createCommentPayload
	if err := readJSON(w, r, &comment); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	storeComment := &storage.Comment{
		ID:        uuid.NewString(),
		Content:   comment.Content,
		UserID:    comment.UserID,
		BlogID:    comment.BlogID,
		CreatedAt: time.Now(),
	}

	if err := a.store.Create(storeComment, r.Context()); err != nil {
		a.internalServerError(w, r, err)
		return
	}

	if err := a.kafkaProducer.Push(storeComment); err != nil {
		log.Println("error pushing a comment")
	}

	w.WriteHeader(http.StatusOK)
}

func (a *app) getCommentByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	comment, err := a.store.GetByID(id, r.Context())
	if err != nil {
		if err == storage.NotFound {
			a.notFoundResponse(w, r, err)
		} else {
			a.internalServerError(w, r, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, &comment)
}
