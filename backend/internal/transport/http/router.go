package http

import (
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/handlers"
	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	auth       *handlers.AuthHandler
	calls      *handlers.CallsHandler
	webrtc     *handlers.WebRTCHandler
	voice      *handlers.VoiceHandler
	history    *handlers.HistoryHandler
	jwtService middleware.JWTService
}

func NewRouter(auth *handlers.AuthHandler, calls *handlers.CallsHandler, webrtc *handlers.WebRTCHandler, voice *handlers.VoiceHandler, history *handlers.HistoryHandler, jwtService middleware.JWTService) *Router {
	return &Router{
		auth:       auth,
		calls:      calls,
		webrtc:     webrtc,
		voice:      voice,
		history:    history,
		jwtService: jwtService,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	engine.GET("/system/health", handlers.Health)

	api := engine.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", r.auth.Register)
			authGroup.POST("/login", r.auth.Login)
			authGroup.POST("/logout", middleware.Auth(r.jwtService), r.auth.Logout)
		}

		callsGroup := api.Group("/calls")
		callsGroup.Use(middleware.Auth(r.jwtService))
		{
			callsGroup.POST("", r.calls.Create)
			callsGroup.PUT("/:id", r.calls.Update)
			callsGroup.GET("/history", r.history.List)
			callsGroup.POST("/initiate", r.webrtc.Initiate)
			callsGroup.POST("/terminate", r.webrtc.Terminate)
		}

		if r.voice != nil {
			api.POST("/voice/token", middleware.Auth(r.jwtService), r.voice.Token)
		}
	}

	if r.voice != nil {
		engine.GET("/api/voice/twiml", r.voice.TwiML)
		engine.POST("/api/voice/twiml", r.voice.TwiML)
	}
}
