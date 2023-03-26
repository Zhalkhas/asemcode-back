package it_bagalau

import (
	"net/http"
)

type MessageResponse struct {
	Text          string `json:"text"`
	MessageAuthor string `json:"message_author"`
	Grade         *int   `json:"grade,omitempty"`
}

func (c MessageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ChatResponse struct {
	ChatId       string            `json:"chat_id"`
	TargetGrade  string            `json:"target_grade"`
	Technologies []string          `json:"technologies"`
	Messages     []MessageResponse `json:"messages"`
}

func (g ChatResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
