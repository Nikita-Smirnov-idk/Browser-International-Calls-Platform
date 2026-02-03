package app

import (
	"fmt"
	"log"
	"os"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/config"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/infrastructure/jwt"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/infrastructure/postgres"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/infrastructure/voip"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/handlers"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/auth"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/calls"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/history"
	"gorm.io/gorm"
)

type App struct {
	userRepo   domain.UserRepository
	callRepo   domain.CallRepository
	voipClient voip.Client
	router     *http.Router
	db         *gorm.DB
	config     *config.Config
}

func New(cfg *config.Config) (*App, error) {
	db, err := postgres.NewConnection(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	migrationsPath := findMigrationsPath()
	if migrationsPath == "" {
		wd, _ := os.Getwd()
		log.Printf("Warning: migrations directory not found in working directory: %s, skipping migrations", wd)
	} else {
		if err := postgres.RunMigrations(db, migrationsPath); err != nil {
			log.Printf("Warning: failed to run migrations: %v", err)
		}
	}

	userRepo := postgres.NewUserRepository(db)
	callRepo := postgres.NewCallRepository(db)

	voipClient, err := voip.NewClient(&voip.Config{
		Provider:   cfg.VoIP.Provider,
		AccountSID: cfg.VoIP.AccountSID,
		AuthToken:  cfg.VoIP.AuthToken,
		APIKey:     cfg.VoIP.APIKey,
		FromNumber: cfg.VoIP.FromNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize voip client: %w", err)
	}

	jwtService := jwt.NewService(cfg.JWT.Secret)

	registerUC := auth.NewRegisterUseCase(userRepo)
	loginUC := auth.NewLoginUseCase(userRepo, jwtService)
	logoutUC := auth.NewLogoutUseCase()
	startCallUC := calls.NewStartCallUseCase(callRepo)
	endCallUC := calls.NewEndCallUseCase(callRepo)
	initiateCallUC := calls.NewInitiateCallUseCase(callRepo, voipClient)
	terminateCallUC := calls.NewTerminateCallUseCase(callRepo, voipClient)
	listHistoryUC := history.NewListHistoryUseCase(callRepo)

	authHandler := handlers.NewAuthHandler(registerUC, loginUC, logoutUC, jwtService)
	callsHandler := handlers.NewCallsHandler(startCallUC, endCallUC)
	webrtcHandler := handlers.NewWebRTCHandler(initiateCallUC, terminateCallUC)
	historyHandler := handlers.NewHistoryHandler(listHistoryUC)

	router := http.NewRouter(authHandler, callsHandler, webrtcHandler, historyHandler, jwtService)

	return &App{
		userRepo:   userRepo,
		callRepo:   callRepo,
		voipClient: voipClient,
		router:     router,
		db:         db,
		config:     cfg,
	}, nil
}

func (a *App) Router() *http.Router {
	return a.router
}

func (a *App) Close() error {
	if err := a.voipClient.Close(); err != nil {
		log.Printf("Error closing VoIP client: %v", err)
	}
	return postgres.Close(a.db)
}

func findMigrationsPath() string {
	possiblePaths := []string{
		"migrations",
		"./migrations",
		"../migrations",
		"backend/migrations",
		"./backend/migrations",
	}
	
	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			if files, err := os.ReadDir(path); err == nil {
				for _, file := range files {
					if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".sql" {
						return path
					}
				}
			}
		}
	}
	return ""
}
