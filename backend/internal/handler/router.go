package handler

import (
	"net/http"
	"time"

	"legalpermit/internal/config"
	"legalpermit/internal/middleware"
	"legalpermit/internal/model"
	"legalpermit/internal/ocr"
	"legalpermit/internal/repository"
	"legalpermit/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter wires repositories, services, handlers and routes together.
func NewRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	stepRepo := repository.NewStepRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	ptRepo := repository.NewPTRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	deadlineRepo := repository.NewDeadlineRepository(db)

	// Infrastructure
	tokenMgr := middleware.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	ocrProvider := ocr.NewMockProvider()

	// Services
	authSvc := service.NewAuthService(userRepo, tokenMgr)
	deadlineSvc := service.NewDeadlineService(deadlineRepo)
	projectSvc := service.NewProjectService(projectRepo, stepRepo, deadlineSvc)
	stepSvc := service.NewStepService(stepRepo)
	docSvc := service.NewDocumentService(docRepo, cfg.UploadDir)
	ptSvc := service.NewPTService(ptRepo, cfg.UploadDir)
	settingSvc := service.NewSettingService(settingRepo)
	dashboardSvc := service.NewDashboardService(stepRepo, deadlineSvc)

	// Seed defaults: DACI/notification settings and Master Deadline rules.
	if err := settingSvc.EnsureDefaults(); err != nil {
		panic(err)
	}
	if err := deadlineSvc.EnsureDefaults(); err != nil {
		panic(err)
	}

	// Handlers
	authH := NewAuthHandler(authSvc)
	projectH := NewProjectHandler(projectSvc)
	stepH := NewStepHandler(stepSvc, docSvc)
	ptH := NewPTHandler(ptSvc)
	settingH := NewSettingHandler(settingSvc)
	dashboardH := NewDashboardHandler(dashboardSvc, docSvc)
	ocrH := NewOCRHandler(ocrProvider)
	deadlineH := NewDeadlineHandler(deadlineSvc)

	r := gin.Default()
	r.MaxMultipartMemory = 16 << 20 // 16 MiB

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		api.POST("/auth/login", authH.Login)

		authed := api.Group("")
		authed.Use(middleware.Auth(tokenMgr))
		{
			authed.GET("/auth/me", authH.Me)

			// Reference metadata (category labels) for the UI.
			authed.GET("/meta/categories", func(c *gin.Context) {
				c.JSON(http.StatusOK, service.CategoryLabels)
			})

			// Projects & steps (Proses A–I).
			authed.GET("/projects", projectH.List)
			authed.POST("/projects", projectH.Create)
			authed.GET("/projects/:id", projectH.Get)
			authed.GET("/projects/:id/progress", projectH.Progress)

			authed.GET("/steps/:id", stepH.Get)
			authed.PUT("/steps/:id", stepH.Update)
			authed.POST("/steps/:id/documents", stepH.UploadDocument)
			authed.GET("/documents/:id/download", stepH.DownloadDocument)

			// Master Data PT (Proses E).
			authed.GET("/pt", ptH.List)
			authed.POST("/pt", ptH.Create)
			authed.GET("/pt/:id", ptH.Get)
			authed.POST("/pt/:id/documents", ptH.UploadDocument)
			authed.GET("/pt-documents/:id/download", ptH.DownloadDocument)

			// Dashboard: early warning + document search.
			authed.GET("/dashboard/warnings", dashboardH.EarlyWarnings)
			authed.GET("/dashboard/documents", dashboardH.SearchDocuments)

			// OCR AI extraction.
			authed.POST("/ocr/extract", ocrH.Extract)

			// Settings (DINAMIS SETTING KADEP) — write restricted to KADEP/DIROPS.
			authed.GET("/settings/daci", settingH.GetDACI)
			authed.GET("/settings/notification", settingH.GetNotification)
			authed.GET("/deadline-master", deadlineH.List)
			settingsWrite := authed.Group("")
			settingsWrite.Use(middleware.RequireRole(model.RoleKadep, model.RoleDirops))
			{
				settingsWrite.PUT("/settings/daci", settingH.SetDACI)
				settingsWrite.PUT("/settings/notification", settingH.SetNotification)
				settingsWrite.PUT("/deadline-master", deadlineH.Update)
			}
		}
	}

	return r
}
