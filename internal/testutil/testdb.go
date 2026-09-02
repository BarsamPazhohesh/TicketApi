package testutil

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"ticket-api/internal/config"
	"ticket-api/internal/db/api_keys"
	"ticket-api/internal/db/api_routes"
	"ticket-api/internal/db/departments"
	"ticket-api/internal/db/roles"
	"ticket-api/internal/db/roles_relations"
	"ticket-api/internal/db/ticket_priorities"
	"ticket-api/internal/db/ticket_statuses"
	"ticket-api/internal/db/ticket_types"
	"ticket-api/internal/db/users"
	"ticket-api/internal/errx"
	"ticket-api/internal/handler"
	"ticket-api/internal/middleware"
	"ticket-api/internal/repository"
	"ticket-api/internal/routes"
	"ticket-api/internal/security"
	"ticket-api/internal/services"
	"ticket-api/internal/services/cache"
	"ticket-api/internal/services/token"
	"ticket-api/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "this-is-a-very-secure-32-byte-long-secret-key-12345")
}

// FindProjectRoot locates the root folder containing config.yaml
func FindProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// SetupTestSQLDB initializes an in-memory SQLite DB with complete project migrations executed
func SetupTestSQLDB(t *testing.T) *sql.DB {
	t.Helper()
	root := FindProjectRoot()
	_ = config.Load(filepath.Join(root, "config.yaml"))

	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}

	errx.NewRegistry(db)

	migrationsDir := filepath.Join(root, "cmd", "migrate", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations dir: %v", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			t.Fatalf("failed to read migration file %s: %v", file, err)
		}
		// Split queries if multiple in one file
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			trimmed := strings.TrimSpace(q)
			if trimmed == "" {
				continue
			}
			if _, err := db.Exec(trimmed); err != nil {
				t.Fatalf("failed executing migration %s: %v\nquery: %s", file, err, trimmed)
			}
		}
	}

	SeedBaseData(t, db)
	return db
}

// SeedBaseData inserts standard roles, departments, ticket types, statuses, and routes
func SeedBaseData(t *testing.T, db *sql.DB) {
	t.Helper()

	seeds := []string{
		`INSERT OR IGNORE INTO roles (id, title, status, created_at, updated_at) VALUES
			(1, 'superadmin', 1, datetime('now'), datetime('now')),
			(2, 'agent', 1, datetime('now'), datetime('now')),
			(3, 'customer', 1, datetime('now'), datetime('now'));`,

		`INSERT OR IGNORE INTO departments (id, title, description, status, created_at, updated_at) VALUES
			(1, 'Technical Support', 'Tech dept', 1, datetime('now'), datetime('now')),
			(2, 'Sales & Billing', 'Billing dept', 1, datetime('now'), datetime('now')),
			(3, 'Inactive Dept', 'Disabled dept', 0, datetime('now'), datetime('now'));`,

		`INSERT OR IGNORE INTO ticket_types (id, title, description, status, created_at, updated_at) VALUES
			(1, 'Bug Report', 'Software bugs', 1, datetime('now'), datetime('now')),
			(2, 'Feature Request', 'New feature', 1, datetime('now'), datetime('now')),
			(3, 'Inactive Type', 'Disabled type', 0, datetime('now'), datetime('now'));`,

		`INSERT OR IGNORE INTO ticket_statuses (id, title, description, status, created_at, updated_at) VALUES
			(1, 'open', 'Open ticket', 1, datetime('now'), datetime('now')),
			(2, 'in_progress', 'Ticket in progress', 1, datetime('now'), datetime('now')),
			(3, 'closed', 'Closed ticket', 1, datetime('now'), datetime('now'));`,

		// Seed initial routes into api_routes table so RBAC can bind
		`INSERT OR IGNORE INTO api_routes (id, route, method, description, status, created_at, updated_at) VALUES
			(1, '/tickets/GetTicketsList/', 'POST', 'Get tickets list', 1, datetime('now'), datetime('now')),
			(2, '/users/GetUsersByIDs/', 'POST', 'Get users list', 1, datetime('now'), datetime('now')),
			(3, '/users/GetUserByID/', 'POST', 'Get user by id', 1, datetime('now'), datetime('now')),
			(4, '/users/GetUserByUsername/', 'POST', 'Get user by username', 1, datetime('now'), datetime('now')),
			(5, '/tickets/GetTicketByID/', 'POST', 'Get ticket by ID', 1, datetime('now'), datetime('now'));`,

		// Allow superadmin (role 1) and agent (role 2) on routes
		`INSERT OR IGNORE INTO api_routes_roles_relation (api_route_id, role_id, status, created_at, updated_at) VALUES
			(1, 1, 1, datetime('now'), datetime('now')), (1, 2, 1, datetime('now'), datetime('now')),
			(2, 1, 1, datetime('now'), datetime('now')), (2, 2, 1, datetime('now'), datetime('now')),
			(3, 1, 1, datetime('now'), datetime('now')), (3, 2, 1, datetime('now'), datetime('now')),
			(4, 1, 1, datetime('now'), datetime('now')), (4, 2, 1, datetime('now'), datetime('now')),
			(5, 1, 1, datetime('now'), datetime('now')), (5, 2, 1, datetime('now'), datetime('now'));`,
	}

	for _, s := range seeds {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("failed seeding data: %v\nquery: %s", err, s)
		}
	}
}

// TestAppBundle bundles all initialized components for testing
type TestAppBundle struct {
	DB               *sql.DB
	Repos            *repository.AppRepositories
	Services         *services.AppServices
	Handlers         *handler.AppHandlers
	SecurityRegistry *security.SecurityRegistry
	Engine           *gin.Engine
}

// SetupTestApp constructs full application layers with in-memory SQLite and optional Mongo
func SetupTestApp(t *testing.T, mongoDB *mongo.Database) *TestAppBundle {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db := SetupTestSQLDB(t)
	cacheSvc := cache.NewCacheService(nil)

	userRepo := repository.NewUsersRepository(users.New(db))
	deptRepo := repository.NewDepartmentsRepository(departments.New(db), cacheSvc)
	typeRepo := repository.NewTicketTypesRepository(ticket_types.New(db), cacheSvc)
	statusRepo := repository.NewTicketStatusesRepository(ticket_statuses.New(db), cacheSvc)
	priorityRepo := repository.NewTicketPrioritiesRepository(ticket_priorities.New(db))
	rolesRepo := repository.NewRolesRepository(roles.New(db))
	apiKeyRepo := repository.NewAPIKeysRepository(api_keys.New(db))
	apiRouteRepo := repository.NewAPIRoutesRepository(api_routes.New(db))
	rolesRelRepo := repository.NewRolesRelationRepository(roles_relations.New(db), api_keys.New(db), api_routes.New(db))

	ticketRepo := repository.NewTicketRepository(mongoDB)
	chatRepo := repository.NewChatRepository(mongoDB)

	repos := &repository.AppRepositories{
		Ticket:           ticketRepo,
		ChatRepository:   chatRepo,
		Roles:            rolesRepo,
		Departments:      deptRepo,
		TicketTypes:      typeRepo,
		TicketPriorities: priorityRepo,
		APIRoutes:        apiRouteRepo,
		APIKeys:          apiKeyRepo,
		RolesRelations:   rolesRelRepo,
		Users:            userRepo,
		TicketStatus:     statusRepo,
	}

	appServices := services.NewAppService(nil, nil, repos)
	appHandlers := handler.NewAppHandlers(repos, appServices)

	secRegistry := security.NewSecurityRegistry(rolesRelRepo)
	_ = secRegistry.Reload(context.Background())

	// Build Gin engine
	g := gin.New()
	g.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phoneNumber", util.ValidatePhoneNumber)
	}

	v1 := g.Group("/api/v1")
	{
		captchaGroup := v1.Group("")
		captchaGroup.Use(middleware.CaptchaMiddleware(appServices.Token))
		{
			captchaGroup.POST(routes.APIRoutes.Auth.SignUp.Path, appHandlers.Auth.SignUpWithPassword)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateTicket.Path, appHandlers.Ticket.CreateTicketHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.GetTicketByTrackCode.Path, appHandlers.Ticket.GetTicketByTrackCodeHandler)
			captchaGroup.POST(routes.APIRoutes.Tickets.CreateChat.Path, appHandlers.Chat.CreateChatHandler)
			captchaGroup.POST(routes.APIRoutes.Auth.LoginWithNoAuth.Path, appHandlers.Auth.LoginWithNoAuth)
		}

		loginGroup := v1.Group("")
		loginGroup.POST(routes.APIRoutes.Auth.Login.Path, appHandlers.Auth.LoginWithPassword)

		authGroup := v1.Group("")
		authGroup.Use(middleware.AuthorizationMiddleware(appServices.Token))
		authGroup.Use(middleware.DynamicRBACGuardMiddleware(secRegistry))
		{
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketsList.Path, appHandlers.Ticket.GetTicketsListHandler)
			authGroup.POST(routes.APIRoutes.Users.GetUsersByIDs.Path, appHandlers.User.GetUsersByIDs)
			authGroup.POST(routes.APIRoutes.Users.GetUserByID.Path, appHandlers.User.GetUserByID)
			authGroup.POST(routes.APIRoutes.Users.GetUserByUsername.Path, appHandlers.User.GetUserByUsername)
			authGroup.POST(routes.APIRoutes.Tickets.GetTicketByID.Path, appHandlers.Ticket.GetTicketByIDHandler)
			authGroup.POST(routes.APIRoutes.Files.UploadTicketFile.Path, appHandlers.File.UploadTicketFileHandler)
			authGroup.POST(routes.APIRoutes.Files.GetDownloadLinkTicketFile.Path, appHandlers.File.GetDownloadLinkTicketFileHandler)
		}

		publicGroup := v1.Group("")
		{
			publicGroup.GET(routes.APIRoutes.Captcha.GetCaptcha.Path, appHandlers.Captcha.GenerateCaptchaHandler)
			publicGroup.POST(routes.APIRoutes.Captcha.VerifyCaptcha.Path, appHandlers.Captcha.VerifyCaptchaHandler)
			publicGroup.GET(routes.APIRoutes.Auth.LoginWithSingleUseToken.Path, appHandlers.Auth.LoginWithOneTimeToken)
			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketTypes.Path, appHandlers.Ticket.GetAllActiveTicketTypesHandler)
			publicGroup.GET(routes.APIRoutes.Tickets.GetAllActiveTicketStatuses.Path, appHandlers.Ticket.GetAllActiveTicketStatusesHandler)
			publicGroup.GET(routes.APIRoutes.Departments.GetAllActiveDepartments.Path, appHandlers.Department.GetAllActiveDepartmentsHandler)
			publicGroup.POST("tickets/CloseTicket/", appHandlers.Ticket.CloseTicketHandler)
		}
	}

	return &TestAppBundle{
		DB:               db,
		Repos:            repos,
		Services:         appServices,
		Handlers:         appHandlers,
		SecurityRegistry: secRegistry,
		Engine:           g,
	}
}

// GenerateTestAuthToken returns a signed JWT for testing authenticated requests
func GenerateTestAuthToken(t *testing.T, tokenSvc *token.TokenService, userID int64, username string, roleIDs []int64) string {
	t.Helper()
	tokenStr, err := tokenSvc.NewAuthToken(token.AuthClaims{
		UserID:   userID,
		Username: username,
		RoleIDs:  roleIDs,
	})
	if err != nil {
		t.Fatalf("failed generating test auth token: %v", err)
	}
	return tokenStr
}

// GenerateTestCaptchaToken generates a valid Captcha token for testing captcha-guarded endpoints
func GenerateTestCaptchaToken(t *testing.T, tokenSvc *token.TokenService, ip string) string {
	t.Helper()
	tokenStr, err := tokenSvc.NewCaptchaToken(ip)
	if err != nil {
		t.Fatalf("failed generating test captcha token: %v", err)
	}
	return tokenStr
}
