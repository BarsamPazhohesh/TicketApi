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
		_ = v.RegisterValidation("phoneNumber", util.ValidatePhoneNumber)
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

	cfgApp := config.Get().App
	cfgRL := config.Get().RateLimit

	v1 := g.Group("/api/v1")
	v1.Use(middleware.LimitRequestBody(cfgApp.MaxJsonRequestSize))
	{
		// 1. Public Tier (Open lookups, Captcha, & Login)
		publicGroup := v1.Group("")
		publicGroup.Use(middleware.RateLimitMiddleware(app.redis, cfgRL.Public))
		{
			publicGroup.GET(routes.APIRoutes.Captcha.GetCaptcha.Path, app.handlers.Captcha.GenerateCaptchaHandler)
			publicGroup.POST(routes.APIRoutes.Captcha.VerifyCaptcha.Path, app.handlers.Captcha.VerifyCaptchaHandler)
			publicGroup.POST(routes.APIRoutes.Auth.Login.Path, app.handlers.Auth.LoginWithPassword)
			publicGroup.GET(routes.APIRoutes.Auth.LoginWithSingleUseToken.Path, app.handlers.Auth.LoginWithOneTimeToken)
			publicGroup.GET(routes.APIRoutes.Auth.CheckToken.Path, app.handlers.Auth.CheckToken)
			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketTypes.Path, app.handlers.Ticket.GetAllActiveTicketTypesHandler)
			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketStatuses.Path, app.handlers.Ticket.GetAllActiveTicketStatusesHandler)
			publicGroup.GET(routes.APIRoutes.Departments.GetAllActiveDepartments.Path, app.handlers.Department.GetAllActiveDepartmentsHandler)
		}

		// 2. Step-Up Level 1: OTP Tier (Human / Captcha Verified)
		otpGroup := v1.Group("")
		otpGroup.Use(middleware.CaptchaMiddleware(app.services.Token, false))
		otpGroup.Use(middleware.RateLimitMiddleware(app.redis, cfgRL.OTP))
		{
			otpGroup.POST(routes.APIRoutes.OTP.SendOTP.Path, app.handlers.OTP.SendOTP)
			otpGroup.POST(routes.APIRoutes.OTP.VerifyOTP.Path, app.handlers.OTP.VerifyOTP)
		}

		// 3. Step-Up Level 2: Customer & Guest Tier (Phone Verified OR Logged-in Auth Token)
		customerGroup := v1.Group("")
		customerGroup.Use(middleware.CaptchaMiddleware(app.services.Token, true))
		customerGroup.Use(middleware.RateLimitMiddleware(app.redis, cfgRL.Customer))
		{
			customerGroup.POST(routes.APIRoutes.Tickets.CreateTicket.Path, app.handlers.Ticket.CreateTicketHandler)
			customerGroup.POST(routes.APIRoutes.Tickets.GetTicketByTrackCode.Path, app.handlers.Ticket.GetTicketByTrackCodeHandler)
			customerGroup.POST(routes.APIRoutes.Tickets.CreateChat.Path, app.handlers.Chat.CreateChatHandler)
			customerGroup.POST(routes.APIRoutes.Files.UploadTicketFile.Path, app.handlers.File.UploadTicketFileHandler)
			customerGroup.POST(routes.APIRoutes.Files.GetDownloadLinkTicketFile.Path, app.handlers.File.GetDownloadLinkTicketFileHandler)
		}

		// 4. Staff & Admin RBAC Tier (Auth Token + Dynamic RBAC)
		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthorizationMiddleware(app.services.Token))
		authGroup.Use(middleware.DynamicRBACGuardMiddleware(app.security))
		authGroup.Use(middleware.RateLimitMiddleware(app.redis, cfgRL.Auth))
		{
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketsList.Path, app.handlers.Ticket.GetTicketsListHandler)
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketByID.Path, app.handlers.Ticket.GetTicketByIDHandler)
			authGroup.POST(routes.APIRoutes.Users.GetUserByID.Path, app.handlers.User.GetUserByID)
			authGroup.POST(routes.APIRoutes.Users.GetUserByUsername.Path, app.handlers.User.GetUserByUsername)
			authGroup.POST(routes.APIRoutes.Users.GetUsersByIDs.Path, app.handlers.User.GetUsersByIDs)
			authGroup.POST(routes.APIRoutes.Auth.SignUp.Path, app.handlers.Auth.SignUpWithPassword)
			authGroup.POST(routes.APIRoutes.Auth.LoginWithNoAuth.Path, app.handlers.Auth.LoginWithNoAuth)
		}

		// 5. Machine-to-Machine Tier (API Key Guard)
		apiKeyGroup := v1.Group("")
		apiKeyGroup.Use(middleware.ApiKeyGuardMiddleware(app.services.Token, app.repos.APIKeys))
		{
			apiKeyGroup.POST(routes.APIRoutes.Auth.GetSingleUseToken.Path, app.handlers.Auth.GetSingleUseToken)
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
