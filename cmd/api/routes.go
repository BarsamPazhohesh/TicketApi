package main

import (
	"net/http"
	_ "ticket-api/docs"
	"ticket-api/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phoneNumber", util.ValidatePhoneNumber)
	}

	v1 := g.Group("/api/v1/")
	{
		// Version route
		v1.GET("", app.handlers.Version.GetCurrentVersionHandler)

		// Ticket routes
		v1.POST("tickets", app.handlers.Ticket.CreateTicketHandler) // create new ticket
		v1.GET("tickets/:id", app.handlers.Ticket.GetTicketHandler) // get ticket by ID

		v1.POST("tickets/:id/chat", app.handlers.Chat.CreateChatHandler) // create chat

		v1.POST("auth/LoginWithNoAuth", app.handlers.User.LoginWithNoAuth)
	}

	// Redirect /swagger → /swagger/index.html
	g.GET("/swagger", redirectSwagger)

	// Serve Swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return g
}

func redirectSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
