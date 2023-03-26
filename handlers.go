package it_bagalau

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type ChatController struct {
	*chi.Mux
	ChatService *ChatService
}

func NewChatController(service *ChatService) *ChatController {
	controller := &ChatController{
		Mux:         chi.NewRouter(),
		ChatService: service,
	}
	controller.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	controller.Use(render.SetContentType(render.ContentTypeJSON))
	controller.Use(middleware.Recoverer)
	controller.Use(middleware.Logger)
	controller.Post("/", controller.createChat)
	controller.Get("/{chatId}", controller.getChat)
	controller.Post("/{chatId}/message", controller.createMessage)
	return controller

}

func (c *ChatController) getChat(rw http.ResponseWriter, req *http.Request) {
	chatID := chi.URLParam(req, "chatId")
	chat, err := c.ChatService.GetChat(chatID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chat == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	response := ChatResponse{
		ChatId:       chat.ChatId,
		TargetGrade:  chat.TargetGrade,
		Technologies: chat.Technologies,
		Messages:     make([]MessageResponse, len(chat.Messages)),
	}
	for i, message := range chat.Messages {
		response.Messages[i] = MessageResponse{
			Text:          message.Text,
			MessageAuthor: message.MessageAuthor,
		}
		grade := message.GradingResult
		if grade != nil {
			response.Messages[i].Grade = &grade.Grade
		}
	}
	err = render.Render(rw, req, response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
		return
	}
}
func (c *ChatController) createChat(rw http.ResponseWriter, req *http.Request) {
	var request CreateChatRequest
	if err := render.Bind(req, &request); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := c.ChatService.CreateChat(request.TargetGrade, request.Stack, request.Technologies)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageResponses := make([]MessageResponse, len(chat.Messages))
	for i, message := range chat.Messages {
		messageResponses[i] = MessageResponse{
			Text:          message.Text,
			MessageAuthor: message.MessageAuthor,
		}
		grade := message.GradingResult
		if grade != nil {
			messageResponses[i].Grade = &grade.Grade
		}
	}
	response := ChatResponse{
		ChatId:       chat.ChatId,
		TargetGrade:  chat.TargetGrade,
		Technologies: chat.Technologies,
		Messages:     messageResponses,
	}
	err = render.Render(rw, req, response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
		return
	}
}
func (c *ChatController) createMessage(rw http.ResponseWriter, req *http.Request) {
	chatID := chi.URLParam(req, "chatId")
	var request CreateMessageRequest
	if err := render.Bind(req, &request); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	message, err := c.ChatService.CreateMessage(chatID, request.Content)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := MessageResponse{
		Text:          message.Text,
		MessageAuthor: message.MessageAuthor,
	}
	grade := message.GradingResult
	if grade != nil {
		response.Grade = &grade.Grade
	}
	err = render.Render(rw, req, response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
