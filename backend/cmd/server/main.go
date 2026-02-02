package main

import (
	"log"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/app"
	"github.com/gin-gonic/gin"
)

func main() {
	application := app.New()
	engine := gin.Default()
	application.Router().Setup(engine)
	log.Fatal(engine.Run(":8080"))
}
