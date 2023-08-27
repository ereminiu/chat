package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ereminiu/chat/entities"
	"github.com/ereminiu/chat/inputs"
	"github.com/ereminiu/chat/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Serivce
}

func NewHandler(s *service.Serivce) (*Handler, error) {
	return &Handler{s: s}, nil
}

func (h *Handler) AddUser(c *gin.Context) {
	var user entities.User
	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err)
	}
	if h.s.UserNameExists(user.Name) {
		h.NameIsUsed("user", user.Name, c)
		return
	}
	err := h.s.AddUser(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong...",
		})
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"id": user.Id,
	})
}

func (h *Handler) AddUserToChat(c *gin.Context) {
	var input inputs.UserToChat
	if err := c.BindJSON(&input); err != nil {
		log.Fatal(err)
	}
	userId, err := strconv.Atoi(input.UserId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong user_id format",
		})
		log.Fatal(err)
	}
	if !h.s.UserIdExists(userId) {
		h.IdNotFound("user", userId, c)
		return
	}
	chatId, err := strconv.Atoi(input.ChatId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong chat_id format",
		})
		log.Fatal(err)
	}
	if !h.s.ChatIdExists(chatId) {
		h.IdNotFound("chat", chatId, c)
		return
	}
	h.s.AddUserToChat(&entities.Chat{Id: chatId}, userId)
	c.JSON(200, gin.H{
		"message": "user has been added to chat",
	})
}

func (h *Handler) CreateChat(c *gin.Context) {
	var input inputs.Chat
	if err := c.BindJSON(&input); err != nil {
		log.Fatal(err)
	}
	if h.s.ChatNameExists(input.Name) {
		h.NameIsUsed("chat", input.Name, c)
		return
	}
	userIds := make([]int, 0, len(input.Users))
	for _, userId := range input.Users {
		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "found invalid id in users list",
				"id":      userId,
			})
			log.Fatal(err)
		}
		if !h.s.UserIdExists(id) {
			h.IdNotFound("user", id, c)
			return
		}
		userIds = append(userIds, id)
	}
	chat := entities.Chat{Name: input.Name}
	h.s.AddChat(&chat)
	for _, userId := range userIds {
		h.s.AddUserToChat(&chat, userId)
	}
	c.JSON(200, gin.H{
		"id": chat.Id,
	})
}

func (h *Handler) GetUsersChats(c *gin.Context) {
	var user inputs.User
	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err)
	}
	if !h.s.UserIdExists(user.Id) {
		h.IdNotFound("user", user.Id, c)
		return
	}
	res, err := h.s.GetUsersChats(user.Id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"chats": res,
	})
}

func (h *Handler) SendMessage(c *gin.Context) {
	var input inputs.Message
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request format",
		})
		log.Fatal(err)
	}
	message := entities.Message{Text: input.Text}
	if err := h.s.AddMessage(&message); err != nil {
		c.JSON(500, gin.H{
			"error":   "Can't add message",
			"message": message,
		})
		log.Fatal(err)
	}
	chatId, err := strconv.Atoi(input.Chat)
	if err != nil {
		log.Fatal(err)
	}
	if !h.s.ChatIdExists(chatId) {
		h.IdNotFound("chat", chatId, c)
		return
	}
	if err := h.s.AddMessageToChat(&message, chatId); err != nil {
		c.JSON(500, gin.H{
			"error":   "Can't add massage to chat",
			"message": message,
			"chat_id": chatId,
		})
		log.Fatal(err)
	}
	userId, err := strconv.Atoi(input.Author)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid input format",
		})
		log.Fatal(err)
	}
	if !h.s.UserIdExists(userId) {
		h.IdNotFound("author", userId, c)
		return
	}
	if err = h.s.AddMessageToUser(&message, userId); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"id": message.Id,
	})
}

func (h *Handler) GetMessagesByChat(c *gin.Context) {
	var input inputs.ChatId
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "wrong input format",
		})
		log.Fatal(err)
	}
	chatId, err := strconv.Atoi(input.Chat)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong input format for chat_id",
		})
		log.Fatal(err)
	}
	if !h.s.ChatIdExists(chatId) {
		h.IdNotFound("chat", chatId, c)
		return
	}
	res, err := h.s.GetMessagesByChat(chatId)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"messages": res,
	})
}

func (h *Handler) DebugGetMessagesByChat(c *gin.Context) {
	var input inputs.ChatId
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "wrong input format",
		})
		log.Fatal(err)
	}
	chatId, err := strconv.Atoi(input.Chat)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong input format",
		})
		log.Fatal(err)
	}
	if !h.s.ChatIdExists(chatId) {
		h.IdNotFound("chat", chatId, c)
		return
	}
	res, err := h.s.DebugGetMessagesByChat(chatId)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"messages": res,
	})
}

func (h *Handler) IdNotFound(role string, id int, c *gin.Context) {
	c.JSON(400, gin.H{
		"error": fmt.Sprintf("%s's id not found", role),
		"id":    id,
	})
}

func (h *Handler) NameNotFound(role, name string, c *gin.Context) {
	c.JSON(400, gin.H{
		"error": fmt.Sprintf("%s's name not found", role),
		"name":  name,
	})
}

func (h *Handler) NameIsUsed(role, name string, c *gin.Context) {
	c.JSON(400, gin.H{
		"error": fmt.Sprintf("%s's name is already used", role),
		"name":  name,
	})
}
