package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kms-connect.com/mock-ehr/pkg/controller"
)

func main() {
	engine := gin.Default()
	controller.CernerController(engine)

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	engine.Run(":9999")
}
