package it_bagalau

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	"os"

	"log"
)

var (
	openaiApiKey = os.Getenv("OPENAI_API_KEY")
)

func init() {
	b := []byte(os.Getenv("FIREBASE_CONFIG"))
	opt := option.WithCredentialsJSON(b)

	config := &firebase.Config{DatabaseURL: "https://it-bagalau-default-rtdb.asia-southeast1.firebasedatabase.app/"}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Panicln(err)
	}
	firebaseDB, err := firebaseApp.Database(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	chatGptClient := openai.NewClient(openaiApiKey)
	chatService := NewChatService(firebaseDB, chatGptClient)
	controller := NewChatController(chatService)
	functions.HTTP("chat", controller.ServeHTTP)
}
