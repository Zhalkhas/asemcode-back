package it_bagalau

import (
	"context"
	"firebase.google.com/go/v4/db"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"regexp"
	"strconv"
	"time"
)

type ChatService struct {
	db            *db.Client
	chatGptClient *openai.Client
}

func NewChatService(db *db.Client, chatGptClient *openai.Client) *ChatService {
	return &ChatService{db: db, chatGptClient: chatGptClient}
}

func (c *ChatService) CreateChat(targetGrade, stack string, technologies []string) (*Chat, error) {
	technologiesString := ""
	for i, technology := range technologies {
		if i != 0 {
			technologiesString += ", "
		}
		technologiesString += technology
	}
	initialPrompt := fmt.Sprintf(
		"You are strict and demanding senior %s developer, "+
			"that is taking interview from %s %s developer. "+
			"Interview should cover %s. "+
			"Ask one question in one message. "+
			"When interview is done, give grade in format \"Grade: x/10\", "+
			"where x is grade out of 10, and give expanded comment on missed topics. "+
			"Give less than 4 grade if most of answers not given.",
		stack, targetGrade, stack, technologiesString,
	)
	chat := &Chat{
		TargetGrade:  targetGrade,
		Technologies: technologies,
		Stack:        stack,
		Messages: []Message{
			{
				Text:          initialPrompt,
				MessageAuthor: openai.ChatMessageRoleSystem,
			},
			{
				Text:          "Hello! Are you ready to start interview?",
				MessageAuthor: openai.ChatMessageRoleAssistant,
			},
		},
	}
	chatRef, err := c.db.NewRef("chats").Push(context.Background(), chat)
	if err != nil {
		log.Println("error creating chat", err)
		return nil, err
	}
	chat.ChatId = chatRef.Key
	return chat, nil
}

func (c *ChatService) GetChat(chatID string) (*Chat, error) {
	chat := &Chat{}
	chatRef := c.db.NewRef("chats").Child(chatID)
	err := chatRef.Get(context.Background(), chat)
	if err != nil {
		log.Println("error getting chat", err)
		return nil, err
	}
	return chat, nil
}

func (c *ChatService) CreateMessage(chatID, messageText string) (*Message, error) {
	messagesRef := c.db.NewRef("chats").Child(chatID).Child("messages")
	var messages []Message
	err := messagesRef.Get(context.Background(), &messages)
	if err != nil {
		log.Println("error getting messages", err)
		return nil, err
	}
	completionMessages := make([]openai.ChatCompletionMessage, len(messages)+1)
	for i, message := range messages {
		completionMessages[i] = openai.ChatCompletionMessage{
			Role:    message.MessageAuthor,
			Content: message.Text,
		}
	}
	completionMessages[len(messages)] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: messageText,
	}
	completion, err := c.chatGptClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		TopP:        1,
		MaxTokens:   256,
		N:           1,
		Temperature: 0.7,
		Messages:    completionMessages,
	})
	if err != nil {
		return nil, err
	}
	newMessage := &Message{
		Text:          completion.Choices[0].Message.Content,
		MessageAuthor: completion.Choices[0].Message.Role,
		CreatedAt:     time.Unix(completion.Created, 0),
	}
	gradeRegex := regexp.MustCompile("([0-9]{1,2})/10")
	gradeMatch := gradeRegex.FindStringSubmatch(newMessage.Text)
	if len(gradeMatch) > 1 {
		grade, err := strconv.ParseInt(gradeMatch[1], 10, 64)
		if err == nil {
			newMessage.GradingResult = &GradingResult{
				Grade: int(grade),
			}
		}
	}
	messages = append(messages, Message{
		Text:          messageText,
		MessageAuthor: openai.ChatMessageRoleUser,
		CreatedAt:     time.Now(),
	})
	messages = append(messages, *newMessage)
	err = messagesRef.Set(context.Background(), messages)
	if err != nil {
		return nil, err
	}
	return newMessage, nil
}
