package main

import (
	"github.com/gin-gonic/gin"
	"kms-connect.com/mock-ehr/pkg/controller"
)

func main() {
	engine := gin.Default()
	controller.CernerController(engine)
	engine.Run(":9999")
}
