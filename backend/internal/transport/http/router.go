package http

import (
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/handlers"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	auth    *handlers.AuthHandler
	calls   *handlers.CallsHandler
	history *handlers.HistoryHandler
}

func NewRouter(auth *handlers.AuthHandler, calls *handlers.CallsHandler, history *handlers.HistoryHandler) *Router {
	return &Router{
		auth:    auth,
		calls:   calls,
		history: history,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	engine.GET("/system/health", handlers.Health)

	authGroup := engine.Group("/auth")
	{
		authGroup.POST("/register", r.auth.Register)
		authGroup.POST("/login", r.auth.Login)
		authGroup.POST("/logout", middleware.Auth(), r.auth.Logout)
	}

	apiGroup := engine.Group("/calls")
	apiGroup.Use(middleware.Auth())
	{
		apiGroup.POST("/start", r.calls.Start)
		apiGroup.POST("/end", r.calls.End)
		apiGroup.GET("/history", r.history.List)
	}
}
