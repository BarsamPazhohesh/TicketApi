package main

import (
	"net/http"
	_ "ticket-api/docs"
	"ticket-api/internal/middleware"
	"ticket-api/internal/routes"
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
			captchaGroup.POST(routes.APIRoutes.Auth.SignUp.Path, app.handlers.Auth.SigupWithPassword)
			captchaGroup.GET(routes.APIRoutes.Auth.Login.Path, app.handlers.Auth.LoginWithPassword)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateTicket.Path, app.handlers.Ticket.CreateTicketHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.GetTicketByTrackCode.Path, app.handlers.Ticket.GetTicketByTrackCodeHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateChat.Path, app.handlers.Chat.CreateChatHandler)
		}

		// Routes protected by Authorization middleware
		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthorizationMiddleware(app.services.Token))
		{
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketByID.Path, app.handlers.Ticket.GetTicketByIDHandler)
			authGroup.POST(routes.APIRoutes.Auth.LoginWithNoAuth.Path, app.handlers.User.LoginWithNoAuth)
		}

		// Auth routes (no middleware for one-time token)
		publicGroup := v1.Group("")
		{
			// Version and public captcha routes (no middleware)
			publicGroup.GET(routes.APIRoutes.Versions.GetCurrentVersion.Path, app.handlers.Version.GetCurrentVersionHandler)
			publicGroup.GET(routes.APIRoutes.Captcha.GetCaptcha.Path, app.handlers.Captcha.GenerateCaptchaHandler)
			publicGroup.POST(routes.APIRoutes.Captcha.VerifyCaptcha.Path, app.handlers.Captcha.VerifyCaptchaHandler)

			publicGroup.POST(routes.APIRoutes.Auth.GetSingleUseToken.Path, app.handlers.Auth.GenerateOneTimeToken)
			publicGroup.GET(routes.APIRoutes.Auth.LoginWithSingleUseToken.Path, app.handlers.Auth.LoginWithOneTimeToken)

			publicGroup.POST(routes.APIRoutes.Tickets.GetTicketsList.Path, app.handlers.Ticket.ListTicketsHandler)
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
