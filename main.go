package main

import (
	"log"

	"github.com/ereminiu/chat/pkg/handlers"
	"github.com/ereminiu/chat/pkg/repository"
	"github.com/ereminiu/chat/pkg/service"
	"github.com/gin-gonic/gin"
)

func main() {
	repos, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	service, err := service.NewService(repos)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := handlers.NewHandler(service)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("/users/add", handler.AddUser)
	r.POST("users/add/tochat", handler.AddUserToChat)
	r.POST("/chats/add", handler.CreateChat)
	r.GET("/chats/get", handler.GetUsersChats)
	r.POST("messages/add", handler.SendMessage)
	r.GET("/messages/get", handler.GetMessagesByChat)
	r.GET("/messages/get/debug", handler.DebugGetMessagesByChat)

	r.Run(":9000")
}
