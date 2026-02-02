package app

import (
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/infrastructure/postgres"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/handlers"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/auth"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/calls"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/history"
)

type App struct {
	userRepo domain.UserRepository
	callRepo domain.CallRepository
	router   *http.Router
}

func New() *App {
	userRepo := postgres.NewUserRepository()
	callRepo := postgres.NewCallRepository()

	registerUC := auth.NewRegisterUseCase(userRepo)
	loginUC := auth.NewLoginUseCase(userRepo)
	logoutUC := auth.NewLogoutUseCase()
	startCallUC := calls.NewStartCallUseCase(callRepo)
	endCallUC := calls.NewEndCallUseCase(callRepo)
	listHistoryUC := history.NewListHistoryUseCase(callRepo)

	authHandler := handlers.NewAuthHandler(registerUC, loginUC, logoutUC)
	callsHandler := handlers.NewCallsHandler(startCallUC, endCallUC)
	historyHandler := handlers.NewHistoryHandler(listHistoryUC)

	router := http.NewRouter(authHandler, callsHandler, historyHandler)

	return &App{
		userRepo: userRepo,
		callRepo: callRepo,
		router:   router,
	}
}

func (a *App) Router() *http.Router {
	return a.router
}
