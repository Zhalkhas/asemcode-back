package it_bagalau

import (
	"net/http"
)

type CreateChatRequest struct {
	TargetGrade  string   `json:"target_grade"`
	Technologies []string `json:"technologies"`
	Stack        string   `json:"stack"`
}

func (c CreateChatRequest) Bind(r *http.Request) error {
	return nil
}

type CreateMessageRequest struct {
	Content string `json:"content"`
}

func (c CreateMessageRequest) Bind(r *http.Request) error {
	return nil
}
