package main

import (
	"net/http"
	_ "ticket-api/docs"
	"ticket-api/internal/middleware"
	"ticket-api/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	// Validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phoneNumber", util.ValidatePhoneNumber)
	}

	v1 := g.Group("/api/v1")
	{
		// Routes protected by Captcha middleware
		captchaGroup := v1.Group("")
		captchaGroup.Use(middleware.CaptchaMiddleware(app.services.Token))
		{
			captchaGroup.POST("auth/sigup", app.handlers.Auth.SigupWithPassword)
			captchaGroup.GET("auth/login", app.handlers.Auth.LoginWithPassword)
			captchaGroup.POST("tickets", app.handlers.Ticket.CreateTicketHandler)
			captchaGroup.POST("ticket/track-code", app.handlers.Ticket.GetTicketByTrackCodeHandler)
			captchaGroup.POST("tickets/:id/chat", app.handlers.Chat.CreateChatHandler)
		}

		// Routes protected by Authorization middleware
		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthorizationMiddleware(app.services.Token))
		{
			authGroup.POST("ticket/id", app.handlers.Ticket.GetTicketByIDHandler)
			authGroup.POST("auth/LoginWithNoAuth", app.handlers.User.LoginWithNoAuth)
		}

		// Auth routes (no middleware for one-time token)
		publicGroup := v1.Group("")
		{
			// Version and public captcha routes (no middleware)
			publicGroup.GET("", app.handlers.Version.GetCurrentVersionHandler)
			publicGroup.GET("captcha/new", app.handlers.Captcha.GenerateCaptchaHandler)
			publicGroup.POST("captcha/verify", app.handlers.Captcha.VerifyCaptchaHandler)

			publicGroup.POST("auth/one-time-token", app.handlers.Auth.GenerateOneTimeToken)
			publicGroup.GET("auth/login/token", app.handlers.Auth.LoginWithOneTimeToken)

			publicGroup.POST("tickets/list", app.handlers.Ticket.ListTicketsHandler)
		}
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
