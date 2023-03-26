package it_bagalau

import "time"

type Chat struct {
	ChatId       string    `json:"chat_id"`
	TargetGrade  string    `json:"target_grade"`
	Technologies []string  `json:"technologies"`
	Stack        string    `json:"stack"`
	Messages     []Message `json:"messages"`
}

type GradingResult struct {
	// Grade, out of 10
	Grade int `json:"grade"`
}

type Message struct {
	Text          string         `json:"text"`
	MessageAuthor string         `json:"message_author"`
	CreatedAt     time.Time      `json:"created_at"`
	GradingResult *GradingResult `json:"grading_result"`
}
