package main

import (
	"net/http"
	_ "ticket-api/docs"
	"ticket-api/internal/config"
	"ticket-api/internal/middleware"
	"ticket-api/internal/routes"
	"ticket-api/internal/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	g.MaxMultipartMemory = config.Get().App.MaxUploadFilesSize << 10

	// Validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phoneNumber", util.ValidatePhoneNumber)
	}

	cfgCORS := config.Get().CORSConfig
	g.Use(cors.New(cors.Config{
		AllowOrigins:     cfgCORS.AllowOrigins,
		AllowMethods:     cfgCORS.AllowMethods,
		AllowHeaders:     cfgCORS.AllowHeaders,
		ExposeHeaders:    cfgCORS.ExposeHeaders,
		AllowCredentials: cfgCORS.AllowCredentials,
		MaxAge:           time.Duration(cfgCORS.MaxAgeHours) * time.Hour,
	}))

	g.Use(middleware.RouteStatusChecker())

	v1 := g.Group("/api/v1")
	{
		captchaGroup := v1.Group("")
		captchaGroup.Use(middleware.CaptchaMiddleware(app.services.Token))
		captchaGroup.Use(middleware.RateLimitMiddleware(app.redis, 30))
		captchaGroup.Use(middleware.LimitRequestBody(config.Get().App.MaxJsonRequestSize))
		{
			captchaGroup.POST(routes.APIRoutes.Auth.SignUp.Path, app.handlers.Auth.SignUpWithPassword)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateTicket.Path, app.handlers.Ticket.CreateTicketHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.GetTicketByTrackCode.Path, app.handlers.Ticket.GetTicketByTrackCodeHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateChat.Path, app.handlers.Chat.CreateChatHandler)
			captchaGroup.POST(routes.APIRoutes.Auth.LoginWithNoAuth.Path, app.handlers.Auth.LoginWithNoAuth)
		}

		LoginGroup := v1.Group("")
		LoginGroup.Use(middleware.RateLimitMiddleware(app.redis, 10))
		LoginGroup.POST(routes.APIRoutes.Auth.Login.Path, app.handlers.Auth.LoginWithPassword)

		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthorizationMiddleware(app.services.Token))
		authGroup.Use(middleware.RateLimitMiddleware(app.redis, 15))
		authGroup.Use(middleware.LimitRequestBody(config.Get().App.MaxJsonRequestSize))
		{
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketsList.Path, app.handlers.Ticket.GetTicketsListHandler)
			authGroup.POST(routes.APIRoutes.Users.GetUsersByIDs.Path, app.handlers.User.GetUsersByIDs)
			authGroup.POST(routes.APIRoutes.Users.GetUserByID.Path, app.handlers.User.GetUserByID)
			authGroup.POST(routes.APIRoutes.Users.GetUserByUsername.Path, app.handlers.User.GetUserByUsername)
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketByID.Path, app.handlers.Ticket.GetTicketByIDHandler)
		}

		publicGroup := v1.Group("")
		publicGroup.Use(middleware.RateLimitMiddleware(app.redis, 30))
		publicGroup.Use(middleware.LimitRequestBody(config.Get().App.MaxJsonRequestSize))
		{
			publicGroup.GET(routes.APIRoutes.Versions.GetCurrentVersion.Path, app.handlers.Version.GetCurrentVersionHandler)
			publicGroup.GET(routes.APIRoutes.Captcha.GetCaptcha.Path, app.handlers.Captcha.GenerateCaptchaHandler)
			publicGroup.POST(routes.APIRoutes.Captcha.VerifyCaptcha.Path, app.handlers.Captcha.VerifyCaptchaHandler)

			publicGroup.GET(routes.APIRoutes.Auth.LoginWithSingleUseToken.Path, app.handlers.Auth.LoginWithOneTimeToken)

			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketTypes.Path, app.handlers.Ticket.GetAllActiveTicketTypesHandler)
			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketStatuses.Path, app.handlers.Ticket.GetAllActiveTicketStatusesHandler)
			publicGroup.GET(routes.APIRoutes.Departments.GetAllActiveDepartments.Path, app.handlers.Department.GetAllActiveDepartmentsHandler)

		}

		fileGroup := v1.Group("")
		fileGroup.Use(middleware.RateLimitMiddleware(app.redis, 10))
		{
			fileGroup.POST(routes.APIRoutes.Files.UploadTicketFile.Path, app.handlers.File.UploadTicketFileHandler)
			fileGroup.POST(routes.APIRoutes.Files.GetDownloadLinkTicketFile.Path, app.handlers.File.GetDownloadLinkTicketFileHandler)
		}

		_APIKeyGroup := v1.Group("")
		_APIKeyGroup.Use(middleware.LimitRequestBody(config.Get().App.MaxJsonRequestSize))
		_APIKeyGroup.Use(middleware.ApiKeyGuardMiddleware(app.services.Token, app.repos.APIKeys))
		{
			_APIKeyGroup.POST(routes.APIRoutes.Auth.GetSingleUseToken.Path, app.handlers.Auth.GetSingleUseToken)
		}
	}

	// Redirect /swagger and /swagger/ → /swagger/index.html
	g.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Serve Swagger UI for all other paths
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return g
}
